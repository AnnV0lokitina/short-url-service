package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInitConfig(t *testing.T) {
	cfg := initConfig()
	assert.IsType(t, cfg, &config{})
	assert.NotEqual(t, cfg.ServerAddress, "")
	assert.NotEqual(t, cfg.BaseURL, "")

	initParams(cfg)
}
