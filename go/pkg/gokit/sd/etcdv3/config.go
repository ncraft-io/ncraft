package etcdv3

type Config struct {
    Cert          string `json:"cert"`
    Key           string `json:"key"`
    CACert        string `json:"caCert"`
    DialTimeout   int    `json:"dialTimeout" default:"3"`
    DialKeepAlive int    `json:"dialKeepAlive" default:"3"`
    Username      string `json:"username"`
    Password      string `json:"password"`
}
