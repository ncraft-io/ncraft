package box

import (
	"github.com/ncraft-io/ncraft/go/pkg/ncraft/config/secrets"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDecryptRejectsTruncatedAndTamperedMessages(t *testing.T) {
	s := NewSecrets(secrets.PublicKey(make([]byte, 32)), secrets.PrivateKey(make([]byte, 32)))
	require.NoError(t, s.Init())
	for length := 0; length < 40; length++ {
		require.NotPanics(t, func() {
			_, err := s.Decrypt(make([]byte, length), secrets.SenderPublicKey(make([]byte, 32)))
			require.Error(t, err)
		})
	}
	_, err := s.Decrypt(make([]byte, 64), secrets.SenderPublicKey(make([]byte, 32)))
	require.Error(t, err)
}
