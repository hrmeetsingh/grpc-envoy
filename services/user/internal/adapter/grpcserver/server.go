package grpcserver

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/harmeetsingh/grpc-envoy/services/user/internal/adapter/grpcserver/pb/user/v1"
	"github.com/harmeetsingh/grpc-envoy/services/user/internal/domain"
	"github.com/harmeetsingh/grpc-envoy/services/user/internal/usecase"
)

type Server struct {
	pb.UnimplementedUserServiceServer
	uc *usecase.UserUseCase
}

func New(uc *usecase.UserUseCase) *Server {
	return &Server{uc: uc}
}

func (s *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	u, err := s.uc.CreateUser(ctx, req.GetName())
	if err != nil {
		return nil, toStatus(err)
	}
	return &pb.CreateUserResponse{Id: u.ID, Name: u.Name}, nil
}

func (s *Server) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	u, err := s.uc.GetUser(ctx, req.GetId())
	if err != nil {
		return nil, toStatus(err)
	}
	return &pb.GetUserResponse{Id: u.ID, Name: u.Name}, nil
}

func toStatus(err error) error {
	switch err {
	case domain.ErrUserNotFound:
		return status.Error(codes.NotFound, err.Error())
	case domain.ErrEmptyName:
		return status.Error(codes.InvalidArgument, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}
