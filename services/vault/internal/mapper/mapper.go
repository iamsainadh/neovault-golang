package mapper

import (
	"github.com/iamsainadh/neovault/vault/internal/model"
	pb "github.com/iamsainadh/neovault/vault/proto"
)

// Proto To model

func CreateVaultRequestToInput(req *pb.CreateVaultRequest) *model.CreateVaultInputs {
	return &model.CreateVaultInputs{
		OrgId:        req.OrgId,
		EncryptedKey: req.EncryptedKey,
		Algorithm:    req.Algorithm,
	}
}
