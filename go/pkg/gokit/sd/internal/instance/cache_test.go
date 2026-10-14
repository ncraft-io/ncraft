package instance

import (
	"errors"
	"testing"

	"github.com/go-kit/kit/sd"
	"github.com/stretchr/testify/require"
)

func TestCacheOwnsEventsAndNotifiesChanges(t *testing.T) {
	c := NewCache()
	one, two := make(chan sd.Event, 8), make(chan sd.Event, 8)
	c.Register(one)
	c.Register(two)
	require.Empty(t, (<-one).Instances)
	<-two
	input := []string{"b", "a"}
	c.Update(sd.Event{Instances: input})
	require.Equal(t, []string{"b", "a"}, input, "Update must not sort caller's slice")
	input[0] = "changed"
	first := <-one
	first.Instances[0] = "observer mutation"
	require.Equal(t, []string{"a", "b"}, (<-two).Instances)
	require.Equal(t, []string{"a", "b"}, c.State().Instances)
	state := c.State()
	state.Instances[0] = "snapshot mutation"
	require.Equal(t, []string{"a", "b"}, c.State().Instances)
	c.Update(sd.Event{Instances: []string{"b", "a"}})
	require.Empty(t, one, "equivalent events must not be broadcast")
	c.Deregister(one)
	err := errors.New("discovery unavailable")
	c.Update(sd.Event{Err: err})
	require.Empty(t, one)
	require.ErrorIs(t, (<-two).Err, err)
	c.Stop()
}
