package main

import (
	context "context"
	"fmt"
	"net"

	pb "gateway/test/grpc"

	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	pb.UnimplementedPingPongServer
}

func (s *server) Ping(ctx context.Context, in *pb.PingRequest) (*pb.PingResponse, error) {
	fmt.Printf("Received request: %s\n", in.GetMessage())
	return &pb.PingResponse{Message: "pong"}, nil
}

func main() {
	addr := "127.0.0.1:50051"
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Printf("Failed to listen on %s: %v\n", addr, err)
		return
	}

	grpcServer := grpc.NewServer()
	pb.RegisterPingPongServer(grpcServer, &server{})
	reflection.Register(grpcServer)

	fmt.Printf("gRPC server listening on %s\n", addr)
	if err := grpcServer.Serve(lis); err != nil {
		fmt.Printf("Failed to serve gRPC server: %v\n", err)
	}
}
