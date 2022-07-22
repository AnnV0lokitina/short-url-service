package main

import (
	"github.com/stretchr/testify/assert"
	"os"
	"strconv"
	"testing"
)

func TestInitParams(t *testing.T) {
	params := InitParams()
	assert.IsType(t, params, &paramsConfig{})
}

func TestInitParamsWithConfig(t *testing.T) {
	os.Clearenv()
	str := "123"
	envStr := "456"
	b := true
	cfg1 := paramsConfig{}
	cfg2 := paramsConfig{
		ServerAddress:   &str,
		BaseURL:         &str,
		FileStoragePath: &str,
		DataBaseDSN:     &str,
		EnableHTTPS:     &b,
		Config:          &str,
		TrustedSubnet:   &str,
	}
	err := os.Setenv("SERVER_ADDRESS", envStr)
	assert.Nil(t, err)
	err = os.Setenv("BASE_URL", envStr)
	assert.Nil(t, err)
	err = os.Setenv("FILE_STORAGE_PATH", envStr)
	assert.Nil(t, err)
	err = os.Setenv("DATABASE_DSN", envStr)
	assert.Nil(t, err)
	err = os.Setenv("ENABLE_HTTPS", strconv.FormatBool(false))
	assert.Nil(t, err)
	err = os.Setenv("TRUSTED_SUBNET", envStr)
	assert.Nil(t, err)

	cfg := initParamsWithConfig(&cfg1)
	assert.Equal(t, envStr, cfg.ServerAddress)
	assert.Equal(t, envStr, cfg.BaseURL)
	assert.Equal(t, envStr, cfg.FileStoragePath)
	assert.Equal(t, envStr, cfg.DataBaseDSN)
	assert.Equal(t, false, cfg.EnableHTTPS)
	assert.Equal(t, "", cfg.Config)
	assert.Equal(t, envStr, cfg.TrustedSubnet)

	cfg = initParamsWithConfig(&cfg2)
	assert.Equal(t, str, cfg.ServerAddress)
	assert.Equal(t, str, cfg.BaseURL)
	assert.Equal(t, str, cfg.FileStoragePath)
	assert.Equal(t, str, cfg.DataBaseDSN)
	assert.Equal(t, true, cfg.EnableHTTPS)
	assert.Equal(t, str, cfg.Config)
	assert.Equal(t, str, cfg.TrustedSubnet)
}

func TestInitConfig(t *testing.T) {
	os.Clearenv()
	params := &paramsConfig{}
	cfg := InitConfig(params)
	assert.IsType(t, cfg, &config{})
	assert.NotEqual(t, cfg.ServerAddress, "")
	assert.NotEqual(t, cfg.BaseURL, "")
	assert.Equal(t, cfg.FileStoragePath, "")
	assert.Equal(t, cfg.DataBaseDSN, "")
	assert.Equal(t, cfg.TrustedSubnet, "")

	path := "defaults/defaults_test.json"
	params2 := &paramsConfig{
		Config: &path,
	}
	cfg = InitConfig(params2)
	assert.IsType(t, cfg, &config{})
	assert.Equal(t, cfg.ServerAddress, "test")
	assert.Equal(t, cfg.BaseURL, "test")
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
