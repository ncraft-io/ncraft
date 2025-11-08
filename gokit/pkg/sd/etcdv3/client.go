package etcdv3

import (
    "context"
    "github.com/go-kit/kit/log"
    kitsd "github.com/go-kit/kit/sd"
    std "github.com/go-kit/kit/sd/etcdv3"
    "github.com/ncraft-io/ncraft/go/pkg/ncraft/logs"
    "net/url"
    "strings"
    "time"
)

type Client struct {
    client    std.Client
    registrar *std.Registrar
    logger    log.Logger
}

func New(urls []string, cfg *Config, logger log.Logger) *Client {
    options := std.ClientOptions{
        // Path to trusted ca file
        CACert: cfg.CACert,

        // Path to certificate
        Cert: cfg.Cert,

        // Path to private key
        Key: cfg.Key,

        // Username if required
        Username: cfg.Username,

        // Password if required
        Password: cfg.Password,

        // If DialTimeout is 0, it defaults to 3s
        DialTimeout: time.Second * time.Duration(cfg.DialTimeout),

        // If DialKeepAlive is 0, it defaults to 3s
        DialKeepAlive: time.Second * time.Duration(cfg.DialKeepAlive),
    }

    // Build the client.
    client, err := std.NewClient(context.Background(), urls, options)
    if err != nil {
        return nil
    }

    return &Client{
        client: client,
        logger: logger,
    }
}

func (c *Client) Register(urlStr, name string, tags []string) error {
    // Build the registrar.

    if !strings.HasPrefix(urlStr, "etcd://") {
        urlStr = "etcd://" + urlStr
    }

    url, err := url.Parse(urlStr)
    if err != nil {
        return err
    }

    registrar := std.NewRegistrar(c.client, std.Service{
        Key:   "/" + name + "/" + url.Host,
        Value: url.Host,
    }, c.logger)

    c.registrar = registrar
    registrar.Register()

    s, _ := c.client.GetEntries("/")
    logs.Info("all key:", s)
    return nil
}

func (c *Client) Deregister() error {
    c.registrar.Deregister()
    return nil
}

func (c *Client) Instancer(service string) kitsd.Instancer {
    if c == nil {
        return nil
    }
    instancer, err := std.NewInstancer(c.client, "/"+service+"/", c.logger)
    if err != nil {
        panic(err)
    }
    return instancer
}
