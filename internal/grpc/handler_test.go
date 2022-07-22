package grpc

import (
	"github.com/stretchr/testify/assert"
	"testing"

	servicePkg "github.com/AnnV0lokitina/short-url-service/internal/mocked_service"
)

func TestNewHandler(t *testing.T) {
	service := servicePkg.NewMockedService()
	h := NewHandler(service)
	assert.IsType(t, &Handler{}, h)
}

func TestGetUserID(t *testing.T) {
	var example uint32
	userID, err := getUserID(1)
	assert.Nil(t, err)
	assert.IsType(t, example, userID)
	assert.Equal(t, uint32(1), userID)
	userID, err = getUserID(0)
	assert.Nil(t, err)
	assert.IsType(t, example, userID)
}
