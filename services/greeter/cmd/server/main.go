package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/harmeetsingh/grpc-envoy/services/greeter/internal/adapter/grpcserver/pb/greeter/v1"
	grpcserver "github.com/harmeetsingh/grpc-envoy/services/greeter/internal/adapter/grpcserver"
	"github.com/harmeetsingh/grpc-envoy/services/greeter/internal/adapter/memory"
	"github.com/harmeetsingh/grpc-envoy/services/greeter/internal/usecase"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "50052"
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	clock := memory.NewRealClock()
	uc := usecase.New(clock)
	srv := grpcserver.New(uc)

	gs := grpc.NewServer()
	pb.RegisterGreeterServiceServer(gs, srv)
	reflection.Register(gs)

	log.Printf("greeter-service listening on :%s", port)
	if err := gs.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
