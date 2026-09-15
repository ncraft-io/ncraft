package storage

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sync"

	"github.com/ncraft-io/ncraft/go/pkg/ncraft/config"
	"github.com/ncraft-io/ncraft/go/pkg/ncraft/logs"

	"github.com/mojo-lang/mojo/go/pkg/mojo/core"
)

type initializer func(cfg *Config) Storage
type ContextInitializer func(context.Context, *Config) (Storage, error)

var initializers = make(map[string]ContextInitializer)
var initializersMu sync.RWMutex

// Register preserves support for backends using the original constructor.
func Register(name string, init initializer) {
	RegisterWithContext(name, func(_ context.Context, cfg *Config) (Storage, error) {
		backend := init(cfg)
		if backend == nil {
			return nil, fmt.Errorf("initialize storage vendor %q", name)
		}
		return backend, nil
	})
}

func RegisterWithContext(name string, init ContextInitializer) {
	initializersMu.Lock()
	defer initializersMu.Unlock()
	initializers[name] = init
}

// StreamingStorage adds cancellable, seekable reads without changing Storage.
// Callers must close the returned reader; metadata does not contain file bytes.
type StreamingStorage interface {
	Storage
	Open(context.Context, string) (io.ReadSeekCloser, *Object, error)
	WriteContext(context.Context, *Object, core.Options) error
}

type Storage interface {
	BucketName() string
	SetBucket(name string) error

	Read(key string, options core.Options) (*Object, error)
	Write(object *Object, options core.Options) error

	Download(key string, path string, options core.Options) error
	Upload(localFile string, key string, options core.Options) error
}

func NewStorage(cfg *Config) Storage {
	backend, err := NewStorageWithContext(context.Background(), cfg)
	if err != nil {
		return nil
	}
	return backend
}

func NewStorageWithContext(ctx context.Context, cfg *Config) (Storage, error) {
	if cfg == nil {
		return nil, errors.New("storage configuration is required")
	}
	vendor := cfg.Vendor
	if vendor == "" {
		vendor = "minio"
	}
	initializersMu.RLock()
	init := initializers[vendor]
	initializersMu.RUnlock()
	if init == nil {
		return nil, fmt.Errorf("unsupported storage vendor %q", vendor)
	}
	backend, err := init(ctx, cfg)
	if err != nil {
		return nil, err
	}
	if backend == nil {
		return nil, fmt.Errorf("storage vendor %q returned no backend", vendor)
	}
	return backend, nil
}

var storage Storage
var storageOnce sync.Once

func GetStorage() Storage {
	storageOnce.Do(func() {
		conf := &Config{}
		if err := config.NcraftGet("storage").Scan(conf); err != nil {
			logs.Warnw("failed to get the server config", "error", err.Error())
			storage = NewDummyStorage()
		} else {
			if storage = NewStorage(conf); storage == nil {
				storage = NewDummyStorage()
			}
		}
	})
	return storage
}

func Read(key string, options core.Options) (*Object, error) {
	return GetStorage().Read(key, options)
}

func Write(object *Object, options core.Options) error {
	return GetStorage().Write(object, options)
}

func Download(key string, path string, options core.Options) error {
	return GetStorage().Download(key, path, options)
}
func Upload(localFile string, key string, options core.Options) error {
	return GetStorage().Upload(localFile, key, options)
}

type DummyStorage struct {
	err error
}

func NewDummyStorage() Storage                                      { return &DummyStorage{err: errors.New("DummyStorage: not implement")} }
func (s *DummyStorage) BucketName() string                          { return "dummy" }
func (s *DummyStorage) SetBucket(string) error                      { return s.err }
func (s *DummyStorage) Read(string, core.Options) (*Object, error)  { return nil, s.err }
func (s *DummyStorage) Write(*Object, core.Options) error           { return s.err }
func (s *DummyStorage) Download(string, string, core.Options) error { return s.err }
func (s *DummyStorage) Upload(string, string, core.Options) error   { return s.err }
