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

func TestSetShortenBatch(t *testing.T) {
	service := servicePkg.NewMockedService()
	service.SetBaseURL("baseURL")
	h := NewHandler(service)
	ctx := context.TODO()

	inList := []*pb.AddrListItemRequest{
		{
			CorrelationId: "1",
			OriginalAddr:  "original1",
		},
		{
			CorrelationId: "2",
			OriginalAddr:  "original2",
		},
	}
	in := &pb.AddrListRequest{
		User: repoPkg.RightUser,
		List: inList,
	}

	outList := make([]*pb.AddrListItemResponse, 0, len(inList))
	for _, item := range inList {
		url := entity.NewURL(item.OriginalAddr, "baseURL")
		outList = append(outList, &pb.AddrListItemResponse{
			CorrelationId: item.CorrelationId,
			ShortAddr:     url.Short,
		})
	}

	out := &pb.AddrListResponse{
		User: repoPkg.RightUser,
		List: outList,
	}

	result, err := h.SetShortenBatch(ctx, in)
	assert.Nil(t, err)
	assert.Equal(t, out, result)

	in.User = repoPkg.UserWithError
	_, err = h.SetShortenBatch(ctx, in)
	assert.NotNil(t, err)
}
