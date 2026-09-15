// Package minio implements MinIO and S3-compatible object storage.
package minio

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"
	"strings"
	"sync"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/minio/minio-go/v7/pkg/s3utils"
	"github.com/mojo-lang/mojo/go/pkg/mojo/core"
	"github.com/ncraft-io/ncraft/go/pkg/ncraft/logs"
	"github.com/ncraft-io/ncraft/go/pkg/ncraft/storage"
)

func init() {
	constructor := func(ctx context.Context, cfg *storage.Config) (storage.Storage, error) { return New(ctx, cfg) }
	storage.RegisterWithContext("minio", constructor)
	storage.RegisterWithContext("s3", constructor)
}

var _ storage.StreamingStorage = (*Minio)(nil)

type Minio struct {
	client     *minio.Client
	mu         sync.RWMutex
	bucketName string
	region     string
	prefix     string
}

// NewMinio preserves the original constructor for existing callers.
func NewMinio(cfg *storage.Config) storage.Storage {
	backend, err := New(context.Background(), cfg)
	if err != nil {
		logs.Errorw("failed to initialize minio storage", "error", err)
		return nil
	}
	return backend
}

// New initializes a reusable S3 client and returns configuration/bucket errors.
func New(ctx context.Context, cfg *storage.Config) (*Minio, error) {
	if cfg == nil {
		return nil, errors.New("storage configuration is required")
	}
	if err := s3utils.CheckValidBucketNameStrict(cfg.BucketName); err != nil {
		return nil, fmt.Errorf("invalid storage bucketName: %w", err)
	}
	endpoint := strings.TrimSpace(cfg.Endpoint)
	if endpoint == "" {
		return nil, errors.New("storage endpoint is required")
	}
	secure := cfg.Secure != nil && *cfg.Secure
	if strings.Contains(endpoint, "://") {
		u, err := url.Parse(endpoint)
		if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" || u.User != nil || (u.Path != "" && u.Path != "/") || u.RawQuery != "" || u.Fragment != "" {
			return nil, errors.New("storage endpoint must be an HTTP(S) origin without credentials, path or query")
		}
		endpoint, secure = u.Host, u.Scheme == "https"
	}
	prefix := strings.Trim(cfg.Prefix, "/")
	if prefix != "" {
		prefix += "/"
	}
	if (cfg.AccessKey == "") != (cfg.SecretKey == "") {
		return nil, errors.New("storage accessKey and secretKey must be configured together")
	}
	creds := credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, cfg.SessionToken)
	if cfg.AccessKey == "" {
		if cfg.SessionToken != "" {
			return nil, errors.New("storage sessionToken requires accessKey and secretKey")
		}
		creds = credentials.NewEnvAWS()
	}
	value, err := creds.Get()
	if err != nil {
		return nil, err
	}
	if value.AccessKeyID == "" || value.SecretAccessKey == "" {
		return nil, errors.New("S3 credentials are required in config or AWS_ACCESS_KEY_ID/AWS_SECRET_ACCESS_KEY")
	}
	client, err := minio.New(endpoint, &minio.Options{Creds: creds, Secure: secure, Region: cfg.Region, BucketLookup: minio.BucketLookupPath})
	if err != nil {
		return nil, fmt.Errorf("initialize S3 client: %w", err)
	}
	m := &Minio{client: client, region: cfg.Region, prefix: prefix, bucketName: cfg.BucketName}
	if cfg.CreateBucket == nil || *cfg.CreateBucket {
		if err := m.ensureBucket(ctx, cfg.BucketName); err != nil {
			return nil, err
		}
	}
	return m, nil
}

func (m *Minio) BucketName() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.bucketName
}

func (m *Minio) ensureBucket(ctx context.Context, name string) error {
	exists, err := m.client.BucketExists(ctx, name)
	if err != nil {
		return fmt.Errorf("check S3 bucket: %w", err)
	}
	if !exists {
		if err = m.client.MakeBucket(ctx, name, minio.MakeBucketOptions{Region: m.region}); err != nil && minio.ToErrorResponse(err).Code != "BucketAlreadyOwnedByYou" {
			return fmt.Errorf("create S3 bucket: %w", err)
		}
	}
	return nil
}

