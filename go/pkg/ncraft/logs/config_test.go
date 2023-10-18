package logs

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ncraft-io/ncraft/go/pkg/ncraft/config"
)

func TestConfig(t *testing.T) {
	err := config.LoadPath("./testdata/configs")
	assert.NoError(t, err)

	logCfg := &Config{}
	err = config.ScanFrom(logCfg, "ncraft.logs", "logs")
	assert.NoError(t, err)

	assert.NotNil(t, logCfg.File)
	assert.True(t, logCfg.File.Compress)
}
