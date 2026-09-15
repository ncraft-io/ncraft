# Service lifecycle

Mojo-generated standalone and unified servers recognize three optional interfaces:

```go
type Starter interface {
    Start(Config) error
}

type Shutdowner interface {
    Shutdown(context.Context) error
}

type ErrorSource interface {
    Errors() <-chan error
}
```

A service can implement any subset. The same instance is used for lifecycle
hooks, custom HTTP routes and RPC endpoints. Lifecycle hooks run on the original
service, before service/endpoint middleware wrapping.

## Startup

The host calls `Start` before binding HTTP, gRPC and debug listeners. In a unified
server, services start in registration order. When a `Start` fails, that method
must clean up its partially acquired resources; the host shuts down previously
started services in reverse order. Services without `Start` are also eligible
for `Shutdown` once reached in the startup sequence.

The host binds all transport listeners before serving requests. A bind failure
closes listeners already bound and shuts down successfully started services.

## Runtime errors

`ErrorSource` is for fatal failures of service-owned background tasks. Ordinary
request errors are returned through RPC/HTTP methods.

The service creates and owns its error channel. A buffered channel (at least one
slot for the first failure) allows reporting during startup without waiting for
the host. Subsequent reports should use a nonblocking send or select on the
service stop signal, since the host stops receiving after the first fatal error.
`Errors` must return the same channel for the service lifetime. The host
never closes it. Nil channels, nil errors and channel closure are ignored.

`WaitForError` waits for transport errors, service errors or context cancellation
without spawning forwarding goroutines. Cancellation, including SIGINT/SIGTERM
handled by the generated host, is normal shutdown. A non-nil error is logged as a
failure before the same shutdown sequence runs.

## Shutdown

The generated host drains HTTP, gRPC and debug servers concurrently, then calls
service `Shutdown` hooks in reverse startup order. Transport draining and service
cleanup share a fresh 30-second deadline (`DefaultShutdownTimeout`); the canceled
signal context is not reused for cleanup. On timeout HTTP connections are closed
and gRPC is force-stopped. Hooks must respect their context, including a deadline
already exhausted while draining requests. Cleanup errors are logged and do not
prevent remaining hooks from running.

Move business cleanup from legacy `handlers.InterruptHandler` functions into
`Shutdown`. Generated hosts now manage signals directly so error-triggered exits
and signal-triggered exits follow the same cleanup path. `Run` retains its void
signature and logs errors; application-specific process exit codes are not part
of these interfaces.

The file service does not need an extra Range server or a background error
channel: Range is handled by the main HTTP listener managed by the host.
