package client

import (
	pb "github.com/iamsainadh/neovault/secret/proto"
	"github.com/iamsainadh/neovault/utils"
)

func NewSecretsClient(address string) (pb.SecretServiceClient, error) {

	conn, err := utils.NewGrpcConnectionClient(address)

	if err != nil {
		return nil, err
	}

	return pb.NewSecretServiceClient(conn), nil
}

func NewSecretVersionClient(address string) (pb.SecretVersionClient, error) {
	conn, err := utils.NewGrpcConnectionClient(address)
	if err != nil {
		return nil, err
	}

	return pb.NewSecretVersionClient(conn), nil
}
