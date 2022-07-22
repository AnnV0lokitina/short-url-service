package grpc

import (
	"context"
	repoPkg "github.com/AnnV0lokitina/short-url-service/internal/mocked_repo"
	servicePkg "github.com/AnnV0lokitina/short-url-service/internal/mocked_service"
	pb "github.com/AnnV0lokitina/short-url-service/proto"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDeleteBatch(t *testing.T) {
	service := servicePkg.NewMockedService()
	service.SetBaseURL("baseURL")
	h := NewHandler(service)
	ctx := context.TODO()

	in := &pb.DeleteBatchRequest{
		User:         repoPkg.RightUser,
		ChecksumList: []string{"checksum"},
	}
	out := &pb.TextResponse{
		Result: "OK",
	}

	result, err := h.DeleteBatch(ctx, in)
	assert.Nil(t, err)
	assert.Equal(t, out, result)

	in.User = servicePkg.WrongUser
	_, err = h.DeleteBatch(ctx, in)
	assert.NotNil(t, err)

}
