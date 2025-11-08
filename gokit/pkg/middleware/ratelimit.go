package middleware

import (
    "context"
    "github.com/go-kit/kit/endpoint"
    "github.com/mojo-lang/core/go/pkg/mojo/core"
    "golang.org/x/time/rate"
    "time"
)

func NewTokenBucketLimitMiddleware(bkt *rate.Limiter) endpoint.Middleware {
    return func(next endpoint.Endpoint) endpoint.Endpoint {
        return func(ctx context.Context, request interface{}) (response interface{}, err error) {
            if !bkt.Allow() {
                return nil, core.NewResourceExhaustedError("Rate limit exceed!")
            }
            return next(ctx, request)
        }
    }
}

func EveryRateLimiter(interval time.Duration, b int) endpoint.Middleware {
    limiter := rate.NewLimiter(rate.Every(interval), b)
    return NewTokenBucketLimitMiddleware(limiter)
}
