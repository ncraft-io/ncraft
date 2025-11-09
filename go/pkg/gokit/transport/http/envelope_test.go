package http

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsEnvelopeStyle(t *testing.T) {
	assert.True(t, IsEnvelopeStyle(context.Background(), "envelope"))
	assert.True(t, IsEnvelopeStyle(context.Background(), "_envelope"))
	assert.True(t, IsEnvelopeStyle(context.WithValue(context.Background(), "envelope", true), ""))
	assert.True(t, IsEnvelopeStyle(context.WithValue(context.Background(), "_envelope", true), ""))
}

func TestIsAIPStyle(t *testing.T) {
	assert.True(t, IsAIPStyle(context.Background(), "AIP"))
	assert.True(t, IsAIPStyle(context.Background(), "_aip"))
	assert.True(t, IsAIPStyle(context.WithValue(context.Background(), "aip", true), ""))
	assert.True(t, IsAIPStyle(context.WithValue(context.Background(), "_aip", true), ""))
}
