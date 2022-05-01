package redis

import (
    "context"
    "errors"
)

func MGet(redis Redis, keys ...string) (interface{}, error) {
    if redis != nil {
        var args []interface{}
        for _, key := range keys {
            args = append(args, key)
        }
        return redis.Do(context.Background(), "MGET", args...)
    }
    return nil, errors.New("the redis handler is nil")
}
