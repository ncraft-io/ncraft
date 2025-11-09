package middleware

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	"github.com/mojo-lang/mojo/go/pkg/mojo/core"
	"github.com/ncraft-io/ncraft/go/pkg/ncraft/auth/jwt"
)

// NewJWT returns an endpoint middleware that logs the
// duration of each invocation, and the resulting error, if any.
func NewJWT() endpoint.Middleware {
	j := jwt.NewJWT()
	if j == nil {
		return func(next endpoint.Endpoint) endpoint.Endpoint {
			return func(ctx context.Context, request interface{}) (response interface{}, err error) {
				return next(ctx, request)
			}
		}
	}

	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			if path, ok := ctx.Value("http-request-path").(string); ok && j.Exceptional(path) {
				return next(ctx, request)
			}

			token := jwt.GetContextToken(ctx)

			validated := false
			if len(token) > 0 {
				if t, err := j.Parse(token); err != nil {
				} else {
					m := t.Normalize()
					for k, v := range m {
						ctx = context.WithValue(ctx, fmt.Sprintf("jwt:%s", k), v)
					}
					if !j.Config.IgnoreTime {
						validated = t.ValidateTime(j.Config.ExpiredDuration)
					} else {
						validated = true
					}
				}
			}

			if !validated {
				return nil, core.NewUnauthenticatedError("auth code is invalid!")
			}

			return next(ctx, request)
		}
	}
}
