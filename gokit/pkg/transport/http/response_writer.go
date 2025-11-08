package http

import (
    "context"
    "net/http"
)

type ResponseWriter interface {
    WriteHttpResponse(ctx context.Context, writer http.ResponseWriter) error
}
