package config

import (
	"time"

	"github.com/ncraft-io/ncraft/go/pkg/ncraft/config/reader"
)

type value struct{}

func newValue() reader.Value {
	return new(value)
}

func (v *value) Null() bool {
	return true
}

func (v *value) Bool(bool) bool {
	return false
}

func (v *value) Int(int) int {
	return 0
}

func (v *value) String(string) string {
	return ""
}

func (v *value) Float64(float64) float64 {
	return 0.0
}

func (v *value) Duration(time.Duration) time.Duration {
	return time.Duration(0)
}

func (v *value) StringSlice([]string) []string {
	return nil
}

func (v *value) StringMap(map[string]string) map[string]string {
	return map[string]string{}
}

func (v *value) Scan(interface{}) error {
	return nil
}

func (v *value) Bytes() []byte {
	return nil
}
