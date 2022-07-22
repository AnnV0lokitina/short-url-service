package main

import (
	handler "github.com/AnnV0lokitina/short-url-service/internal/grpc"
	"google.golang.org/grpc"
)

type GRPCService struct {
	Handler *handler.Handler
	Server  *grpc.Server
}
