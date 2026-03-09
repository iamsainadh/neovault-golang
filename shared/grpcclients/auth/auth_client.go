package client

import (
	pb "github.com/iamsainadh/neovault/auth/proto"
	"github.com/iamsainadh/neovault/utils"
)

func NewAuthClient(address string) (pb.UserServiceClient, error) {
	conn, err := utils.NewGrpcConnectionClient(address)
	if err != nil {
		return nil, err
	}
	return pb.NewUserServiceClient(conn), nil
}
