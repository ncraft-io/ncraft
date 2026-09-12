# Client configuration

`NewConfig()` reads common client settings. `NewConfig("book")` additionally
applies the settings for `client.book`. `LoadConfig` has the same behavior and
returns an error instead of logging it and returning nil.

```go
import clientconfig "github.com/ncraft-io/ncraft/go/pkg/gokit/client"

cfg, err := clientconfig.LoadConfig("book")
if err != nil {
    return err
}
// cfg.Config is the merged sd.Config, ready to pass to the generated NewClient.
// Existing access through cfg.Mode, cfg.Transport, cfg.Retry, etc. still works.
```

Configuration is read through the NCraft configuration manager. Its existing
file loading supports separate `configs/sd.yaml` and `configs/client.yaml`
files. The client loader itself reads the manager's current values, including
values loaded by the application; it does not create a separate file watcher.

The merge order, from lowest to highest precedence, is:

1. `config.NcraftGet("sd")`, used as the default `client.sd` settings.
2. `config.NcraftGet("client")`, the common client settings.
3. `config.NcraftGet("client.book")`, the selected service's settings.

`NcraftGet` prefers the corresponding `ncraft.*` section and otherwise falls
back to the root section. Nested objects merge by field. Lists replace inherited
lists, including `[]`; explicit `false`, `0`, `""`, and `null` override inherited
values. An empty object retains inherited fields; use `null` to clear an entire
subconfiguration. An absent service inherits the common settings. If none of
these sections exist, loading returns an empty `Config` without an error.

For example, `configs/sd.yaml` can provide discovery defaults:

```yaml
sd:
  mode: direct
  transport: grpc
  retry:
    enable: true
    timeout: 1000
    max: 3
  direct:
    sample.v1.BookService:
      urls: ["localhost:9000"]
```

And `configs/client.yaml` can override them:

```yaml
client:
  http:
    headers: [Authorization, X-Request-ID]
    enveloped: true
    timeout: 5000
  grpc:
    dialTimeout: 3000
    waitForReady: false
    maxRecvMsgSize: 4194304
    maxSendMsgSize: 4194304
  sd:
    retry:
      timeout: 2000

  book:
    serviceName: sample.v1.BookService
    enveloped: false
    headers: [Authorization]
    sd:
      transport: http
      retry:
        max: 5
      direct:
        sample.v1.BookService:
          urls: ["http://localhost:8080"]
```

`LoadConfig("book")` produces direct HTTP discovery with retry timeout `2000`
and maximum attempts `5`. It inherits `retry.enable`, HTTP timeout and gRPC
settings. Its headers list is only `[Authorization]`, and envelope decoding is
disabled. Loading any other service still uses the common discovery defaults.

The service selector is separate from `serviceName`, which overrides the registry
key. `LoadConfig("sample.v1.BookService")` first looks for that literal key in
`client`; otherwise it reads the nested path. Multiple arguments, such as
`LoadConfig("sample", "v1", "BookService")`, form the same service name.
The argument identifies a service under `client`, not an arbitrary config path.
Service keys share the client section with common setting keys, so use service
names distinct from the setting names below.

| Setting               | Consumer behavior |
|-----------------------| --- |
| `sd`                  | NCraft discovery configuration, available as `cfg.Config` |
| `serviceName`         | Registry key override, corresponding to `WithServiceName` |
| `http.headers`        | Context keys to send, corresponding to `CtxValuesToSend` |
| `http.enveloped`      | Decode the HTTP response's `data` field, corresponding to `WithEnvelope` |
| `http.timeout`        | `http.Client.Timeout`, in milliseconds; zero means no limit |
| `grpc.dialTimeout`    | Dial context timeout, in milliseconds; zero uses the caller's context |
| `grpc.waitForReady`   | `grpc.WaitForReady` call option |
| `grpc.maxRecvMsgSize` | `grpc.MaxCallRecvMsgSize`, in bytes; zero uses the gRPC default |
| `grpc.maxSendMsgSize` | `grpc.MaxCallSendMsgSize`, in bytes; zero uses the gRPC default |

This package loads settings; the consuming client applies them to transport
constructors and options. Logger, middleware, HTTP client instances and gRPC
option functions (including transport credentials) remain runtime arguments.

