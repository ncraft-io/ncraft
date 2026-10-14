package source

import (
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
)

func TestNoopWatcherStopIsIdempotent(t *testing.T) {
	w, err := NewNoopWatcher()
	require.NoError(t, err)
	var wg sync.WaitGroup
	for range 10 {
		wg.Go(func() { _ = w.Stop() })
	}
	wg.Wait()
	_, err = w.Next()
	require.ErrorIs(t, err, ErrWatcherStopped)
}
