package hcl

import (
	jsoniter "github.com/json-iterator/go"

	"github.com/ncraft-io/ncraft/go/pkg/ncraft/config/encoder"

	"github.com/hashicorp/hcl"
)

type hclEncoder struct{}

func (h hclEncoder) Encode(v interface{}) ([]byte, error) {
	return jsoniter.Marshal(v)
}

func (h hclEncoder) Decode(d []byte, v interface{}) error {
	return hcl.Unmarshal(d, v)
}

func (h hclEncoder) String() string {
	return "hcl"
}

func NewEncoder() encoder.Encoder {
	return hclEncoder{}
}
