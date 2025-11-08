package config

import (
	"github.com/ncraft-io/ncraft/go/pkg/ncraft/config/reader"
)

func NcraftGet(defaultPath string, path ...string) reader.Value {
	if len(path) == 0 {
		if val := Get("ncraft." + defaultPath); !val.Null() {
			return val
		}
		return Get(defaultPath)
	} else {
		return Get(path...)
	}
}
