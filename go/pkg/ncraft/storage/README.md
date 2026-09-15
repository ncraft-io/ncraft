# Object storage

The `minio` implementation uses `github.com/minio/minio-go/v7` v7.3.0 and registers
both `minio` and `s3` vendors. Blank-import the implementation when using the
storage factory:

```go
import (
    "context"
    "github.com/ncraft-io/ncraft/go/pkg/ncraft/storage"
    _ "github.com/ncraft-io/ncraft/go/pkg/ncraft/storage/minio"
)

createBucket := false
backend, err := storage.NewStorageWithContext(ctx, &storage.Config{
    Vendor: "s3",
    Endpoint: "https://s3.example.com",
    BucketName: "files",
    Region: "us-east-1",
    Prefix: "uploads",
    CreateBucket: &createBucket,
})
```

Omitting `AccessKey` and `SecretKey` reads `AWS_ACCESS_KEY_ID`,
`AWS_SECRET_ACCESS_KEY`, and optional `AWS_SESSION_TOKEN`. Explicit credentials
must include both key and secret. Endpoint URLs determine TLS; bare host:port
endpoints use `Secure` (false by default for compatibility). Bucket addressing
uses path style. The backend adds `Prefix` to object keys for every operation.

`CreateBucket: false` performs no bucket existence/creation calls at startup.
Omitting it keeps the previous behavior of checking and creating the bucket.
Use an explicit region when bucket region discovery is not permitted.

## Cancellable reads and writes

The factory returns configuration and initialization errors. S3 implements
`StreamingStorage` in addition to the existing `Storage` interface:

```go
streaming := backend.(storage.StreamingStorage)
reader, metadata, err := streaming.Open(ctx, "manual.pdf")
if err != nil {
    return err
}
defer reader.Close()
// reader implements io.Reader, io.Seeker and io.Closer.
// metadata includes size, MIME type, ETag and last modification time.
```

`Open` resolves metadata and missing-object errors before returning. The reader
fetches bytes on demand and supports HTTP Range requests through seeking.
`WriteContext` uploads an `Object` with its MIME type, using the actual content
length rather than a caller-supplied `Size` value.

The concrete MinIO implementation also offers `ReadContext`, `DownloadContext`
and `UploadContext`. The existing methods and constructors remain available;
they use `context.Background()`. Prefer the context variants for request handlers
and give initialization a deadline when bucket provisioning is enabled.

`NoSuchKey` is mapped to `core.NewNotFoundError`; permission, connection and
bucket configuration failures retain their original errors. Concurrent callers
reuse the SDK client. Do not change buckets during a multi-operation workflow;
construct separate backends for separate buckets.

## Validation

```sh
go test ./pkg/ncraft/storage/...
go test -race ./pkg/ncraft/storage/...
```

Tests use the real SDK against a local HTTP protocol fixture. They check signed
requests, session tokens, region, bucket provisioning, object prefixes, special
characters, metadata, empty objects, complete reads, seeking, cancellation and
missing/forbidden objects. No running MinIO server is required.
