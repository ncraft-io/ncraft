package logs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLevelEnabled(t *testing.T) {
	SetLevel(ErrorLevel)
	assert.False(t, LevelEnabled(WarnLevel))
	assert.True(t, LevelEnabled(ErrorLevel))
}
