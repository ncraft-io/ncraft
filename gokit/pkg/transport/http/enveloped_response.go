package http

import "github.com/mojo-lang/core/go/pkg/mojo/core"

type EnvelopedResponse struct {
	Error *core.Error `json:"error"`
	Data  interface{} `json:"data"`

	TotalCount    int32  `json:"totalCount,omitempty"`
	NextPageToken string `json:"nextPageToken,omitempty"`
}

func (r *EnvelopedResponse) ToErrorWrapped() *ErrorWrappedEnvelopedResponse {
	if r != nil {
		return (*ErrorWrappedEnvelopedResponse)(r)
	}
	return nil
}
