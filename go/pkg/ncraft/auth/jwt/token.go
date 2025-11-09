package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mojo-lang/mojo/go/pkg/mojo/core"

	"github.com/ncraft-io/ncraft/go/pkg/ncraft/logs"
)

type Token map[string]any

func (t Token) Normalize() Token {
	normalized := make(map[string]any)
	for k, v := range t {
		switch k {
		case "exp":
			dt, _ := t.GetExpirationTime()
			normalized["expiration_time"] = dt
			normalized["exp"] = dt
		case "nbf":
			dt, _ := t.GetNotBefore()
			normalized["not_before"] = dt
			normalized["nbf"] = dt
		case "iat":
			dt, _ := t.GetIssuedAt()
			normalized["issued_at"] = dt
			normalized["iat"] = dt
		case "aud":
			normalized["audience"] = v
			normalized["aud"] = v
		case "iss":
			normalized["issuer"] = v
			normalized["iss"] = v
		case "sub":
			normalized["subject"] = v
			normalized["sub"] = v
		default:
			normalized[k] = v
		}
	}
	return normalized
}

func (t Token) ValidateTime(expiredDuration *core.Duration) bool {
	now := time.Now()
	if dt, err := t.GetExpirationTime(); err != nil {
		return false
	} else if dt != nil {
		return now.Before(*dt)
	}

	if dt, err := t.GetNotBefore(); err != nil {
		return false
	} else if dt != nil {
		return now.After(*dt)
	}

	if dt, err := t.GetIssuedAt(); err != nil {
		return false
	} else if dt != nil {
		if now.Before(*dt) {
			return false
		}
		logs.Debugw("the issue time", "issue", dt.String())

		if expiredDuration != nil {
			duration := expiredDuration.ToDuration()
			logs.Debugw("the issue time", "duration", expiredDuration.Format())
			return now.Before(dt.Add(duration))
		}
	}

	return true
}

func (t Token) GetExpirationTime() (*time.Time, error) {
	date, err := jwt.MapClaims(t).GetExpirationTime()
	if err != nil {
		return nil, err
	}
	if date != nil {
		return &date.Time, nil
	}
	return nil, nil
}

// GetNotBefore implements the Claims interface.
func (t Token) GetNotBefore() (*time.Time, error) {
	date, err := jwt.MapClaims(t).GetNotBefore()
	if err != nil {
		return nil, err
	}
	if date != nil {
		return &date.Time, nil
	}
	return nil, nil
}

// GetIssuedAt implements the Claims interface.
func (t Token) GetIssuedAt() (*time.Time, error) {
	date, err := jwt.MapClaims(t).GetIssuedAt()
	if err != nil {
		return nil, err
	}
	if date != nil {
		return &date.Time, nil
	}
	return nil, nil
}

// GetAudience implements the Claims interface.
func (t Token) GetAudience() ([]string, error) {
	audience, err := jwt.MapClaims(t).GetAudience()
	if err != nil {
		return nil, err
	}
	return audience, nil
}

// GetIssuer implements the Claims interface.
func (t Token) GetIssuer() (string, error) {
	return jwt.MapClaims(t).GetIssuer()
}

// GetSubject implements the Claims interface.
func (t Token) GetSubject() (string, error) {
	return jwt.MapClaims(t).GetSubject()
}

func (t Token) Get(key string) any {
	if v, ok := t[key]; ok {
		return v
	}
	return nil
}

func (t Token) Set(key string, value any) Token {
	if len(key) > 0 && value != nil {
		t[key] = value
	}
	return t
}

func (t Token) SetExpirationTime(duration time.Duration) Token {
	t["exp"] = float64(time.Now().Add(duration).Unix())
	return t
}

func (t Token) SetAudience(audience []string) Token {
	t["aud"] = audience
	return t
}

func (t Token) SetIssuedAt() Token {
	t["iat"] = float64(time.Now().Unix())
	return t
}

func (t Token) SetIssuer(issuer string) Token {
	t["iss"] = issuer
	return t
}

func (t Token) SetSubject(subject string) Token {
	t["sub"] = subject
	return t
}
