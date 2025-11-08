package http

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResponseJsonWriter_WriteHttpResponse(t *testing.T) {
	const out = "string response"
	handler := func(w http.ResponseWriter, r *http.Request) {
		NewResponseJsonWriter(out).WriteHttpResponse(context.Background(), w)
	}

	req := httptest.NewRequest("GET", "https://example.com/foo", nil)
	r := httptest.NewRecorder()
	handler(r, req)

	resp := r.Result()
	body, _ := io.ReadAll(resp.Body)
	assert.Equal(t, `"string response"`, string(body))
}
