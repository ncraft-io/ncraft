package storage

import (
	"github.com/mojo-lang/mojo/go/pkg/mojo/core"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestObject_GetHttpHeaders(t *testing.T) {
	obj := &Object{
		Etag:         "1234",
		Key:          "",
		LastModified: nil,
		Size:         0,
		ContentType: &core.MediaType{
			Type:    "application",
			Subtype: "zip",
		},
	}

	headers := obj.GetHttpHeaders()
	assert.Equal(t, 3, len(headers.Vals))
}
