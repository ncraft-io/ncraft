package direct

import (
	"github.com/go-kit/kit/sd"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestInstancer(t *testing.T) {
	var empty *Client
	require.Nil(t, empty.Instancer("svc"))
	c := New(map[string]*Config{"nil": nil, "empty": {}, "svc": {Urls: []string{"host:1", "host:2"}}})
	for _, name := range []string{"missing", "nil", "empty"} {
		require.Nil(t, c.Instancer(name))
	}
	i := c.Instancer("svc")
	require.NotNil(t, i)
	c.m["svc"].Urls[0] = "changed"
	events := make(chan sd.Event, 1)
	i.Register(events)
	require.Equal(t, []string{"host:1", "host:2"}, (<-events).Instances)
	i.Deregister(events)
	i.Stop()
}
