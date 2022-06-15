package main

import (
	"context"
	repoPkg "github.com/AnnV0lokitina/short-url-service/internal/repo"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInitRepo(t *testing.T) {
	cfg := &config{}
	repo, err := initRepo(context.TODO(), cfg)
	assert.Nil(t, err)
	assert.IsType(t, repo, repoPkg.NewMemoryRepo())
}
