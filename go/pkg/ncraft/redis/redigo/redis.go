package redigo

import (
    "context"
    "errors"
    redigo "github.com/gomodule/redigo/redis"
    "github.com/ncraft-io/ncraft/go/pkg/ncraft/redis"
)

func init() {
    redis.RegisterRedis("redigo", NewRedis)
}

type Redis struct {
    pool *redigo.Pool
}

type Command struct {
    name string
    args []interface{}
}

type Batch struct {
    pool *redigo.Pool
    cmds []Command
}

func NewRedis(cfg *redis.Config) redis.Redis {
    if cfg.MaxIdleConnections == 0 {
        cfg.MaxIdleConnections = 100
    }

    if cfg.MaxActiveConnections == 0 {
        cfg.MaxActiveConnections = 300
    }

    //if cfg.IdleTimeout == 0 {
    //	cfg.IdleTimeout = 180 * time.Second
    //}

    dialFunc := func() (c redigo.Conn, err error) {
        c, err = redigo.Dial("tcp", cfg.Connections[0])
        if err != nil {
            panic(err)
            return nil, err
        }

        if len(cfg.Password) > 0 {
            if _, err := c.Do("AUTH", cfg.Password); err != nil {
                c.Close()
                return nil, err
            }
        }

        _, err = c.Do("SELECT", cfg.DbNumber)
        if err != nil {
            c.Close()
            return nil, err
        }
        return
    }

    // initialize a new pool
    p := &redigo.Pool{
        MaxIdle:     cfg.MaxIdleConnections,
        IdleTimeout: cfg.IdleTimeout,
        MaxActive:   cfg.MaxActiveConnections,
        Dial:        dialFunc,
    }

    return &Redis{pool: p}
}

func (p *Redis) Do(ctx context.Context, cmd string, arguments ...interface{}) (interface{}, error) {
    connection := p.pool.Get()
    defer connection.Close()

    return connection.Do(cmd, arguments...)
}

func (p *Redis) Close() error {
    return p.pool.Close()
}

func (p *Redis) Stats() *redis.Stats {
    return &redis.Stats{}
}

// NewBatch implement the Cluster NewBatch method.
func (p *Redis) NewBatch() redis.Batch {
    return &Batch{
        pool: p.pool,
        cmds: make([]Command, 0),
    }
}

// RunBatch implement the Cluster RunBatch method.
func (p *Redis) RunBatch(ctx context.Context, batch redis.Batch) ([]interface{}, error) {
    bat := batch.(*Batch)

    connection := p.pool.Get()
    defer connection.Close()

    for _, cmd := range bat.cmds {
        err := connection.Send(cmd.name, cmd.args...)
        if err != nil {
            return nil, err
        }
    }

    err := connection.Flush()
    if err != nil {
        return nil, err
    }

    var replies []interface{}
    for i := 0; i < len(bat.cmds); i++ {
        reply, err := connection.Receive()
        if err != nil {
            return nil, err
        }

        replies = append(replies, reply)
    }

    return replies, nil
}

func (b *Batch) Put(cmd string, args ...interface{}) error {
    if len(args) < 1 {
        return errors.New("no key found in args")
    }

    b.cmds = append(b.cmds, Command{name: cmd, args: args})
    return nil
}
