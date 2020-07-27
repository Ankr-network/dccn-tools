package main

import (
	"context"
	greeter "github.com/Ankr-network/dccn-tools/protoc-gen-ankr/examples/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net"
)

//go:generate protoc --ankr_out=plugins=ankr:. pb/greeter.proto

type server struct {
	greeter.UnimplementedGreeterServer
}

func (*server) Hello(_ context.Context, _ *greeter.Request) (*greeter.Response, error) {
	return nil, status.Errorf(codes.Internal, "test error")
}
func (*server) World(_ context.Context, req *greeter.Request) (*greeter.Response, error) {
	return &greeter.Response{Msg: req.GetName()}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":9090")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	greeter.RegisterGreeterServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
