package reader

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ncraft-io/ncraft/go/pkg/ncraft/config/encoder/json"
)

func TestWithEncoder(t *testing.T) {
	options := &Options{}
	WithEncoder(json.NewEncoder())(options)

	assert.NotEmpty(t, options.Encodings)
	assert.NotNil(t, options.Encodings["json"])
}
