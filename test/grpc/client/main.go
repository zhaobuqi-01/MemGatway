package main

import (
	"context"
	"fmt"

	pb "gateway/test/grpc"

	grpc "google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("127.0.0.1:50051", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("Failed to connect to server: %v\n", err)
		return
	}
	defer conn.Close()

	client := pb.NewPingPongClient(conn)

	resp, err := client.Ping(context.Background(), &pb.PingRequest{Message: "ping"})
	if err != nil {
		fmt.Printf("Failed to ping server: %v\n", err)
		return
	}

	fmt.Printf("Response: %s\n", resp.GetMessage())
}
