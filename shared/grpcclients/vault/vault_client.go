package client

import (
	"github.com/iamsainadh/neovault/utils"
	pb "github.com/iamsainadh/neovault/vault/proto"
)

func NewVaultClient(address string) (pb.VaultServiceClient, error) {
	conn, err := utils.NewGrpcConnectionClient(address)

	if err != nil {
		return nil, err
	}

	return pb.NewVaultServiceClient(conn), nil
}
