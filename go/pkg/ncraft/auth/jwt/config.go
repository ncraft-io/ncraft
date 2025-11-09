package jwt

import "github.com/mojo-lang/mojo/go/pkg/mojo/core"

type Config struct {
	Enable          bool           `json:"enable,omitempty"`
	Algorithm       string         `json:"algorithm,omitempty"`
	Key             string         `json:"key,omitempty"`
	ExpiredDuration *core.Duration `json:"expiredDuration,omitempty"` // 100s, 10m,  24h
	IgnoreTime      bool           `json:"ignoreTime,omitempty"`
	Exceptions      []string       `json:"exceptions,omitempty"`
}
