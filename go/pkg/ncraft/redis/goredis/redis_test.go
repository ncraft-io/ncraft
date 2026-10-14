package goredis

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"testing"
	"time"

	goRedis "github.com/go-redis/redis/v8"
	redigo "github.com/gomodule/redigo/redis"
	"github.com/ncraft-io/ncraft/go/pkg/ncraft/redis"
	"github.com/stretchr/testify/require"
)

type scriptStep struct {
	commands [][]string
	reply    string
}

// Use the real Redis wire protocol over an in-memory connection. No local Redis
// server, fixed ports, or sleeps are needed.
func scriptedClient(t *testing.T, steps ...scriptStep) *Redis {
	t.Helper()
	clientConn, serverConn := net.Pipe()
	done := make(chan error, 1)
	go func() {
		defer serverConn.Close()
		decoder := redigo.NewConn(serverConn, time.Second, time.Second)
		for _, step := range steps {
			for _, want := range step.commands {
				got, err := redigo.Strings(decoder.Receive())
				if err != nil {
					done <- err
					return
				}
				if fmt.Sprint(got) != fmt.Sprint(want) {
					done <- fmt.Errorf("command = %v, want %v", got, want)
					return
				}
			}
			if _, err := io.WriteString(serverConn, step.reply); err != nil {
				done <- err
				return
			}
		}
		done <- nil
	}()
	c := &Redis{Client: goRedis.NewClient(&goRedis.Options{
		Dialer:     func(context.Context, string, string) (net.Conn, error) { return clientConn, nil },
		MaxRetries: -1, PoolSize: 1,
	})}
	t.Cleanup(func() {
		require.NoError(t, c.Close())
		_ = serverConn.Close()
		require.NoError(t, <-done)
	})
	return c
}

func TestDoAndStats(t *testing.T) {
	c := scriptedClient(t,
		scriptStep{[][]string{{"SET", "key", "value"}}, "+OK\r\n"},
		scriptStep{[][]string{{"GET", "key"}}, "$5\r\nvalue\r\n"},
		scriptStep{[][]string{{"GET", "missing"}}, "$-1\r\n"},
	)
	got, err := c.Do(context.Background(), "SET", "key", "value")
	require.NoError(t, err)
	require.Equal(t, "OK", got)
	got, err = c.Do(context.Background(), "GET", "key")
	require.NoError(t, err)
	require.Equal(t, "value", got)
	_, err = c.Do(context.Background(), "GET", "missing")
	require.ErrorIs(t, err, goRedis.Nil)
	require.GreaterOrEqual(t, c.Stats().Hits, int64(1))
}

func TestBatchReuseAndErrors(t *testing.T) {
	c := scriptedClient(t,
		scriptStep{[][]string{{"PING"}, {"GET", "key"}}, "+PONG\r\n$5\r\nvalue\r\n"},
		scriptStep{[][]string{{"INCR", "count"}}, ":1\r\n"},
		scriptStep{[][]string{{"GET", "missing"}, {"PING"}}, "$-1\r\n+PONG\r\n"},
	)
	batch := c.NewBatch()
	require.NoError(t, batch.Put("PING"))
	require.NoError(t, batch.Put("GET", "key"))
	replies, err := c.RunBatch(context.Background(), batch)
	require.NoError(t, err)
	require.Equal(t, []any{"PONG", "value"}, replies)
	require.NoError(t, batch.Put("INCR", "count"))
	replies, err = c.RunBatch(context.Background(), batch)
	require.NoError(t, err)
	require.Equal(t, []any{int64(1)}, replies)
	require.NoError(t, batch.Put("GET", "missing"))
	require.NoError(t, batch.Put("PING"))
	replies, err = c.RunBatch(context.Background(), batch)
	require.ErrorIs(t, err, goRedis.Nil)
	require.Equal(t, []any{nil, "PONG"}, replies)
}

type foreignBatch struct{}

func (foreignBatch) Put(string, ...any) error { return nil }

func TestInvalidBatchAndCancelledRequest(t *testing.T) {
	client := goRedis.NewClient(&goRedis.Options{Dialer: func(ctx context.Context, _, _ string) (net.Conn, error) { return nil, ctx.Err() }, MaxRetries: -1})
	c := &Redis{Client: client}
	defer c.Close()
	for _, batch := range []redis.Batch{nil, foreignBatch{}, (*pipeline)(nil), &pipeline{}} {
		require.NotPanics(t, func() {
			_, err := c.RunBatch(context.Background(), batch)
			require.Error(t, err)
		})
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := c.Do(ctx, "PING")
	require.ErrorIs(t, err, context.Canceled)
}

func TestNewRedisOptions(t *testing.T) {
	cfg := &redis.Config{Connections: []string{"unused:6379"}, MinIdleConnections: -1, Password: "secret", DbNumber: 2, IdleTimeout: time.Minute}
	c := NewRedis(cfg).(*Redis)
	defer c.Close()
	opts := c.Client.(*goRedis.Client).Options()
	require.Equal(t, "unused:6379", opts.Addr)
	require.Equal(t, "secret", opts.Password)
	require.Equal(t, 2, opts.DB)
	require.Equal(t, 300, opts.PoolSize)
	require.Equal(t, time.Second, opts.ReadTimeout)
	require.Equal(t, time.Minute, opts.IdleTimeout)
	cluster := NewRedis(&redis.Config{Connections: []string{"one:6379", "two:6379"}, MinIdleConnections: -1}).(*Redis)
	defer cluster.Close()
	require.Equal(t, []string{"one:6379", "two:6379"}, cluster.Client.(*goRedis.ClusterClient).Options().Addrs)
}

func TestDoPropagatesConnectionError(t *testing.T) {
	expected := errors.New("dial failed")
	client := goRedis.NewClient(&goRedis.Options{Dialer: func(context.Context, string, string) (net.Conn, error) { return nil, expected }, MaxRetries: -1})
	c := &Redis{Client: client}
	defer c.Close()
	_, err := c.Do(context.Background(), "PING")
	require.ErrorIs(t, err, expected)
}
