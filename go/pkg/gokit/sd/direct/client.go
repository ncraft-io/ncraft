package direct

import (
	kitsd "github.com/go-kit/kit/sd"
)

type Client struct {
	m map[string]*Config
}

func New(m map[string]*Config) *Client {
	ret := &Client{
		m: m,
	}
	return ret
}

func (c *Client) Register(urlStr, name string, tags []string) error {
	return nil
}

func (c *Client) Deregister() error {
	return nil
}

func (c *Client) Instancer(service string) kitsd.Instancer {
	if c == nil {
		return nil
	}
	cfg, ok := c.m[service]
	if !ok || cfg == nil {
		return nil
	}

	if len(cfg.Urls) == 0 {
		return nil
	}
	var ret kitsd.FixedInstancer
	return append(ret, cfg.Urls...)
}
