package main

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestInitParams(t *testing.T) {
	params := InitParams()
	assert.IsType(t, params, &paramsConfig{})
}

func TestInitConfig(t *testing.T) {
	os.Clearenv()
	params := &paramsConfig{}
	cfg := InitConfig(params)
	assert.IsType(t, cfg, &config{})
	assert.NotEqual(t, cfg.ServerAddress, "")
	assert.NotEqual(t, cfg.BaseURL, "")
}

func TestSetEnvFromJSON(t *testing.T) {
	os.Clearenv()
	err := setEnvFromJSON("defaults/defaults_test.json")
	assert.Nil(t, err)
	params := &paramsConfig{}
	cfg := InitConfig(params)
	assert.Equal(t, "test", cfg.FileStoragePath)
	assert.Equal(t, "test", cfg.ServerAddress)
	assert.Equal(t, "test", cfg.BaseURL)
	assert.Equal(t, "test", cfg.DataBaseDSN)
	assert.Equal(t, false, cfg.EnableHTTPS)
	os.Clearenv()
}
