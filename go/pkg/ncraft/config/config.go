// Package config is an interface for dynamic configuration.
package config

import (
	"context"
	"strings"

	"github.com/ncraft-io/ncraft/go/pkg/ncraft/config/loader"
	"github.com/ncraft-io/ncraft/go/pkg/ncraft/config/reader"
	"github.com/ncraft-io/ncraft/go/pkg/ncraft/config/source"
	"github.com/ncraft-io/ncraft/go/pkg/ncraft/config/source/file"
	"github.com/ncraft-io/ncraft/go/pkg/ncraft/config/source/flag"
)

// Config is an interface abstraction for dynamic configuration
type Config interface {
	// Values provide the reader.Values interface
	reader.Values
	// Init the config
	Init(opts ...Option) error
	// Options in the config
	Options() Options
	// Close Stop the config loader/watcher
	Close() error
	// Load config sources
	Load(source ...source.Source) error
	// Sync Force a source change-set sync
	Sync() error
	// Watch a value for changes
	Watch(path ...string) (Watcher, error)
}

// Watcher is the config watcher
type Watcher interface {
	Next() (reader.Value, error)
	Stop() error
}

type Options struct {
	Loader loader.Loader
	Reader reader.Reader
	Source []source.Source

	// for alternative data
	Context context.Context
}

type Option func(o *Options)

var (
	// Default Config Manager
	DefaultConfig, _ = NewConfig()
)

// NewConfig returns new config
func NewConfig(opts ...Option) (Config, error) {
	return newConfig(opts...)
}

// Bytes Return config as raw json
func Bytes() []byte {
	return DefaultConfig.Bytes()
}

// Map Return config as a map
func Map() map[string]interface{} {
	return DefaultConfig.Map()
}

// Scan values to a go type
func Scan(v interface{}) error {
	return DefaultConfig.Scan(v)
}

// ScanFrom scan from the specifier keys to a go type
func ScanFrom(v interface{}, key string, alternatives ...string) error {
	val := Get(key)
	for _, alter := range alternatives {
		if !val.Null() {
			break
		}

		val = Get(alter)
	}
	return val.Scan(v)
}

// Sync Force a source change set sync
func Sync() error {
	return DefaultConfig.Sync()
}

// Get a value from the config
func Get(path ...string) reader.Value {
	return DefaultConfig.Get(normalizePath(path...)...)
}

// Deprecated: Use Get instead.
func GetValue(path ...string) reader.Value {
	return Get(path...)
}

// Load config sources
func Load(sources ...source.Source) error {
	return DefaultConfig.Load(sources...)
}

// Watch a value for changes
func Watch(path ...string) (Watcher, error) {
	return DefaultConfig.Watch(normalizePath(path...)...)
}

// LoadFile is shorthand for creating a file source and loading it
func LoadFile(path string) error {
	return Load(file.NewSource(
		file.WithPath(path),
	))
}

func LoadPath(path string) error {
	return Load(newFileSources(path, "")...)
}

func LoadPathWithSuffix(path string, suffix string) error {
	return Load(newFileSources(path, suffix)...)
}

func LoadFlag() error {
	return Load(flag.NewSource())
}

func normalizePath(path ...string) []string {
	var segments []string
	for _, p := range path {
		s := strings.Split(p, ".")
		segments = append(segments, s...)
	}
	return segments
}
