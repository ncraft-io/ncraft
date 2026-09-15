package minio

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/mojo-lang/mojo/go/pkg/mojo/core"
	"github.com/ncraft-io/ncraft/go/pkg/ncraft/storage"
)

func boolPtr(v bool) *bool { return &v }

func TestS3ReadWriteAndRanges(t *testing.T) {
	var mu sync.Mutex
	objects := map[string][]byte{}
	types := map[string]string{}
	bucketCreated := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		auth := r.Header.Get("Authorization")
		if !strings.Contains(auth, "Credential=test-key/") || !strings.Contains(auth, "/test-region/s3/aws4_request") {
			t.Error("missing Signature V4 credential/region")
		}
		if r.Header.Get("X-Amz-Security-Token") != "test-token" {
			t.Error("missing session token")
		}
		if r.URL.Path == "/test-bucket/" || r.URL.Path == "/test-bucket" {
			if r.Method == http.MethodPut {
				bucketCreated = true
				return
			}
			if !bucketCreated {
				w.WriteHeader(404)
			}
			return
		}
		if !strings.HasPrefix(r.URL.Path, "/test-bucket/tenant/") {
			t.Errorf("wrong bucket/prefix: %s", r.URL.Path)
		}
		if strings.HasSuffix(r.URL.Path, "/denied") {
			w.Header().Set("X-Amz-Error-Code", "AccessDenied")
			w.WriteHeader(403)
			fmt.Fprint(w, "<Error><Code>AccessDenied</Code><Message>denied</Message></Error>")
			return
		}
		if r.Method == http.MethodPut {
			var reader io.Reader = r.Body
			if strings.Contains(r.Header.Get("Content-Encoding"), "aws-chunked") {
				reader = httputil.NewChunkedReader(reader)
			}
			data, err := io.ReadAll(reader)
			if err != nil {
				t.Errorf("read PUT: %v", err)
			}
			objects[r.URL.Path], types[r.URL.Path] = data, r.Header.Get("Content-Type")
			w.Header().Set("ETag", `"test-etag"`)
			return
		}
		data, ok := objects[r.URL.Path]
		if !ok {
			w.Header().Set("X-Amz-Error-Code", "NoSuchKey")
			w.WriteHeader(404)
			fmt.Fprint(w, "<Error><Code>NoSuchKey</Code></Error>")
			return
		}
		w.Header().Set("ETag", `"test-etag"`)
		w.Header().Set("Content-Type", types[r.URL.Path])
		http.ServeContent(w, r, "object", time.Unix(1700000000, 0), bytes.NewReader(data))
	}))
	defer server.Close()
	cfg := &storage.Config{Vendor: "s3", Endpoint: server.URL, BucketName: "test-bucket", Region: "test-region", AccessKey: "test-key", SecretKey: "test-secret", SessionToken: "test-token", Prefix: "/tenant/"}
	backend, err := storage.NewStorageWithContext(context.Background(), cfg)
	if err != nil {
		t.Fatal(err)
	}
	m := backend.(*Minio)
	mu.Lock()
	created := bucketCreated
	mu.Unlock()
	if !created {
		t.Fatal("bucket was not created")
	}
	content := bytes.Repeat([]byte("Hello, S3!\n"), 32768)
	mediaType, _ := core.ParseMediaType("text/plain; charset=utf-8")
	obj := &storage.Object{Key: "nested/文件 #?.txt", Content: content, Size: 1, ContentType: mediaType}
	if err := m.WriteContext(context.Background(), obj, nil); err != nil {
		t.Fatal(err)
	}
	got, err := m.Read(obj.Key, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got.Content, content) || got.Size != int64(len(content)) || got.Key != obj.Key || got.ContentType.Format() != mediaType.Format() {
		t.Fatalf("incorrect object metadata or bytes: size=%d bytes=%d", got.Size, len(got.Content))
	}
	reader, metadata, err := m.Open(context.Background(), obj.Key)
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()
	if len(metadata.Content) != 0 {
		t.Fatal("Open buffered the object")
	}
	if _, err := reader.Seek(7, io.SeekStart); err != nil {
		t.Fatal(err)
	}
	part := make([]byte, 13)
	if _, err := io.ReadFull(reader, part); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(part, content[7:20]) {
		t.Fatal("incorrect ranged content")
	}
	if _, err := m.Read("missing", nil); !core.IsNotFoundError(err) {
		t.Fatalf("missing object: %v", err)
	}
	if _, err := m.Read("denied", nil); err == nil || core.IsNotFoundError(err) {
		t.Fatalf("access denied must not become not found: %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, _, err := m.Open(ctx, obj.Key); !errors.Is(err, context.Canceled) {
		t.Fatalf("cancellation: %v", err)
	}
	if err := m.WriteContext(ctx, obj, nil); !errors.Is(err, context.Canceled) {
		t.Fatalf("write cancellation: %v", err)
	}
	if err := m.Write(&storage.Object{Key: "empty"}, nil); err != nil {
		t.Fatal(err)
	}
	empty, err := m.Read("empty", nil)
	if err != nil || len(empty.GetContent()) != 0 {
		t.Fatalf("empty object: %v", err)
	}
	if err := m.Write(nil, nil); err == nil {
		t.Fatal("nil object accepted")
	}
}

func TestConfiguration(t *testing.T) {
	t.Setenv("AWS_ACCESS_KEY_ID", "env-key")
	t.Setenv("AWS_SECRET_ACCESS_KEY", "env-secret")
	t.Setenv("AWS_SESSION_TOKEN", "env-token")
	for _, vendor := range []string{"s3", "minio", ""} {
		cfg := &storage.Config{Vendor: vendor, Endpoint: "https://localhost:9000", BucketName: "test-bucket", CreateBucket: boolPtr(false)}
		backend, err := storage.NewStorageWithContext(context.Background(), cfg)
		if err != nil {
			t.Fatal(err)
		}
		if backend.BucketName() != cfg.BucketName {
			t.Fatal("bucket mismatch")
		}
	}
	for _, endpoint := range []string{"", "ftp://localhost", "http://localhost/bucket", "https://user:pass@localhost", "http://localhost?key=value"} {
		if _, err := New(context.Background(), &storage.Config{Endpoint: endpoint, BucketName: "test-bucket", CreateBucket: boolPtr(false)}); err == nil {
			t.Fatalf("accepted endpoint %q", endpoint)
		}
	}
	for _, cfg := range []*storage.Config{nil, {Endpoint: "localhost:9000", BucketName: "bad bucket"}, {Endpoint: "localhost:9000", BucketName: "test-bucket", AccessKey: "key"}} {
		if _, err := New(context.Background(), cfg); err == nil {
			t.Fatal("accepted invalid config")
		}
	}
	for _, key := range []string{"AWS_ACCESS_KEY_ID", "AWS_ACCESS_KEY", "AWS_SECRET_ACCESS_KEY", "AWS_SECRET_KEY"} {
		t.Setenv(key, "")
	}
	if _, err := New(context.Background(), &storage.Config{Endpoint: "localhost:9000", BucketName: "test-bucket", CreateBucket: boolPtr(false)}); err == nil {
		t.Fatal("missing credentials accepted")
	}
}
