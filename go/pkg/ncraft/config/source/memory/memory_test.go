package memory

import (
	"sync"
	"testing"
	"time"

	"github.com/ncraft-io/ncraft/go/pkg/ncraft/config/source"
	"github.com/stretchr/testify/require"
)

func TestEmptyMemorySource(t *testing.T) {
	for _, s := range []source.Source{NewSource(), NewSource(WithChangeSet(nil))} {
		require.NotPanics(t, func() {
			cs, err := s.Read()
			require.NoError(t, err)
			require.NotNil(t, cs)
		})
	}
}

func TestMemorySourceOwnsChangeSets(t *testing.T) {
	data := []byte(`{"key":1}`)
	s := NewSource(WithJSON(data))
	data[0] = '!'
	cs, err := s.Read()
	require.NoError(t, err)
	require.JSONEq(t, `{"key":1}`, string(cs.Data))
	cs.Data[0] = '!'
	cs, err = s.Read()
	require.NoError(t, err)
	require.JSONEq(t, `{"key":1}`, string(cs.Data))

	w1, err := s.Watch()
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, w1.Stop()) })
	w2, err := s.Watch()
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, w2.Stop()) })
	require.NoError(t, s.Write(&source.ChangeSet{Format: "json", Data: []byte(`{"key":2}`)}))
	first, err := w1.Next()
	require.NoError(t, err)
	first.Data[0] = '!'
	second, err := w2.Next()
	require.NoError(t, err)
	require.JSONEq(t, `{"key":2}`, string(second.Data))
	cs, err = s.Read()
	require.NoError(t, err)
	require.Equal(t, cs.Sum(), cs.Checksum)
	require.Equal(t, "memory", cs.Source)
}

func TestStopUnblocksMemoryWatcher(t *testing.T) {
	w, err := NewSource(WithJSON([]byte(`{}`))).Watch()
	require.NoError(t, err)
	done := make(chan error, 1)
	go func() { _, err := w.Next(); done <- err }()
	var wg sync.WaitGroup
	for range 10 {
		wg.Go(func() { _ = w.Stop() })
	}
	wg.Wait()
	select {
	case err := <-done:
		require.ErrorIs(t, err, source.ErrWatcherStopped)
	case <-time.After(time.Second):
		t.Fatal("Stop did not unblock Next")
	}
}
