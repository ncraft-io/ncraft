package reader

import (
	"github.com/ncraft-io/ncraft/go/pkg/ncraft/config/encoder"
	"github.com/ncraft-io/ncraft/go/pkg/ncraft/config/encoder/hcl"
	"github.com/ncraft-io/ncraft/go/pkg/ncraft/config/encoder/json"
	"github.com/ncraft-io/ncraft/go/pkg/ncraft/config/encoder/toml"
	"github.com/ncraft-io/ncraft/go/pkg/ncraft/config/encoder/xml"
	"github.com/ncraft-io/ncraft/go/pkg/ncraft/config/encoder/yaml"
)

type Options struct {
	Encodings map[string]encoder.Encoder
}

type Option func(o *Options)

func NewOptions(opts ...Option) Options {
	options := Options{
		Encodings: map[string]encoder.Encoder{
			"hcl":  hcl.NewEncoder(),
			"json": json.NewEncoder(),
			"toml": toml.NewEncoder(),
			"xml":  xml.NewEncoder(),
			"yaml": yaml.NewEncoder(),
			"yml":  yaml.NewEncoder(),
		},
	}
	for _, o := range opts {
		o(&options)
	}
	return options
}

func WithEncoder(e encoder.Encoder) Option {
	return func(o *Options) {
		if o.Encodings == nil {
			o.Encodings = make(map[string]encoder.Encoder)
		}
		o.Encodings[e.String()] = e
	}
}
