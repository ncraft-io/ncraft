package http

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWrapErrorPreservesAllHeaders(t *testing.T) {
	cause := errors.New("unavailable")
	err := WrapError(cause, http.StatusServiceUnavailable, "request failed", "Retry-After", "10", "X-Reason", "busy", "X-Reason", "retry")
	require.Equal(t, http.StatusServiceUnavailable, err.StatusCode())
	require.Equal(t, "request failed: unavailable", err.Error())
	require.Equal(t, "10", err.Headers().Get("Retry-After"))
	require.Equal(t, []string{"busy", "retry"}, err.Headers().Values("X-Reason"))
}
