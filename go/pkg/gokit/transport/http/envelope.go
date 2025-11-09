package http

import (
	"context"
	"strings"
)

func IsEnvelopeStyle(ctx context.Context, style string) bool {
	enveloped := false
	style = strings.ToLower(style)
	if style == EnvelopStyle || style == underScoreEnvelopStyle {
		enveloped = true
	}
	if val, ok := ctx.Value(EnvelopStyle).(bool); ok {
		enveloped = val
	}
	if val, ok := ctx.Value(underScoreEnvelopStyle).(bool); ok {
		enveloped = val
	}
	return enveloped
}

func IsAIPStyle(ctx context.Context, style string) bool {
	aip := false
	style = strings.ToLower(style)
	if style == AIPStyle || style == underScoreAIPStyle {
		aip = true
	}
	if val, ok := ctx.Value(AIPStyle).(bool); ok {
		aip = val
	}
	if val, ok := ctx.Value(underScoreAIPStyle).(bool); ok {
		aip = val
	}
	return aip
}
