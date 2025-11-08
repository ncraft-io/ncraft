package sd

import (
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	"github.com/ncraft-io/ncraft/gokit/pkg/sd/direct"
	"github.com/ncraft-io/ncraft/gokit/pkg/sd/etcdv3"
	"github.com/ncraft-io/ncraft/gokit/pkg/sd/nacos"
	"strings"
)

type Client interface {
	// Register our instance.
	Register(url, service string, tags []string) error

	// Deregister At the end of our service lifecycle, for example at the end of func main,
	// we should make sure to deregister ourselves. This is important! Don't
	// accidentally skip this step by invoking a log.Fatal or os.Exit in the
	// interim, which bypasses the defer stack.
	Deregister() error

	// Instancer It's likely that we'll also want to connect to other services and call
	// their methods. We can build an Instancer to listen for changes from sd,
	// create Endpointer, wrap it with a load-balancer to pick a single
	// endpoint, and finally wrap it with a retry strategy to get something that
	// can be used as an endpoint directly.
	Instancer(service string) sd.Instancer
}

const (
	EtcdMode   = "etcd"
	DirectMode = "direct"
	NacosMode  = "nacos"
)

func New(cfg *Config, logger log.Logger) Client {
	if cfg == nil || len(cfg.Url) == 0 {
		return nil
	}

	mode := strings.ToLower(cfg.Mode)
	urls := strings.Split(cfg.Url, ";")
	if cfg.EtcdV3 != nil && (mode == EtcdMode || len(mode) == 0) {
		return etcdv3.New(urls, cfg.EtcdV3, logger)
	} else if mode == DirectMode {
		return direct.New(cfg.Direct)
	} else if cfg.Nacos != nil && mode == NacosMode {
		return nacos.NewClient(urls, cfg.Nacos, logger)
	}
	return nil
}