func (m *Minio) SetBucket(name string) error {
	if name == "" || name == m.BucketName() {
		return nil
	}
	if err := s3utils.CheckValidBucketNameStrict(name); err != nil {
		return err
	}
	if err := m.ensureBucket(context.Background(), name); err != nil {
		return err
	}
	m.mu.Lock()
	m.bucketName = name
	m.mu.Unlock()
	return nil
}

func (m *Minio) key(key string) (string, error) {
	if key == "" {
		return "", core.NewInvalidArgumentError("object key is required")
	}
	key = m.prefix + key
	return key, s3utils.CheckValidObjectName(key)
}

// Open checks metadata first because GetObject defers network errors until use.
func (m *Minio) Open(ctx context.Context, key string) (io.ReadSeekCloser, *storage.Object, error) {
	objectKey, err := m.key(key)
	if err != nil {
		return nil, nil, err
	}
	obj, err := m.client.GetObject(ctx, m.BucketName(), objectKey, minio.GetObjectOptions{})
	if err != nil {
		return nil, nil, objectError(err, key)
	}
	info, err := obj.Stat()
	if err != nil {
		obj.Close()
		return nil, nil, objectError(err, key)
	}
	metadata := &storage.Object{Key: key, Etag: info.ETag, Size: info.Size, LastModified: core.FromTime(info.LastModified)}
	if info.ContentType != "" {
		metadata.ContentType, _ = core.ParseMediaType(info.ContentType)
	}
	return obj, metadata, nil
}

func objectError(err error, key string) error {
	if minio.ToErrorResponse(err).Code == "NoSuchKey" {
		return core.NewNotFoundError("object %s not found", key)
	}
	return err
}

func (m *Minio) Read(key string, options core.Options) (*storage.Object, error) {
	return m.ReadContext(context.Background(), key, options)
}

func (m *Minio) ReadContext(ctx context.Context, key string, options core.Options) (*storage.Object, error) {
	reader, metadata, err := m.Open(ctx, key)
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	metadata.Content, err = io.ReadAll(reader)
	if err != nil {
		return nil, objectError(err, key)
	}
	if int64(len(metadata.Content)) != metadata.Size {
		return nil, io.ErrUnexpectedEOF
	}
	return metadata, nil
}

func (m *Minio) Write(object *storage.Object, options core.Options) error {
	return m.WriteContext(context.Background(), object, options)
}

func (m *Minio) WriteContext(ctx context.Context, object *storage.Object, options core.Options) error {
	if object == nil {
		return core.NewInvalidArgumentError("object is required")
	}
	key, err := m.key(object.Key)
	if err != nil {
		return err
	}
	opts := minio.PutObjectOptions{}
	if object.ContentType != nil {
		opts.ContentType = object.ContentType.Format()
	}
	// Derive the length from the bytes; stale caller metadata must not truncate uploads.
	_, err = m.client.PutObject(ctx, m.BucketName(), key, bytes.NewReader(object.Content), int64(len(object.Content)), opts)
	return err
}

func (m *Minio) Download(key string, path string, options core.Options) error {
	return m.DownloadContext(context.Background(), key, path, options)
}

func (m *Minio) DownloadContext(ctx context.Context, key, path string, options core.Options) error {
	objectKey, err := m.key(key)
	if err != nil {
		return err
	}
	return objectError(m.client.FGetObject(ctx, m.BucketName(), objectKey, path, minio.GetObjectOptions{}), key)
}

func (m *Minio) Upload(localFile string, key string, options core.Options) error {
	return m.UploadContext(context.Background(), localFile, key, options)
}

func (m *Minio) UploadContext(ctx context.Context, localFile, key string, options core.Options) error {
	objectKey, err := m.key(key)
	if err != nil {
		return err
	}
	_, err = m.client.FPutObject(ctx, m.BucketName(), objectKey, localFile, minio.PutObjectOptions{})
	return err
}
