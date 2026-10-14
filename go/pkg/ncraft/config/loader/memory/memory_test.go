package memory

import (
	"sync"
	"testing"
	"time"

	"github.com/ncraft-io/ncraft/go/pkg/ncraft/config/source"
	sourceMemory "github.com/ncraft-io/ncraft/go/pkg/ncraft/config/source/memory"
	"github.com/stretchr/testify/require"
)

func TestSyncPreservesSourcesOnLaterLoad(t *testing.T) {
	first := sourceMemory.NewSource(sourceMemory.WithJSON([]byte(`{"first":1,"shared":"old"}`)))
	m := NewLoader(WithSource(first))
	defer m.Close()
	require.NoError(t, m.Sync())
	snapshot, err := m.Snapshot()
	require.NoError(t, err)
	snapshot.ChangeSet.Data[0] = '!'
	snapshot, err = m.Snapshot()
	require.NoError(t, err)
	require.JSONEq(t, `{"first":1,"shared":"old"}`, string(snapshot.ChangeSet.Data))
	require.NoError(t, m.Load(sourceMemory.NewSource(sourceMemory.WithJSON([]byte(`{"second":2,"shared":"new"}`)))))
	snapshot, err = m.Snapshot()
	require.NoError(t, err)
	require.JSONEq(t, `{"first":1,"second":2,"shared":"new"}`, string(snapshot.ChangeSet.Data))
}

func TestWatcherUpdatesAndStop(t *testing.T) {
	s := sourceMemory.NewSource(sourceMemory.WithJSON([]byte(`{"value":1}`)))
	m := NewLoader(WithSource(s))
	defer m.Close()
	w, err := m.Watch("value")
	require.NoError(t, err)
	require.NoError(t, s.Write(&source.ChangeSet{Format: "json", Data: []byte(`{"value":2}`)}))
	require.NoError(t, m.Sync())
	done := make(chan error, 1)
	go func() {
		snap, err := w.Next()
		if err == nil && string(snap.ChangeSet.Data) != "2" {
			t.Errorf("unexpected value: %s", snap.ChangeSet.Data)
		}
		done <- err
	}()
	select {
	case err := <-done:
		require.NoError(t, err)
	case <-time.After(time.Second):
		t.Fatal("watcher did not receive update")
	}
	require.NoError(t, w.Stop())
	_, err = w.Next()
	require.Error(t, err)
}

func TestConcurrentWatcherStopAndSync(t *testing.T) {
	m := NewLoader(WithSource(sourceMemory.NewSource(sourceMemory.WithJSON([]byte(`{}`)))))
	defer m.Close()
	for range 50 {
		w, err := m.Watch()
		require.NoError(t, err)
		var wg sync.WaitGroup
		wg.Go(func() { _ = m.Sync() })
		for range 4 {
			wg.Go(func() { _ = w.Stop() })
		}
		wg.Wait()
	}
}

func TestCloseUnblocksLoaderWatcher(t *testing.T) {
	m := NewLoader()
	w, err := m.Watch()
	require.NoError(t, err)
	done := make(chan error, 1)
	go func() { _, err := w.Next(); done <- err }()
	var wg sync.WaitGroup
	for range 5 {
		wg.Go(func() { _ = m.Close() })
	}
	wg.Wait()
	select {
	case err := <-done:
		require.Error(t, err)
	case <-time.After(time.Second):
		t.Fatal("Close did not stop loader watcher")
	}
}
