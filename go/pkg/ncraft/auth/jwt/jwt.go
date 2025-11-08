package jwt

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"

	"github.com/ncraft-io/ncraft/go/pkg/ncraft/config"
	"github.com/ncraft-io/ncraft/go/pkg/ncraft/logs"
)

type JWT struct {
	Config *Config

	Key        []byte
	Exceptions []*UrlPath
}

func NewJWT() *JWT {
	cfg := &Config{}

	if err := config.NcraftGet("jwt").Scan(cfg); err != nil {
		logs.Warnw("failed to get the ncraft.jwt config", "error", err)
		return nil
	}

	if cfg.Enable {
		j := &JWT{
			Config: cfg,
			Key:    []byte(cfg.Key),
		}

		for _, e := range cfg.Exceptions {
			j.Exceptions = append(j.Exceptions, NewUrlPath(e))
		}

		return j
	}

	return nil
}

func (j *JWT) Generate(token Token) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(token))

	// Sign a JWT!
	signed, err := t.SignedString(j.Key)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %s\n", err.Error())
	}
	return signed, nil
}

func (j *JWT) Parse(token string) (Token, error) {
	tk, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
		return j.Key, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return nil, err
	}
	if tk == nil || !tk.Valid {
		return nil, fmt.Errorf("the token (%s) is invalid", tk.Raw)
	}

	return Token(tk.Claims.(jwt.MapClaims)), nil
}

func (j *JWT) Exceptional(path string) bool {
	if j != nil {
		for _, ex := range j.Exceptions {
			if ex.Match(path) {
				return true
			}
		}
	}
	return false
}
