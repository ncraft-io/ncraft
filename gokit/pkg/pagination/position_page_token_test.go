package pagination

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPositionPageToken_Create(t *testing.T) {
	token := (&PositionPageToken{}).Create(100).Format()
	_ = token

	assert.NotEmpty(t, token)

	token2 := &PositionPageToken{}
	err := token2.Parse(token)
	assert.NoError(t, err)
	assert.Equal(t, int32(100), token2.Position)
}
