package grpc

import (
	"context"
	servicePkg "github.com/AnnV0lokitina/short-url-service/internal/mocked_service"
	pb "github.com/AnnV0lokitina/short-url-service/proto"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
	"testing"
)

func TestGetStats(t *testing.T) {
	service := servicePkg.NewMockedService()
	service.SetBaseURL("baseURL")

	h := NewHandler(service)
	ctx := context.TODO()

	in := &pb.StatRequest{}
	out := &pb.StatsResponse{
		NAddr:  1,
		NUsers: 1,
	}
	_, err := h.GetStats(ctx, in)
	assert.NotNil(t, err)

	md := metadata.New(map[string]string{headerIP: "101.101.101.1"})
	ctx = metadata.NewIncomingContext(ctx, md)

	result, err := h.GetStats(ctx, in)
	assert.Nil(t, err)
	assert.Equal(t, out, result)
}
