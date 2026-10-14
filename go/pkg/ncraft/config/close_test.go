package config

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/ncraft-io/ncraft/go/pkg/ncraft/config/loader"
	"github.com/stretchr/testify/require"
)

type closeLoader struct {
	loader.Loader
	calls atomic.Int64
	err   error
}

func (l *closeLoader) Close() error { l.calls.Add(1); return l.err }

func TestCloseStopsLoaderExactlyOnce(t *testing.T) {
	expected := errors.New("close failed")
	l := &closeLoader{err: expected}
	c := &config{exit: make(chan bool), opts: Options{Loader: l}}
	results := make(chan error, 10)
	var wg sync.WaitGroup
	for range 10 {
		wg.Go(func() { results <- c.Close() })
	}
	wg.Wait()
	close(results)
	for err := range results {
		require.ErrorIs(t, err, expected)
	}
	require.Equal(t, int64(1), l.calls.Load())
	select {
	case <-c.exit:
	default:
		t.Fatal("config still running")
	}
}
