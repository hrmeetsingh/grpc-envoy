package grpcserver

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/harmeetsingh/grpc-envoy/services/greeter/internal/adapter/grpcserver/pb/greeter/v1"
	"github.com/harmeetsingh/grpc-envoy/services/greeter/internal/domain"
	"github.com/harmeetsingh/grpc-envoy/services/greeter/internal/usecase"
)

type Server struct {
	pb.UnimplementedGreeterServiceServer
	uc *usecase.GreeterUseCase
}

func New(uc *usecase.GreeterUseCase) *Server {
	return &Server{uc: uc}
}

func (s *Server) SayHello(ctx context.Context, req *pb.SayHelloRequest) (*pb.SayHelloResponse, error) {
	g, err := s.uc.SayHello(ctx, req.GetName())
	if err != nil {
		return nil, toStatus(err)
	}
	return &pb.SayHelloResponse{Message: g.Message}, nil
}

func (s *Server) SayGoodbye(ctx context.Context, req *pb.SayGoodbyeRequest) (*pb.SayGoodbyeResponse, error) {
	g, err := s.uc.SayGoodbye(ctx, req.GetName())
	if err != nil {
		return nil, toStatus(err)
	}
	return &pb.SayGoodbyeResponse{Message: g.Message}, nil
}

func toStatus(err error) error {
	switch err {
	case domain.ErrEmptyName:
		return status.Error(codes.InvalidArgument, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}
