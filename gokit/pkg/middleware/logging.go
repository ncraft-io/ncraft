package middleware

import (
    "context"
    "github.com/go-kit/kit/endpoint"
    "github.com/go-kit/kit/log"
    "time"
)

// Logging returns an endpoint middleware that logs the
// duration of each invocation, and the resulting error, if any.
func Logging(logger log.Logger) endpoint.Middleware {
    return func(next endpoint.Endpoint) endpoint.Endpoint {
        return func(ctx context.Context, request interface{}) (response interface{}, err error) {
            defer func(begin time.Time) {
                logger.Log("error", err, "took", time.Since(begin))
            }(time.Now())
            return next(ctx, request)
        }
    }
}
