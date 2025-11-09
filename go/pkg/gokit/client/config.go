package client

import (
	"github.com/ncraft-io/ncraft/go/pkg/gokit/sd"
	"github.com/ncraft-io/ncraft/go/pkg/ncraft/config"
	"strings"

	"github.com/ncraft-io/ncraft/go/pkg/ncraft/logs"
)

type Config struct {
	sd.Config
}

func NewConfig(path ...string) *Config {
	cfg := &Config{}

	if err := config.NcraftGet("client").Scan(cfg); err != nil {
		logs.Warnw("failed to get the ncraft.client config from ", "path", strings.Join(path, "."), "error", err)
		return nil
	}
	return cfg
}
