package http

import (
	"net/http"
)

type RequestDecoder interface {
	DecodeHttpRequest(request *http.Request) error
}
