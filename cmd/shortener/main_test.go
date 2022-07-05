package main

import (
	"context"
	repoPkg "github.com/AnnV0lokitina/short-url-service/internal/repo"
	servicePkg "github.com/AnnV0lokitina/short-url-service/internal/service"
	"github.com/AnnV0lokitina/short-url-service/internal/sqlrepo"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestInitRepo(t *testing.T) {
	cfg := &config{}
	repo, err := initRepo(context.TODO(), cfg)
	assert.Nil(t, err)
	assert.IsType(t, repo, repoPkg.NewMemoryRepo())

	os.Clearenv()
	err = setEnvFromJSON("defaults/defaults_run_test.json")
	assert.Nil(t, err)
	cfgEnv := InitConfig()

	cfg.FileStoragePath = cfgEnv.FileStoragePath
	repo, err = initRepo(context.TODO(), cfg)
	assert.Nil(t, err)
	fileRepo, err := repoPkg.NewFileRepo(cfg.FileStoragePath)
	assert.Nil(t, err)
	assert.IsType(t, repo, fileRepo)

	cfg = &config{}
	cfg.DataBaseDSN = cfgEnv.DataBaseDSN
	repo, err = initRepo(context.TODO(), cfg)
	assert.Nil(t, err)
	sqlRepo, err := sqlrepo.NewSQLRepo(context.TODO(), cfg.DataBaseDSN)
	assert.Nil(t, err)
	assert.IsType(t, repo, sqlRepo)
	os.Clearenv()
}

func TestCreateApp(t *testing.T) {
	cfg := &config{}
	repo := repoPkg.NewMemoryRepo()

	app, service := createApp(cfg, repo)
	assert.IsType(t, app, &App{})
	assert.IsType(t, service, servicePkg.NewService("", repo))
}
