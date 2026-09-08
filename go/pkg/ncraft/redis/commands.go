package redis

import (
	"context"
	"errors"
	"time"
)

func Get(redis Redis, key string) (string, error) {
	if redis != nil {
		return String(redis.Do(context.Background(), "GET", key))
	}
	return "", errors.New("the redis handler is nil")
}

func Set(redis Redis, key string, value string) (string, error) {
	if redis != nil {
		_, err := redis.Do(context.Background(), "SET", key, value)
		return "", err
	}
	return "", errors.New("the redis handler is nil")
}

func SetEx(redis Redis, key string, value string, duration time.Duration) (string, error) {
	if redis != nil {
		_, err := redis.Do(context.Background(), "SETEX", key, int(duration.Seconds()), value)
		return "", err
	}
	return "", errors.New("the redis handler is nil")
}

func Del(redis Redis, keys ...string) (interface{}, error) {
	if redis != nil {
		return redis.Do(context.Background(), "DEL", keys)
	}
	return 0, errors.New("the redis handler is nil")
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
