package nats

import (
	"github.com/ncraft-io/ncraft/go/pkg/ncraft/messaging"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
)

func TestShutdownWithoutJetStreamAndRepeatedly(t *testing.T) {
	for _, n := range []*Nats{nil, {}, {shutdownCh: make(chan struct{})}} {
		var wg sync.WaitGroup
		for range 10 {
			wg.Go(func() { n.Shutdown() })
		}
		wg.Wait()
		if n != nil && n.shutdownCh != nil {
			select {
			case <-n.shutdownCh:
			default:
				t.Fatal("shutdown channel still open")
			}
		}
	}
}

func TestConfigurationHelpers(t *testing.T) {
	_, err := New(nil)
	require.Error(t, err)
	var n *Nats
	require.Empty(t, n.StreamName())
	require.Nil(t, n.SubjectNames())
	require.Zero(t, n.MaxMsgs())
	require.Zero(t, n.MaxAge())
	require.Nil(t, n.GetConn())
	require.Nil(t, n.GetJetStream())
	n = &Nats{config: &messaging.Config{}}
	n.AppendSubjectName("one")
	n.AppendSubjectName("two")
	require.Equal(t, []string{"one", "two"}, n.SubjectNames())
	n.SetSubjectNames("three")
	require.Equal(t, []string{"three"}, n.SubjectNames())
}
