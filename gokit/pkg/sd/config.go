package sd

import (
	"github.com/ncraft-io/ncraft/go/pkg/ncraft/config"
	"strings"

	"github.com/ncraft-io/ncraft/go/pkg/ncraft/logs"

	"github.com/ncraft-io/ncraft/gokit/pkg/retry"
	"github.com/ncraft-io/ncraft/gokit/pkg/sd/direct"
	"github.com/ncraft-io/ncraft/gokit/pkg/sd/etcdv3"
	"github.com/ncraft-io/ncraft/gokit/pkg/sd/nacos"
)

type Config struct {
	// sd mode, like etcd, nacos, direct
	Mode      string                    `json:"mode" yaml:"mode" db:"mode"`
	Transport string                    `json:"transport" yaml:"transport" db:"transport"` // http, or grpc
	Url       string                    `json:"url" yaml:"url"`
	Retry     *retry.Config             `json:"retry" yaml:"retry" db:"retry"`
	EtcdV3    *etcdv3.Config            `json:"etcd" yaml:"etcd"`
	Nacos     *nacos.Config             `json:"nacos" yaml:"nacos"`
	Direct    map[string]*direct.Config `json:"direct" yaml:"direct" db:"direct"`
}

func NewConfig(path ...string) *Config {
	cfg := &Config{}
	if err := config.NcraftGet("sd").Scan(cfg); err != nil {
		logs.Errorw("failed to get the ncraft.sd config from "+strings.Join(path, "."), "error", err)
		return nil
	}
	return cfg
}
