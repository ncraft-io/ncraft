package metrics

import (
	"github.com/ncraft-io/ncraft/go/pkg/ncraft/config"
	"strings"

	"github.com/ncraft-io/ncraft/go/pkg/ncraft/logs"
)

type Config struct {
	Enable     bool   `json:"enable" default:"true"`
	Department string `json:"department"`
	Project    string `json:"project"`
}

func (c *Config) Enabled() bool {
	if c != nil {
		return c.Enable
	}
	return false
}

func NewConfig(path ...string) *Config {
	cfg := &Config{}

	if err := config.NcraftGet("metrics").Scan(cfg); err != nil {
		logs.Warnw("failed to get the ncraft.metrics config from ", "path", strings.Join(path, "."), "error", err)
		return nil
	}

	return cfg
}
