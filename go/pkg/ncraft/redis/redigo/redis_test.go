package redigo

import (
	"context"
	"errors"
	"testing"

	driver "github.com/gomodule/redigo/redis"
	"github.com/ncraft-io/ncraft/go/pkg/ncraft/redis"
	"github.com/stretchr/testify/require"
)

type fakeConn struct {
	driver.Conn
	sent    []string
	replies []any
	err     error
}

func (c *fakeConn) Err() error                              { return nil }
func (c *fakeConn) Close() error                            { return nil }
func (c *fakeConn) Do(cmd string, args ...any) (any, error) { return nil, c.err }
func (c *fakeConn) Send(cmd string, args ...any) error      { c.sent = append(c.sent, cmd); return c.err }
func (c *fakeConn) Flush() error                            { return c.err }
func (c *fakeConn) Receive() (any, error) {
	if c.err != nil {
		return nil, c.err
	}
	reply := c.replies[0]
	c.replies = c.replies[1:]
	return reply, nil
}

type otherBatch struct{}

func (otherBatch) Put(string, ...any) error { return nil }

func TestBatchCommandsAndReuse(t *testing.T) {
	conn := &fakeConn{replies: []any{"PONG", int64(1)}}
	r := &Redis{pool: &driver.Pool{MaxIdle: 1, Dial: func() (driver.Conn, error) { return conn, nil }}}
	defer r.Close()
	b := r.NewBatch()
	require.NoError(t, b.Put("PING"))
	replies, err := r.RunBatch(context.Background(), b)
	require.NoError(t, err)
	require.Equal(t, []any{"PONG"}, replies)
	require.NoError(t, b.Put("INCR", "key"))
	replies, err = r.RunBatch(context.Background(), b)
	require.NoError(t, err)
	require.Equal(t, []any{int64(1)}, replies)
	require.Equal(t, []string{"PING", "INCR"}, conn.sent)
	for _, batch := range []redis.Batch{nil, otherBatch{}, (*Batch)(nil), &Batch{}} {
		_, err := r.RunBatch(context.Background(), batch)
		require.Error(t, err)
	}
}

func TestConnectionAndBatchErrors(t *testing.T) {
	expected := errors.New("dial failed")
	r := &Redis{pool: &driver.Pool{Dial: func() (driver.Conn, error) { return nil, expected }}}
	defer r.Close()
	_, err := r.Do(context.Background(), "PING")
	require.ErrorIs(t, err, expected)
	b := r.NewBatch()
	require.NoError(t, b.Put("GET", "key"))
	_, err = r.RunBatch(context.Background(), b)
	require.ErrorIs(t, err, expected)
}

func TestNewRedisReturnsDialErrors(t *testing.T) {
	r := NewRedis(&redis.Config{Connections: []string{"missing-port"}})
	defer r.Close()
	require.NotPanics(t, func() {
		_, err := r.Do(context.Background(), "PING")
		require.Error(t, err)
	})
}
