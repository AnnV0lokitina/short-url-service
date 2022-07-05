package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInitConfig(t *testing.T) {
	cfg := initConfig()
	assert.NotEqual(t, cfg.ServerAddress, "")
	assert.NotEqual(t, cfg.BaseURL, "")

	initParams(cfg)
}
