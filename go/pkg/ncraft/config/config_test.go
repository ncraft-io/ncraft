package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadPath(t *testing.T) {
	err := LoadPath("./testdata/configs")
	assert.NoError(t, err)

	v := ""
	err = ScanFrom(&v, "ncraft.foo.value")
	assert.NoError(t, err)
	assert.Equal(t, "val", v)

	err = ScanFrom(&v, "ncraft.bar.value")
	assert.NoError(t, err)
	assert.Equal(t, "val", v)
}

func TestScanFrom(t *testing.T) {
	err := LoadPath("./testdata/configs")
	assert.NoError(t, err)

	v := ""
	err = ScanFrom(&v, "foo.value", "ncraft.foo.value")
	assert.NoError(t, err)
	assert.Equal(t, "val", v)

	v = ""
	err = ScanFrom(&v, "foo.value")
	assert.Empty(t, v)
}

func TestLoadPathWithSuffix(t *testing.T) {
	err := LoadPathWithSuffix("./testdata/configs", "dev")
	assert.NoError(t, err)

	v := ""
	err = ScanFrom(&v, "ncraft.bar.value")
	assert.NoError(t, err)
	assert.Equal(t, "dev", v)
}

func TestLoadFile(t *testing.T) {
	err := LoadFile("./testdata/configs/bar.yaml")
	assert.NoError(t, err)

	v := ""
	err = ScanFrom(&v, "ncraft.bar.value")
	assert.NoError(t, err)
	assert.Equal(t, "val", v)

	val := GetValue("ncraft.bar.value")
	assert.True(t, !val.Null())
}

func TestBytes(t *testing.T) {
	err := LoadFile("./testdata/configs/bar.yaml")
	assert.NoError(t, err)

	bytes := Bytes()
	assert.NotEmpty(t, bytes)
}

func TestMap(t *testing.T) {
	err := LoadFile("./testdata/configs/bar.yaml")
	assert.NoError(t, err)

	val := Map()
	assert.NotEmpty(t, val)
	assert.NotNil(t, val["ncraft"])
}
