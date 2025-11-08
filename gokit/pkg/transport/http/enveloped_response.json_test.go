package http

import (
	"testing"

	jsoniter "github.com/json-iterator/go"
	"github.com/mojo-lang/core/go/pkg/mojo/core"
	"github.com/stretchr/testify/assert"
)

func TestEnvelopedResponseCodec_Decode(t *testing.T) {
	json := `{
		"code": "400",
		"message": "Invalid Arguments",
		"data": null
	}`

	resp := &EnvelopedResponse{}
	err := jsoniter.UnmarshalFromString(json, resp)
	assert.NoError(t, err)
	assert.NotNil(t, resp.Error)
}

func TestEnvelopedResponseCodec_Encode(t *testing.T) {
	resp := &EnvelopedResponse{
		Error: core.NewErrorFrom(400, "Invalid Arguments"),
	}

	out, err := jsoniter.Marshal(resp)
	assert.NoError(t, err)
	assert.NotEmpty(t, out)

	response := &EnvelopedResponse{}
	err = jsoniter.Unmarshal(out, response)
	assert.NoError(t, err)
	assert.Equal(t, resp.Error.Message, response.Error.Message)
}
