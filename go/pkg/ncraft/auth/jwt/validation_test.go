package jwt

import (
	"testing"
	"time"

	jwtlib "github.com/golang-jwt/jwt/v5"
	"github.com/mojo-lang/mojo/go/pkg/mojo/core"
	"github.com/stretchr/testify/require"
)

func TestValidateTimeChecksEveryClaim(t *testing.T) {
	past := float64(time.Now().Add(-time.Hour).Unix())
	future := float64(time.Now().Add(time.Hour).Unix())
	for _, tc := range []struct {
		name  string
		token Token
		valid bool
	}{
		{"empty", Token{}, true},
		{"valid", Token{"exp": future, "nbf": past, "iat": past}, true},
		{"expired", Token{"exp": past}, false},
		{"not yet valid", Token{"exp": future, "nbf": future}, false},
		{"issued in future", Token{"exp": future, "iat": future}, false},
		{"not before masks issue", Token{"nbf": past, "iat": future}, false},
		{"invalid exp", Token{"exp": "bad"}, false},
		{"invalid nbf", Token{"exp": future, "nbf": "bad"}, false},
		{"invalid iat", Token{"exp": future, "iat": "bad"}, false},
	} {
		t.Run(tc.name, func(t *testing.T) { require.Equal(t, tc.valid, tc.token.ValidateTime(nil)) })
	}
	require.False(t, Token{"exp": future, "iat": past}.ValidateTime(&core.Duration{Seconds: 60}))
}

func TestJWTSignatureAndClaims(t *testing.T) {
	j := &JWT{Key: []byte("test-key")}
	token := Token{}.SetSubject("user-1").SetIssuer("issuer").SetAudience([]string{"api"}).SetIssuedAt().SetExpirationTime(time.Hour)
	signed, err := j.Generate(token)
	require.NoError(t, err)
	parsed, err := j.Parse(signed)
	require.NoError(t, err)
	subject, err := parsed.GetSubject()
	require.NoError(t, err)
	require.Equal(t, "user-1", subject)
	audience, err := parsed.GetAudience()
	require.NoError(t, err)
	require.Equal(t, []string{"api"}, audience)
	require.Equal(t, "issuer", parsed.Normalize()["issuer"])
	_, err = (&JWT{Key: []byte("wrong-key")}).Parse(signed)
	require.Error(t, err)
	for _, method := range []jwtlib.SigningMethod{jwtlib.SigningMethodHS384, jwtlib.SigningMethodHS512} {
		signed, err = jwtlib.NewWithClaims(method, jwtlib.MapClaims(token)).SignedString(j.Key)
		require.NoError(t, err)
		_, err = j.Parse(signed)
		require.Error(t, err)
	}
	_, err = j.Parse("malformed")
	require.Error(t, err)
}

func TestJWTIgnoreTimeStillVerifiesSignature(t *testing.T) {
	j := &JWT{Key: []byte("test-key")}
	signed, err := j.Generate(Token{}.SetExpirationTime(-time.Hour))
	require.NoError(t, err)
	_, err = j.Parse(signed)
	require.Error(t, err)
	j.Config = &Config{IgnoreTime: true}
	_, err = j.Parse(signed)
	require.NoError(t, err)
	j.Key = []byte("wrong")
	_, err = j.Parse(signed)
	require.Error(t, err)
}
