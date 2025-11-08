package server

import (
	"github.com/ncraft-io/ncraft/go/pkg/ncraft/config"
	"strings"

	"github.com/ncraft-io/ncraft/go/pkg/ncraft/logs"
)

type Perf struct {
	MaxProcess int `json:"maxProcess"`
}

type Config struct {
	DebugAddr string `json:"debugAddr" yaml:"debugAddr" default:":20170"`
	HttpAddr  string `json:"httpAddr" yaml:"httpAddr" default:":20171"`
	GrpcAddr  string `json:"grpcAddr" yaml:"grpcAddr" default:":20172"`
	TcpAddr   string `json:"tcpAddr" yaml:"tpcAddr" default:":20173"`
	UdpAddr   string `json:"udpAddr" yaml:"udpAddr" default:":20174"`
	Perf      Perf   `json:"Perf" yaml:"Perf"`
}

func NewConfig(path ...string) *Config {
	cfg := &Config{}
	err := config.NcraftGet("server").Scan(cfg)
	if err != nil {
		logs.Errorw("failed to get the ncraft.server config from "+strings.Join(path, "."), "error", err)
		return nil
	}
	return cfg
}
