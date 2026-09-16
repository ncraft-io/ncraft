# HTTP response bindings

Services with a fixed protobuf return type can provide a separate HTTP writer
without changing their RPC interface. Generated Mojo Go-kit transports install
`RequestToContext` using `ServerBefore` and resolve `BoundResponseWriter` before
the existing `ResponseWriter` interface and JSON/envelope encoding.

```go
if request, ok := nhttp.RequestFromContext(ctx); ok {
    // Select HTTP behavior using request.Method, headers, etc.
    // result is the non-nil protobuf response pointer returned by the service.
    if err := nhttp.BindResponseWriter(ctx, result, streamWriter); err != nil {
        return nil, err
    }
}
return result, nil
```

## Contract

- Context initialization creates an independent attachment for each request.
  Derived contexts share that attachment; a new `RequestToContext` call resets it.
  Merely deriving a context inside a handler cannot pass a new value back to its
  caller; the transport must install the attachment before calling the endpoint.
- The response must be a non-nil pointer. Binding uses pointer identity, not value
  equality. Copies and replacement responses do not inherit a writer. The latest
  binding wins, allowing synchronous endpoint retries. Bind before returning from
  the service; attachment locking does not make the response or writer thread-safe.
- The request is read-only. Use the current endpoint/encoder context for
  cancellation and tracing; a stored request may carry an earlier context.
- Writers acquire resources only inside `WriteHttpResponse`, then close them
  before returning. No cleanup callback runs when middleware discards a response.
- Endpoint errors use the normal error encoder. An attached writer is used only
  after the endpoint succeeds and returns the bound pointer. Errors before writing
  a response can be returned; after committing headers, log transfer failures
  instead of returning an error that would append a JSON error response.
- Bindings are request-local transport state. Cached/cloned protobuf values do
  not carry them. Middleware that caches responses must explicitly account for
  streaming methods. RPC callers without an HTTP context keep normal behavior.
- Endpoint metrics finish before response encoding. Use an HTTP finalizer or
  outer HTTP middleware to measure complete transfer duration and byte counts.

The runtime does not hold a global response registry or take ownership of service
resources. HTTP writers may use `http.ServeContent` for seekable downloads, while
leaving CORS and preflight handling to the hosting HTTP middleware.
