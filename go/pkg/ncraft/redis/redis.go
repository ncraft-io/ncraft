package redis

import (
	"context"
	"fmt"
	"sync"

	"github.com/ncraft-io/ncraft/go/pkg/ncraft/config"
	"github.com/ncraft-io/ncraft/go/pkg/ncraft/logs"
)

type Redis interface {
	// Do execute a redis command with random number arguments. First argument will
	// be used as key to hash to a slot, so it only supports a subset of redis
	// commands.
	//
	// SUPPORTED: most commands of keys, strings, lists, sets, sorted sets, hashes.
	// NOT SUPPORTED: scripts, transactions, clusters.
	//
	// Particularly, MSET/MSETNX/MGET are supported using result aggregation.
	// To MSET/MSETNX, there's no atomicity gurantee that given keys are set at once.
	// It's possible that some keys are set, while others not.
	//
	// See README.md for more details.
	// See full redis command list: http://www.redis.io/commands
	Do(ctx context.Context, cmd string, arguments ...interface{}) (interface{}, error)

	Close() error

	// NewBatch create a new redisBatch to pack mutiple commands.
	NewBatch() Batch

	// RunBatch execute commands in redisBatch simutaneously. If multiple commands are
	// directed to the same node, they will be merged and sent at once using pipeling.
	RunBatch(ctx context.Context, batch Batch) ([]interface{}, error)

	Stats() *Stats
}

type Batch interface {
	// Put add a redis command to redisBatch.
	Put(cmd string, arguments ...interface{}) error
}

func New(cfg *Config) Redis {
	if cfg == nil {
		cfg = &Config{}
		if err := config.ScanFrom(cfg, "ncraft.redis", "redis"); err != nil {
			logs.Warnw("not set the config and can't read from the config file, will try to use the default config")
			cfg.Connections = []string{":6379"}
		}
	}

	if len(cfg.Connections) == 0 {
		cfg.Connections = []string{":6379"}
	}

	// if cfg.MinIdleConnections == 0 && cfg.MaxIdleConnections != 0 {
	//	cfg.MinIdleConnections = cfg.MaxIdleConnections
	// }

	// if len(cfg.Connections) == 1 { // signal redis mode
	//	return NewPool(cfg)
	// } else { // cluster redis mode
	//	return NewCluster(cfg)
	// }

	if len(cfg.Implementor) == 0 {
		cfg.Implementor = "goredis"
	}

	if p, ok := plugins[cfg.Implementor]; ok {
		return p(cfg)
	}
	panic(fmt.Sprintf("the redis client implementor (%s) not found", cfg.Implementor))
}

var plugins map[string]func(*Config) Redis
var pluginsOnce sync.Once

func RegisterRedis(name string, r func(*Config) Redis) {
	pluginsOnce.Do(func() {
		plugins = make(map[string]func(*Config) Redis)
	})

	plugins[name] = r
}
