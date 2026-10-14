package pagination

import (
	"testing"
	"time"

	"github.com/mr-tron/base58"
	"github.com/stretchr/testify/require"
)

func TestPositionTokenRoundTrip(t *testing.T) {
	for _, position := range []int32{0, 1, 255, 256, 65536, 0xffffff} {
		original := &PositionPageToken{Time: time.Unix(1<<33, 0), Position: position}
		parsed, err := ParsePageToken(original.Format())
		require.NoError(t, err)
		require.Equal(t, position, parsed.Value())
		require.Equal(t, original.Unix(), parsed.(*PositionPageToken).Unix())
	}
}

func TestPositionTokenRejectsMalformedInput(t *testing.T) {
	tokens := []string{"", "x123", "p", "p0OIl"}
	for _, size := range []int{1, 3, 7, 9, 16} {
		tokens = append(tokens, "p"+base58.Encode(make([]byte, size)))
	}
	for _, input := range tokens {
		t.Run(input, func(t *testing.T) {
			token := PositionPageToken{Time: time.Unix(123, 0), Position: 42}
			before := token
			require.NotPanics(t, func() { require.Error(t, token.Parse(input)) })
			require.Equal(t, before, token, "failed parse must not change the token")
			require.NotPanics(t, func() {
				got, err := ParsePageToken(input)
				require.Error(t, err)
				require.Nil(t, got)
			})
		})
	}
}

func TestCreatePositionTokenNumericTypes(t *testing.T) {
	for _, value := range []any{int(42), int8(42), int16(42), int32(42), int64(42), uint(42), uint8(42), uint16(42), uint32(42), uint64(42)} {
		require.Equal(t, int32(42), CreatePositionPageToken(value).Value())
	}
	for _, value := range []any{nil, "42", 42.0} {
		require.Equal(t, int32(0), CreatePositionPageToken(value).Value())
	}
}
