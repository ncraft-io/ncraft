package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"sync"
	"time"

	"google.golang.org/grpc"
)

// DefaultShutdownTimeout bounds transport draining and service cleanup together.
const DefaultShutdownTimeout = 30 * time.Second

// Starter is an optional service startup hook. Start completes before transports
// accept requests. On failure, Start must release any resources it acquired;
// Shutdown is called only for services whose startup completed successfully.
type Starter interface {
	Start(Config) error
}

// Shutdowner releases service resources after the generated transports drain.
// Implementations must respect the deadline and cancellation of ctx.
type Shutdowner interface {
	Shutdown(context.Context) error
}

// ErrorSource reports fatal background failures to the hosting server. Ordinary
// request errors should be returned to the caller instead. The service owns the
// channel and should buffer it so reporting the first failure cannot block during
// startup. Errors must return the same channel throughout a service's lifetime.
// Nil errors, nil channels and channel closure do not request server shutdown.
type ErrorSource interface {
	Errors() <-chan error
}

// WaitForError waits for a non-nil transport/service error, or returns nil when
// ctx is canceled (for example by SIGTERM). It does not close service-owned
// channels or create forwarding goroutines. Closed channels are disabled.
func WaitForError(ctx context.Context, transportErrors <-chan error, services ...interface{}) error {
	cases := []reflect.SelectCase{
		{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ctx.Done())},
		{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(transportErrors)},
	}
	names := []string{"", "transport"}
	for _, service := range services {
		if source, ok := service.(ErrorSource); ok {
			cases = append(cases, reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(source.Errors())})
			names = append(names, fmt.Sprintf("service %T", service))
		}
	}
	for {
		selected, value, ok := reflect.Select(cases)
		if selected == 0 {
			return nil
		}
		if !ok {
			cases[selected].Chan = reflect.Value{}
			continue
		}
		if value.IsNil() {
			continue
		}
		return fmt.Errorf("%s: %w", names[selected], value.Interface().(error))
	}
}

// ShutdownServers drains HTTP and gRPC concurrently under one deadline. HTTP
// connections are closed and gRPC is force-stopped if graceful draining times out.
// Call service Shutdown hooks only after this function returns.
func ShutdownServers(ctx context.Context, grpcServer *grpc.Server, httpServers ...*http.Server) error {
	var wg sync.WaitGroup
	errc := make(chan error, len(httpServers)+1)
	for _, server := range httpServers {
		if server == nil {
			continue
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := server.Shutdown(ctx); err != nil {
				closeErr := server.Close()
				errc <- fmt.Errorf("shutdown HTTP server %s: %w", server.Addr, errors.Join(err, closeErr))
			}
		}()
	}
	if grpcServer != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			done := make(chan struct{})
			go func() { grpcServer.GracefulStop(); close(done) }()
			select {
			case <-done:
			case <-ctx.Done():
				grpcServer.Stop()
				errc <- fmt.Errorf("shutdown gRPC server: %w", ctx.Err())
			}
		}()
	}
	wg.Wait()
	close(errc)
	var errs []error
	for err := range errc {
		errs = append(errs, err)
	}
	return errors.Join(errs...)
}
