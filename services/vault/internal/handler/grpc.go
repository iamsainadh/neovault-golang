package handler

import (
	"context"

	"github.com/iamsainadh/neovault/vault/internal/app_error"
	"github.com/iamsainadh/neovault/vault/internal/mapper"
	"github.com/iamsainadh/neovault/vault/internal/model"
	"github.com/iamsainadh/neovault/vault/internal/service"
	pb "github.com/iamsainadh/neovault/vault/proto"
)

type VaultHandler struct {
	pb.UnimplementedVaultServiceServer
	service service.VaultService
}

func NewVaultHandler(svc service.VaultService) *VaultHandler {
	return &VaultHandler{
		service: svc,
	}
}

func (h *VaultHandler) CreateVault(ctx context.Context, req *pb.CreateVaultRequest) (*pb.CreateVaultResponse, error) {
	vault := mapper.CreateVaultRequestToInput(req)

	created, err := h.service.CreateVaultKeyService(ctx, *vault)

	if err != nil {

		return nil, app_error.ToGRPCError(err)
	}

	return &pb.CreateVaultResponse{
		Id:           created.ID,
		OrgId:        created.OrgId,
		KeyVersion:   created.KeyVersion,
		EncryptedKey: created.EncryptedKey,
	}, nil
}

func (h *VaultHandler) GetVaultByOrgId(ctx context.Context, req *pb.GetVaultRequest) (*pb.GetVaultResponse, error) {
	inputs := model.GetVaultByOrgInputs{
		OrgId: req.OrgId,
	}

	res, err := h.service.GetVaultByOrgId(ctx, inputs)

	if err != nil {

		return nil, app_error.ToGRPCError(err)
	}

	return &pb.GetVaultResponse{
		Id:           res.Id,
		Version:      res.Version,
		EncryptedKey: res.Encrypted_key,
	}, nil
}
