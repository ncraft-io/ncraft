package file

import (
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/ncraft-io/ncraft/go/pkg/ncraft/config/source"
	"github.com/stretchr/testify/require"
)

func TestWatcherStopAndMissingFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	_, err := NewSource(WithPath(path)).Watch()
	require.Error(t, err)
	require.NoError(t, os.WriteFile(path, []byte(`{}`), 0600))
	w, err := NewSource(WithPath(path)).Watch()
	require.NoError(t, err)
	done := make(chan error, 1)
	go func() { _, err := w.Next(); done <- err }()
	var wg sync.WaitGroup
	for range 5 {
		wg.Go(func() { _ = w.Stop() })
	}
	wg.Wait()
	select {
	case err := <-done:
		require.ErrorIs(t, err, source.ErrWatcherStopped)
	case <-time.After(time.Second):
		t.Fatal("stopped watcher is blocked")
	}
	_, err = w.Next()
	require.ErrorIs(t, err, source.ErrWatcherStopped)
}
