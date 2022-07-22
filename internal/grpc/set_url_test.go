package grpc

import (
	"context"
	"github.com/AnnV0lokitina/short-url-service/internal/entity"
	repoPkg "github.com/AnnV0lokitina/short-url-service/internal/mocked_repo"
	servicePkg "github.com/AnnV0lokitina/short-url-service/internal/mocked_service"
	pb "github.com/AnnV0lokitina/short-url-service/proto"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSetURL(t *testing.T) {
	service := servicePkg.NewMockedService()
	service.SetBaseURL("baseURL")
	h := NewHandler(service)
	ctx := context.TODO()

	in := &pb.SetAddrRequest{
		User:         repoPkg.RightUser,
		OriginalAddr: "original",
	}

	url := entity.NewURL("original", "baseURL")
	out := &pb.SetAddrResponse{
		User:   repoPkg.RightUser,
		Result: url.Short,
	}

	result, err := h.SetURL(ctx, in)
	assert.Nil(t, err)
	assert.Equal(t, out, result)

	in.User = 0
	_, err = h.SetURL(ctx, in)
	assert.Nil(t, err)

	in = &pb.SetAddrRequest{
		User:         repoPkg.RightUser,
		OriginalAddr: "original#*%",
	}
	_, err = h.SetURL(ctx, in)
	assert.NotNil(t, err)

	in = &pb.SetAddrRequest{
		User:         repoPkg.RightUser,
		OriginalAddr: "conflict",
	}
	_, err = h.SetURL(ctx, in)
	assert.NotNil(t, err)

	in = &pb.SetAddrRequest{
		User:         repoPkg.RightUser,
		OriginalAddr: "internal error",
	}
	_, err = h.SetURL(ctx, in)
	assert.NotNil(t, err)

}
