package service

import (
	"context"
	"errors"
	repoPkg "github.com/AnnV0lokitina/short-url-service/internal/mocked_repo"
	labelError "github.com/AnnV0lokitina/short-url-service/pkg/error"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetStats(t *testing.T) {
	repo := new(repoPkg.MockedRepo)
	s := NewService("baseURL", repo, "101.101.101")
	_, err := s.GetStats(context.TODO(), "")
	assert.NotNil(t, err)
	var labelErr *labelError.LabelError
	assert.True(t, errors.As(err, &labelErr))
	assert.Equal(t, labelError.TypeForbidden, labelErr.Label)

	_, err = s.GetStats(context.TODO(), "111.101.101.2")
	assert.NotNil(t, err)
	assert.True(t, errors.As(err, &labelErr))
	assert.Equal(t, labelError.TypeForbidden, labelErr.Label)

	stats, err := s.GetStats(context.TODO(), "101.101.101.2")
	assert.Nil(t, err)
	assert.Equal(t, 1, stats.URLs)
	assert.Equal(t, 1, stats.Users)
}

func TestGetIPSubnet(t *testing.T) {
	repo := new(repoPkg.MockedRepo)
	s := NewService("baseURL", repo, "")
	sub, err := s.getIPSubnet("101.101.101.2")
	assert.Nil(t, err)
	assert.Equal(t, "101.101.101", sub)
}
