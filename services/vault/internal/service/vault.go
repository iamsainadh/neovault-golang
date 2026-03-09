package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/iamsainadh/neovault/vault/internal/app_error"
	"github.com/iamsainadh/neovault/vault/internal/model"
	"github.com/iamsainadh/neovault/vault/internal/repository"
)

type VaultService interface {
	CreateVaultKeyService(ctx context.Context, input model.CreateVaultInputs) (model.VaultKey, error)
	GetVaultByOrgId(ctx context.Context, input model.GetVaultByOrgInputs) (model.GetVaultByOrgResponse, error)
}

type vaultService struct {
	repo repository.VaultRepository
}

func NewVaultService(repo repository.VaultRepository) VaultService {
	return &vaultService{repo: repo}
}

func (s *vaultService) CreateVaultKeyService(ctx context.Context, input model.CreateVaultInputs) (model.VaultKey, error) {

	if err := validateCreateInput(input); err != nil {
		return model.VaultKey{}, err
	}

	vaultModel := model.VaultKey{
		OrgId:        input.OrgId,
		KeyVersion:   input.KeyVersion,
		EncryptedKey: input.EncryptedKey,
		Algorithm:    input.Algorithm,
	}

	return s.repo.CreateVaultKeyRepo(ctx, vaultModel)
}

func (s *vaultService) GetVaultByOrgId(ctx context.Context, input model.GetVaultByOrgInputs) (model.GetVaultByOrgResponse, error) {
	if input.OrgId == "" {
		return model.GetVaultByOrgResponse{}, errors.New("Organization Id is required")
	}
	return s.repo.GetVaultByOrgId(ctx, input)
}

func validateCreateInput(input model.CreateVaultInputs) error {
	if strings.TrimSpace(input.OrgId) == "" {
		return fmt.Errorf("%w: org_id is required", app_error.ErrInvalidInput)
	}

	if len(input.EncryptedKey) == 0 {
		return fmt.Errorf("%w: encrypted_key is required", app_error.ErrInvalidInput)
	}

	if strings.TrimSpace(input.Algorithm) == "" {
		return fmt.Errorf("%w: algorithm is required", app_error.ErrInvalidInput)
	}

	return nil
}

type CreateVaultInputs struct {
	OrgId        string
	KeyVersion   int
	EncryptedKey int
	Algorithm    string
}
