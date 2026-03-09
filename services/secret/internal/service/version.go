package service

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/iamsainadh/neovault/secret/internal/app_error"
	"github.com/iamsainadh/neovault/secret/internal/model"
	"github.com/iamsainadh/neovault/secret/internal/repository"
	"github.com/iamsainadh/neovault/utils"
	vault "github.com/iamsainadh/neovault/vault/proto"
)

type SecretVersionService interface {
	CreateSecretVersion(ctx context.Context, input model.CreateSecretVersionInput) (model.CreateSecretVersionResponse, error)
	GetSecretVersionBySecretId(ctx context.Context, secretId string) (model.GetSecretVersionBySecretIdResponse, error)
}

type secretVersionService struct {
	repo        repository.SecretVersionRepository
	vaultClient vault.VaultServiceClient
}

func NewSecretVersionService(repo repository.SecretVersionRepository, vaultClient vault.VaultServiceClient) SecretVersionService {
	return &secretVersionService{
		repo:        repo,
		vaultClient: vaultClient,
	}
}

func (s *secretVersionService) CreateSecretVersion(ctx context.Context, input model.CreateSecretVersionInput) (model.CreateSecretVersionResponse, error) {

	if err := validateCreateVersionInput(input); err != nil {
		return model.CreateSecretVersionResponse{}, err
	}

	//! Get The Vault Info based on the organization
	createdVault, err := s.vaultClient.GetVaultByOrgId(ctx, &vault.GetVaultRequest{
		OrgId: input.OrgId,
	})

	if err != nil {
		slog.Warn("Failed To Get Key:", "Vault Key Error", err)
		return model.CreateSecretVersionResponse{}, err
	}

	// GENERATE DEK
	dek, err := utils.GenerateKey()

	if err != nil {
		slog.Warn("Failed To Generate The DEK", "DEK Error", err)
		return model.CreateSecretVersionResponse{}, err
	}

	// Encrypt The Plain Text with DEK
	encryptedPayload, err := utils.Encrypt([]byte(input.Value), dek)

	if err != nil {
		slog.Warn("Failed To Encrypt", "Encryption Error", err)
		return model.CreateSecretVersionResponse{}, err
	}

	rmk := utils.ConvertToAESKeyFromRMK(os.Getenv("RMK_KEY"))

	// Decrypt The KEK
	decryptKek, err := utils.DecryptKEK(createdVault.EncryptedKey, rmk)

	if err != nil {
		slog.Warn("Failed To Encrypt", "Encryption Error", err)
		return model.CreateSecretVersionResponse{}, err
	}

	encryptedDek, err := utils.Encrypt(dek, decryptKek)

	//! Create Secret Version
	createSecretVersion := model.CreateSecretVersionInput{
		Secret_id:         input.Secret_id,
		Version:           createdVault.Version,
		Encrypted_payload: encryptedPayload,
		Encrypted_dek:     encryptedDek,
		Vault_key_id:      createdVault.Id,
		Created_by:        input.Created_by,
	}

	return s.repo.CreateSecretVersion(ctx, createSecretVersion)
}

func validateCreateVersionInput(input model.CreateSecretVersionInput) error {
	if strings.TrimSpace(input.Secret_id) == "" {
		return fmt.Errorf("%w: secret_id is required", app_error.ErrInvalidInput)
	}
	return nil
}

func (s *secretVersionService) GetSecretVersionBySecretId(ctx context.Context, secretId string) (model.GetSecretVersionBySecretIdResponse, error) {
	return s.repo.GetSecretVersionBySecretId(ctx, secretId)
}
