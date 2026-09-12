package client

import (
	"bytes"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/ncraft-io/ncraft/go/pkg/gokit/sd"
	"strings"

	"github.com/ncraft-io/ncraft/go/pkg/ncraft/config"
	"github.com/ncraft-io/ncraft/go/pkg/ncraft/config/reader"
	"github.com/ncraft-io/ncraft/go/pkg/ncraft/logs"
)

// Config contains serializable client settings. Logger, middleware, HTTP client
// instances and gRPC option functions are supplied by the caller at runtime.
type Config struct {
	ServiceName string      `json:"serviceName" yaml:"serviceName"`
	HTTP        *HTTPConfig `json:"http" yaml:"http"`
	GRPC        *GRPCConfig `json:"grpc" yaml:"grpc"`
	SD          sd.Config   `json:"sd" yaml:"sd"`
}

type HTTPConfig struct {
	// Timeout is the total request timeout in milliseconds; zero means no limit.
	Timeout   int      `json:"timeout" yaml:"timeout"`
	Headers   []string `json:"headers" yaml:"headers"`
	Enveloped bool     `json:"enveloped" yaml:"enveloped"`
}

type GRPCConfig struct {
	// DialTimeout is the connection setup timeout in milliseconds; zero leaves
	// it to the caller's context. Transport credentials remain runtime options.
	DialTimeout  int  `json:"dialTimeout" yaml:"dialTimeout"`
	WaitForReady bool `json:"waitForReady" yaml:"waitForReady"`
	// Message size limits are bytes; zero leaves the gRPC defaults in effect.
	MaxRecvMsgSize int `json:"maxRecvMsgSize" yaml:"maxRecvMsgSize"`
	MaxSendMsgSize int `json:"maxSendMsgSize" yaml:"maxSendMsgSize"`
}

// NewConfig loads common settings, optionally overridden by one service's
// settings. Service path segments are joined with dots. On failure it logs the
// error and returns nil; use LoadConfig to handle the error directly.
func NewConfig(service ...string) *Config {
	cfg, err := LoadConfig(service...)
	if err != nil {
		logs.Warnw("failed to load client config", "service", strings.Join(service, "."), "error", err)
		return nil
	}
	return cfg
}

// LoadConfig merges settings in ascending precedence:
// NcraftGet("sd"), NcraftGet("client"), NcraftGet("client.<service>").
// Objects merge recursively; explicit scalar values, arrays and null replace
// inherited values. An absent service inherits the common settings. No config
// returns an empty Config, as NewConfig did before service overrides existed.
func LoadConfig(service ...string) (*Config, error) {
	name := strings.Join(service, ".")
	if len(service) > 0 {
		for _, segment := range strings.Split(name, ".") {
			if strings.TrimSpace(segment) == "" {
				return nil, fmt.Errorf("client service name contains an empty path segment")
			}
		}
	}

	standalone, err := readConfigObject(config.NcraftGet("sd"), "sd")
	if err != nil {
		return nil, err
	}
	common, err := readConfigObject(config.NcraftGet("client"), "client")
	if err != nil {
		return nil, err
	}
	merged := make(map[string]interface{})
	if standalone != nil {
		merged["sd"] = standalone
	}
	mergeConfigObjects(merged, clientSettings(common))
	if name != "" {
		var specific map[string]interface{}
		// Registry names often contain dots. Prefer a literal service key before
		// interpreting the name as a dotted path through nested config objects.
		if value, exists := common[name]; exists {
			if value != nil {
				var ok bool
				specific, ok = value.(map[string]interface{})
				if !ok {
					return nil, fmt.Errorf("client.%s config must be an object", name)
				}
			}
		} else {
			specific, err = readConfigObject(config.NcraftGet("client."+name), "client."+name)
			if err != nil {
				return nil, err
			}
		}
		mergeConfigObjects(merged, clientSettings(specific))
	}

	data, err := jsoniter.Marshal(merged)
	if err != nil {
		return nil, fmt.Errorf("encode merged client config: %w", err)
	}
	cfg := &Config{}
	if err := jsoniter.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("decode client config for %q: %w", name, err)
	}
	return cfg, nil
}

// Read raw values instead of scanning each layer into Config: struct scans
// apply defaults and lose the distinction between absent and explicit zero.
func readConfigObject(value reader.Value, path string) (map[string]interface{}, error) {
	var raw jsoniter.RawMessage
	if err := value.Scan(&raw); err != nil {
		return nil, fmt.Errorf("read %s config: %w", path, err)
	}
	if len(raw) == 0 {
		return nil, nil
	}
	decoder := jsoniter.NewDecoder(bytes.NewReader(raw))
	decoder.UseNumber()
	var object map[string]interface{}
	if err := decoder.Decode(&object); err != nil {
		return nil, fmt.Errorf("read %s config (expected an object): %w", path, err)
	}
	return object, nil
}

// Only common setting keys participate in the merge; sibling services must
// never leak into the selected client. Accept the old flattened SD fields as
// well, with explicit nested "sd" fields taking precedence within each layer.
func clientSettings(object map[string]interface{}) map[string]interface{} {
	settings := make(map[string]interface{})
	legacySD := make(map[string]interface{})
	for _, key := range []string{"mode", "transport", "url", "retry", "etcd", "nacos", "direct"} {
		if value, exists := object[key]; exists {
			legacySD[key] = value
		}
	}
	if len(legacySD) > 0 {
		settings["sd"] = legacySD
	}
	for _, key := range []string{"sd", "serviceName", "headers", "enveloped", "http", "grpc"} {
		if value, exists := object[key]; exists {
			mergeConfigObjects(settings, map[string]interface{}{key: value})
		}
	}
	return settings
}

func mergeConfigObjects(dst, src map[string]interface{}) {
	for key, value := range src {
		if object, ok := value.(map[string]interface{}); ok {
			child, ok := dst[key].(map[string]interface{})
			if !ok {
				child = make(map[string]interface{})
			}
			mergeConfigObjects(child, object)
			dst[key] = child
		} else {
			dst[key] = value
		}
	}
}
