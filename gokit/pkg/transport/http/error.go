package http

import (
    "github.com/pkg/errors"
    "net/http"
)

// Error satisfies the Headerer and StatusCoder interfaces in
// package github.com/go-kit/kit/transport/http.
type Error struct {
    error
    statusCode int
    headers    http.Header
}

func WrapError(e error, code int, message string, headers ...string) *Error {
    err := &Error{
        error:      errors.Wrap(e, message),
        statusCode: code,
        headers:    make(http.Header),
    }

    length := len(headers)
    if length > 0 && length%2 == 0 {
        for i := 0; i < length/2; i += 2 {
            err.headers.Add(headers[i], headers[i+1])
        }
    }
    return err
}

func (e Error) StatusCode() int {
    return e.statusCode
}

func (e Error) Headers() http.Header {
    return e.headers
}
