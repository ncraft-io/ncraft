package storage

import (
    "context"
    "fmt"
    "github.com/mojo-lang/http/go/pkg/mojo/http"
    gohttp "net/http"
)

func (x *Object) GetHttpHeaders() *http.Headers {
    if x != nil {
        headers := http.NewHeaders()
        if x.ContentType != nil {
            headers.Set("Content-Type", x.ContentType.Format())
            if x.ContentType.Subtype == "zip" {
                headers.Set("Content-Encoding", "gzip")
            }
        }
        if len(x.Etag) > 0 {
            headers.Set("ETag", x.Etag)
        }
        return headers
    }
    return nil
}

func (x *Object) WriteHttpResponse(ctx context.Context, writer gohttp.ResponseWriter) error {
    _ = ctx
    if x != nil {
        x.GetHttpHeaders().SyncTo(writer.Header())
        if count, err := writer.Write(x.Content); err != nil {
            return err
        } else if count < len(x.Content) {
            return fmt.Errorf("failed to write completed content, expect :%d, actual: %d", len(x.Content), count)
        }
        return nil
    }

    return nil
}
