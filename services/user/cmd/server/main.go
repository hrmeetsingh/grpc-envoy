package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/harmeetsingh/grpc-envoy/services/user/internal/adapter/grpcserver/pb/user/v1"
	grpcserver "github.com/harmeetsingh/grpc-envoy/services/user/internal/adapter/grpcserver"
	"github.com/harmeetsingh/grpc-envoy/services/user/internal/adapter/memory"
	"github.com/harmeetsingh/grpc-envoy/services/user/internal/usecase"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "50051"
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	repo := memory.New()
	uc := usecase.New(repo)
	srv := grpcserver.New(uc)

	gs := grpc.NewServer()
	pb.RegisterUserServiceServer(gs, srv)
	reflection.Register(gs)

	log.Printf("user-service listening on :%s", port)
	if err := gs.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
