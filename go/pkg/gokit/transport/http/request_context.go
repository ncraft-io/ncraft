package http

import (
	"context"
	"errors"
	"net/http"
	"reflect"
	"sync"
)

type requestContextKey struct{}

type requestContext struct {
	request  *http.Request
	mu       sync.RWMutex
	response interface{}
	writer   ResponseWriter
}

// RequestToContext is a go-kit ServerBefore hook. Each invocation creates an
// independent response attachment, even when the incoming context already has one.
func RequestToContext(ctx context.Context, request *http.Request) context.Context {
	return context.WithValue(ctx, requestContextKey{}, &requestContext{request: request})
}

// RequestFromContext returns the HTTP request installed by RequestToContext.
// Treat the request as read-only; use ctx for cancellation and downstream calls.
func RequestFromContext(ctx context.Context) (*http.Request, bool) {
	state, _ := ctx.Value(requestContextKey{}).(*requestContext)
	if state == nil || state.request == nil {
		return nil, false
	}
	return state.request, true
}

// BindResponseWriter attaches an HTTP representation to a non-nil response
// pointer without changing the service's return type. Only the exact pointer
// returned through the endpoint may use this writer. A later binding replaces
// the previous binding (for example when endpoint middleware retries a call).
//
// Writers must acquire resources lazily in WriteHttpResponse and release them
// before returning: middleware may discard the response without encoding it.
// Binding does not serialize mutations to the response or writer themselves.
func BindResponseWriter(ctx context.Context, response interface{}, writer ResponseWriter) error {
	state, _ := ctx.Value(requestContextKey{}).(*requestContext)
	if state == nil || state.request == nil {
		return errors.New("HTTP request context is not initialized")
	}
	value := reflect.ValueOf(response)
	if value.Kind() != reflect.Ptr || value.IsNil() {
		return errors.New("HTTP response binding requires a non-nil pointer")
	}
	if writer == nil || isNilWriter(writer) {
		return errors.New("HTTP response writer is nil")
	}
	state.mu.Lock()
	defer state.mu.Unlock()
	state.response, state.writer = response, writer
	return nil
}

func isNilWriter(writer ResponseWriter) bool {
	value := reflect.ValueOf(writer)
	switch value.Kind() {
	case reflect.Ptr, reflect.Map, reflect.Slice, reflect.Func, reflect.Chan, reflect.Interface:
		return value.IsNil()
	}
	return false
}

// BoundResponseWriter resolves an attachment only for the exact response pointer.
// Encoders must call this before checking ResponseWriter or applying JSON envelopes.
// Endpoint errors must use the ordinary error encoder instead.
func BoundResponseWriter(ctx context.Context, response interface{}) (ResponseWriter, bool) {
	state, _ := ctx.Value(requestContextKey{}).(*requestContext)
	if state == nil {
		return nil, false
	}
	value := reflect.ValueOf(response)
	if value.Kind() != reflect.Ptr || value.IsNil() {
		return nil, false
	}
	state.mu.RLock()
	defer state.mu.RUnlock()
	if state.writer == nil || state.response != response {
		return nil, false
	}
	return state.writer, true
}
