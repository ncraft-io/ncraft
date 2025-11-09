package tracing

import (
	"github.com/ncraft-io/ncraft/go/pkg/ncraft/config"
	"strings"

	"github.com/ncraft-io/ncraft/go/pkg/ncraft/logs"
)

type Config struct {
	Enable bool    `json:"enable" yaml:"Enable" default:"false"`
	Url    string  `json:"url" yaml:"url" default:"localhost:6831"`
	Param  float64 `json:"param" json:"param" default:"100000"`
}

func NewConfig(path ...string) *Config {
	cfg := &Config{}
	if err := config.NcraftGet("tracing").Scan(cfg); err != nil {
		logs.Errorw("failed to get the ncraft.tracing config from "+strings.Join(path, "."), "error", err)
		return nil
	}
	return cfg
}
