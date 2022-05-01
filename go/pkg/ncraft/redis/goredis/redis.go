package goredis

import (
    "context"
    "errors"
    goredis "github.com/go-redis/redis/v8"
    "github.com/ncraft-io/ncraft/go/pkg/ncraft/redis"
    "time"
)

type Redis struct {
    Client goredis.UniversalClient
}

func init() {
    redis.RegisterRedis("goredis", NewRedis)
}

func NewRedis(cfg *redis.Config) redis.Redis {
    if cfg.MinIdleConnections == 0 {
        if cfg.MaxIdleConnections != 0 {
            cfg.MinIdleConnections = cfg.MaxIdleConnections
        } else {
            cfg.MinIdleConnections = 100
        }
    }

    if cfg.MaxActiveConnections == 0 {
        cfg.MaxActiveConnections = 300
    }
    if cfg.ReadTimeout == 0 {
        cfg.ReadTimeout = time.Second
    }
    if cfg.WriteTimeout == 0 {
        cfg.WriteTimeout = cfg.ReadTimeout
    }

    //if len(cfg.ClusterSlots) != 0 {
    //	option := &go_redis.ClusterOptions{
    //		ClusterSlots:       func() ([]go_redis.ClusterSlot, error) { return cfg.ClusterSlots, nil },
    //		ReadOnly:           true,
    //		RouteRandomly:      true,
    //		Password:           cfg.Password,
    //		MinIdleConns:       cfg.MinIdleConnections,
    //		PoolSize:           cfg.MaxActiveConnections,
    //		IdleTimeout:        cfg.IdleTimeout,
    //		IdleCheckFrequency: cfg.IdleCheckFrequency,
    //		ReadTimeout:        cfg.ReadTimeout,
    //		WriteTimeout:       cfg.WriteTimeout,
    //		MaxConnAge:         cfg.ConnectionTimeout,
    //		MaxRetries:         cfg.MaxRetries,
    //		MaxRetryBackoff:    cfg.MaxRetryBackoff,
    //		MinRetryBackoff:    cfg.MinRetryBackoff,
    //	}
    //	return &Redis{cmdable: go_redis.NewClusterClient(option)}
    //}

    if len(cfg.Connections) == 1 {
        option := &goredis.Options{
            Addr:               cfg.Connections[0],
            Password:           cfg.Password, // no password set
            DB:                 cfg.DbNumber, // use default DB
            MinIdleConns:       cfg.MinIdleConnections,
            PoolSize:           cfg.MaxActiveConnections,
            IdleTimeout:        cfg.IdleTimeout,
            IdleCheckFrequency: cfg.IdleCheckFrequency,
            ReadTimeout:        cfg.ReadTimeout,
            WriteTimeout:       cfg.WriteTimeout,
            MaxConnAge:         cfg.ConnectionTimeout,
            MaxRetries:         cfg.MaxRetries,
            MaxRetryBackoff:    cfg.MaxRetryBackoff,
            MinRetryBackoff:    cfg.MinRetryBackoff,
        }
        return &Redis{Client: goredis.NewClient(option)}
    } else {
        option := &goredis.ClusterOptions{
            Addrs:              cfg.Connections,
            Password:           cfg.Password,
            MinIdleConns:       cfg.MinIdleConnections,
            PoolSize:           cfg.MaxActiveConnections,
            IdleTimeout:        cfg.IdleTimeout,
            IdleCheckFrequency: cfg.IdleCheckFrequency,
            ReadTimeout:        cfg.ReadTimeout,
            WriteTimeout:       cfg.WriteTimeout,
            MaxConnAge:         cfg.ConnectionTimeout,
            MaxRetries:         cfg.MaxRetries,
            MaxRetryBackoff:    cfg.MaxRetryBackoff,
            MinRetryBackoff:    cfg.MinRetryBackoff,
        }
        return &Redis{Client: goredis.NewClusterClient(option)}
    }
}

func (r *Redis) Do(ctx context.Context, cmd string, arguments ...interface{}) (interface{}, error) {
    args := make([]interface{}, 0, len(arguments)+1)
    args = append(args, cmd)
    args = append(args, arguments...)
    return r.Client.Do(ctx, args...).Result()
}

func (r *Redis) Close() error {
    return r.Client.Close()
}

func (r *Redis) NewBatch() redis.Batch {
    return &pipeline{
        pipeliner: r.Client.Pipeline(),
        cmds:      []*goredis.Cmd{},
    }
}

func (r *Redis) RunBatch(ctx context.Context, batch redis.Batch) ([]interface{}, error) {
    p := batch.(*pipeline)
    if p != nil {
        _, err := p.pipeliner.Exec(ctx)
        ret := make([]interface{}, 0, len(p.cmds))
        for _, cmder := range p.cmds {
            ret = append(ret, cmder.Val())
        }
        return ret, err
    } else {
        return nil, errors.New("input wrong Batch type")
    }
}

func (r *Redis) Stats() *redis.Stats {
    stats := r.Client.PoolStats()
    return &redis.Stats{
        Hits:       int64(stats.Hits),
        Misses:     int64(stats.Misses),
        Timeouts:   int64(stats.Timeouts),
        TotalCount: int64(stats.TotalConns),
        IdleCount:  int64(stats.IdleConns),
        StaleCount: int64(stats.StaleConns),
    }
}

type pipeline struct {
    pipeliner goredis.Pipeliner
    cmds      []*goredis.Cmd
}

func (p *pipeline) Put(cmd string, args ...interface{}) error {
    newArgs := []interface{}{cmd}
    newArgs = append(newArgs, args...)
    cmder := p.pipeliner.Do(context.Background(), newArgs...)
    p.cmds = append(p.cmds, cmder)
    return nil
}
