package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	httptransport "github.com/go-kit/kit/transport/http"
)

type testResponse struct{ value string }

func (r *testResponse) WriteHttpResponse(_ context.Context, w http.ResponseWriter) error {
	_, err := fmt.Fprint(w, r.value)
	return err
}

func TestResponseBindingValidationAndIdentity(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	ctx := RequestToContext(context.Background(), r)
	if request, ok := RequestFromContext(ctx); !ok || request != r {
		t.Fatal("request missing")
	}
	response, writer := &testResponse{"original"}, &testResponse{"stream"}
	var nilResponse *testResponse
	for _, bad := range []interface{}{nil, nilResponse, 1, "value", []byte("x"), map[string]int{}} {
		if err := BindResponseWriter(ctx, bad, writer); err == nil {
			t.Errorf("accepted response %T", bad)
		}
		if _, ok := BoundResponseWriter(ctx, bad); ok {
			t.Fatal("resolved invalid response")
		}
	}
	for _, bad := range []ResponseWriter{nil, nilResponse} {
		if err := BindResponseWriter(ctx, response, bad); err == nil {
			t.Fatal("accepted nil writer")
		}
	}
	if err := BindResponseWriter(context.Background(), response, writer); err == nil {
		t.Fatal("accepted non-HTTP context")
	}
	if err := BindResponseWriter(ctx, response, writer); err != nil {
		t.Fatal(err)
	}
	if got, ok := BoundResponseWriter(ctx, response); !ok || got != writer {
		t.Fatal("binding missing")
	}
	if _, ok := BoundResponseWriter(ctx, &testResponse{"original"}); ok {
		t.Fatal("matched a copied response")
	}
	other := RequestToContext(ctx, r)
	if _, ok := BoundResponseWriter(other, response); ok {
		t.Fatal("inherited another request's binding")
	}
	if got, ok := BoundResponseWriter(context.WithValue(ctx, struct{ name string }{"child"}, true), response); !ok || got != writer {
		t.Fatal("derived context lost binding")
	}
	newResponse := &testResponse{"retry"}
	if err := BindResponseWriter(ctx, newResponse, writer); err != nil {
		t.Fatal(err)
	}
	if _, ok := BoundResponseWriter(ctx, response); ok {
		t.Fatal("old retry binding remained")
	}
}

func TestBoundWriterThroughGoKitServer(t *testing.T) {
	for _, mode := range []string{"stream", "replace", "error"} {
		t.Run(mode, func(t *testing.T) {
			server := httptransport.NewServer(
				func(ctx context.Context, _ interface{}) (interface{}, error) {
					result := &testResponse{"original"}
					if err := BindResponseWriter(ctx, result, &testResponse{"stream"}); err != nil {
						return nil, err
					}
					switch mode {
					case "replace":
						return &testResponse{"replacement"}, nil
					case "error":
						return result, errors.New("denied")
					}
					return result, nil
				},
				func(context.Context, *http.Request) (interface{}, error) { return nil, nil },
				func(ctx context.Context, w http.ResponseWriter, result interface{}) error {
					if writer, ok := BoundResponseWriter(ctx, result); ok {
						return writer.WriteHttpResponse(ctx, w)
					}
					return result.(ResponseWriter).WriteHttpResponse(ctx, w)
				},
				httptransport.ServerBefore(RequestToContext),
				httptransport.ServerErrorEncoder(func(_ context.Context, _ error, w http.ResponseWriter) { http.Error(w, "denied", http.StatusForbidden) }),
			)
			w := httptest.NewRecorder()
			server.ServeHTTP(w, httptest.NewRequest("GET", "/", nil))
			want := map[string]string{"stream": "stream", "replace": "replacement", "error": "denied\n"}[mode]
			if w.Body.String() != want {
				t.Fatalf("body = %q, want %q", w.Body.String(), want)
			}
			if mode == "error" && w.Code != http.StatusForbidden {
				t.Fatalf("status = %d", w.Code)
			}
		})
	}
}

func TestResponseBindingsConcurrentRequests(t *testing.T) {
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			ctx := RequestToContext(context.Background(), httptest.NewRequest("GET", "/", nil))
			result := &testResponse{fmt.Sprint(i)}
			if err := BindResponseWriter(ctx, result, result); err != nil {
				t.Error(err)
				return
			}
			if writer, ok := BoundResponseWriter(ctx, result); !ok || writer != result {
				t.Error("request binding crossed")
			}
		}(i)
	}
	wg.Wait()
}
