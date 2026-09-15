package server

import (
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

type testErrorSource struct{ errors <-chan error }

func (s testErrorSource) Errors() <-chan error { return s.errors }

func TestWaitForError(t *testing.T) {
	fatal := errors.New("background failure")
	closed := make(chan error)
	close(closed)
	source := make(chan error, 2)
	source <- nil
	source <- fatal
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := WaitForError(ctx, closed, struct{}{}, testErrorSource{}, testErrorSource{closed}, testErrorSource{source}); !errors.Is(err, fatal) {
		t.Fatalf("service error = %v", err)
	}
	transport := make(chan error, 2)
	transport <- nil
	transport <- fatal
	if err := WaitForError(ctx, transport); !errors.Is(err, fatal) {
		t.Fatalf("transport error = %v", err)
	}
	// Closing all sources is not itself a failure and must not cause a busy loop.
	canceled, stop := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- WaitForError(canceled, closed, testErrorSource{closed}) }()
	select {
	case err := <-done:
		t.Fatalf("closed sources ended wait: %v", err)
	case <-time.After(20 * time.Millisecond):
	}
	stop()
	select {
	case err := <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(time.Second):
		t.Fatal("cancellation did not stop wait")
	}
}

func serveHTTP(t *testing.T, handler http.Handler) (*http.Server, string, <-chan error) {
	t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	server := &http.Server{Handler: handler}
	t.Cleanup(func() { server.Close() })
	done := make(chan error, 1)
	go func() { done <- server.Serve(listener) }()
	return server, "http://" + listener.Addr().String(), done
}

func TestShutdownServersDrainsHTTPBeforeReturning(t *testing.T) {
	entered := make(chan struct{})
	release := make(chan struct{})
	server, url, served := serveHTTP(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		close(entered)
		<-release
		io.WriteString(w, "finished")
	}))
	response := make(chan error, 1)
	go func() {
		resp, err := http.Get(url)
		if err == nil {
			defer resp.Body.Close()
			var body []byte
			body, err = io.ReadAll(resp.Body)
			if string(body) != "finished" {
				err = errors.New("incomplete response")
			}
		}
		response <- err
	}()
	<-entered
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	stopped := make(chan error, 1)
	go func() { stopped <- ShutdownServers(ctx, grpc.NewServer(), server) }()
	select {
	case err := <-stopped:
		t.Fatalf("shutdown returned before request finished: %v", err)
	case <-time.After(20 * time.Millisecond):
	}
	close(release)
	if err := <-response; err != nil {
		t.Fatal(err)
	}
	if err := <-stopped; err != nil {
		t.Fatal(err)
	}
	if err := <-served; !errors.Is(err, http.ErrServerClosed) {
		t.Fatalf("HTTP Serve: %v", err)
	}
}

// This RPC remains active until the connection is forcefully stopped.
type blockingHealth struct {
	healthpb.UnimplementedHealthServer
	entered chan struct{}
}

func (s *blockingHealth) Check(ctx context.Context, _ *healthpb.HealthCheckRequest) (*healthpb.HealthCheckResponse, error) {
	close(s.entered)
	<-ctx.Done()
	return nil, ctx.Err()
}

func TestShutdownServersForcesHTTPAndGRPCAtDeadline(t *testing.T) {
	entered := make(chan struct{})
	httpServer, url, _ := serveHTTP(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { close(entered); <-r.Context().Done() }))
	httpDone := make(chan struct{})
	go func() {
		resp, _ := http.Get(url)
		if resp != nil {
			resp.Body.Close()
		}
		close(httpDone)
	}()
	<-entered
	grpcServer := grpc.NewServer()
	health := &blockingHealth{entered: make(chan struct{})}
	healthpb.RegisterHealthServer(grpcServer, health)
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(grpcServer.Stop)
	go grpcServer.Serve(listener)
	conn, err := grpc.NewClient(listener.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	rpcDone := make(chan struct{})
	go func() {
		healthpb.NewHealthClient(conn).Check(context.Background(), &healthpb.HealthCheckRequest{})
		close(rpcDone)
	}()
	<-health.entered
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer cancel()
	if err := ShutdownServers(ctx, grpcServer, httpServer); !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("shutdown error: %v", err)
	}
	for _, done := range []<-chan struct{}{httpDone, rpcDone} {
		select {
		case <-done:
		case <-time.After(time.Second):
			t.Fatal("request remained open after forced shutdown")
		}
	}
}
