package redis

import (
    "context"
    "errors"
)

func Get(redis Redis, key string) (string, error) {
    if redis != nil {
        return String(redis.Do(context.Background(), "GET", key))
    }
    return "", errors.New("the redis handler is nil")
}

//type SetOptions struct {
//}

func Set(redis Redis, key string, value string) (string, error) {
    if redis != nil {
        redis.Do(context.Background(), "SET", key, value)
        return "", nil
    }
    return "", errors.New("the redis handler is nil")
}

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
