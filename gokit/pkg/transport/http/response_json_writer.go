package http

import (
    "context"
    jsoniter "github.com/json-iterator/go"
    "net/http"
)

type ResponseJsonWriter struct {
    Response interface{}
}

func NewResponseJsonWriter(response interface{}) *ResponseJsonWriter {
    return &ResponseJsonWriter{Response: response}
}

func (r *ResponseJsonWriter) WriteHttpResponse(ctx context.Context, writer http.ResponseWriter) error {
    _ = ctx
    stream := jsoniter.NewStream(jsoniter.ConfigFastest, writer, 512)
    stream.WriteVal(r.Response)
    if err := stream.Flush(); err != nil {
        return err
    }

    return stream.Error
}
