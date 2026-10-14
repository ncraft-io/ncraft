package encoder_test

import (
	"testing"

	"github.com/ncraft-io/ncraft/go/pkg/ncraft/config/encoder"
	"github.com/ncraft-io/ncraft/go/pkg/ncraft/config/encoder/hcl"
	"github.com/ncraft-io/ncraft/go/pkg/ncraft/config/encoder/json"
	"github.com/ncraft-io/ncraft/go/pkg/ncraft/config/encoder/toml"
	"github.com/ncraft-io/ncraft/go/pkg/ncraft/config/encoder/xml"
	"github.com/ncraft-io/ncraft/go/pkg/ncraft/config/encoder/yaml"
	"github.com/stretchr/testify/require"
)

type settings struct {
	Name    string
	Count   int
	Enabled bool
}

func TestEncodersRoundTripAndInvalidInput(t *testing.T) {
	for _, tc := range []struct {
		name, invalid string
		codec         encoder.Encoder
	}{
		{"json", "{", json.NewEncoder()},
		{"yaml", "name: [", yaml.NewEncoder()},
		{"toml", "name = [", toml.NewEncoder()},
		{"xml", "<settings>", xml.NewEncoder()},
		{"hcl", "name = [", hcl.NewEncoder()},
	} {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.name, tc.codec.String())
			want := settings{"service", 3, true}
			data, err := tc.codec.Encode(want)
			require.NoError(t, err)
			var got settings
			require.NoError(t, tc.codec.Decode(data, &got))
			require.Equal(t, want, got)
			require.Error(t, tc.codec.Decode([]byte(tc.invalid), &got))
		})
	}
}
