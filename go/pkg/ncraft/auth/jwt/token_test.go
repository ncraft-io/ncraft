package jwt

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestToken(t *testing.T) {
	token := make(Token)
	token.SetIssuedAt()

	iss, err := token.GetIssuedAt()
	assert.NoError(t, err)
	assert.True(t, iss.Sub(time.Now()) < time.Second)
}
