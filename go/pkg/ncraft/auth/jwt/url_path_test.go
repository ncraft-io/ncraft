package jwt

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUrlPathMatch(t *testing.T) {
	assert.True(t, NewUrlPath("/path/to/{param}/*").Match("/path/to/test/tail/1"))
	assert.False(t, NewUrlPath("/path/to/{param}").Match("/path/to/test/tail"))
}
