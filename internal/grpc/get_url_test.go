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

func TestGetURL(t *testing.T) {
	service := servicePkg.NewMockedService()
	service.SetBaseURL("baseURL")
	h := NewHandler(service)
	ctx := context.TODO()
	repoPkg.TmpURL = &entity.URL{
		Short:    entity.CreateShortURL("checksum", "baseURL"),
		Original: "original",
	}

	in := &pb.GetAddrRequest{
		Checksum: "checksum",
	}
	out := &pb.GetAddrResponse{
		OriginalAddr: "original",
	}
	response, err := h.GetURL(ctx, in)
	assert.Nil(t, err)
	assert.Equal(t, out, response)

	in = &pb.GetAddrRequest{
		Checksum: "URL deleted",
	}
	_, err = h.GetURL(ctx, in)
	assert.Error(t, err)

	in = &pb.GetAddrRequest{
		Checksum: "Invalid request",
	}
	_, err = h.GetURL(ctx, in)
	assert.Error(t, err)
}
