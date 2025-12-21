package grpc

import (
	"context"
	user "messenger-project/internal/repository/user-service"
	pb "messenger-project/proto/user"
)

type UserServer struct {
	pb.UnimplementedUserServiceServer
}

func (s *UserServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	user, err := user.GetUserByUsername(req.Username)
	if err != nil {
		return &pb.GetUserResponse{Exists: false}, nil
	}
	return &pb.GetUserResponse{
		Id:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Exists:   true,
	}, nil
}

func NewUserServer() *UserServer {
	return &UserServer{}
}
