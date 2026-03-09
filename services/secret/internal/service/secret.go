package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/iamsainadh/neovault/secret/internal/app_error"
	"github.com/iamsainadh/neovault/secret/internal/model"
	"github.com/iamsainadh/neovault/secret/internal/repository"
)

// SecretService is the contract the handler layer depends on.
type SecretService interface {
	CreateSecretService(ctx context.Context, input model.CreateSecretInput) (model.Secret, error)
	GetSecretByKey(ctx context.Context, orgID, key string) (model.Secret, error)
	CheckSecretRotationScheduler(ctx context.Context, cs model.CheckSecretRotationInput) (model.CheckSecretRotationResponse, error)
	UpdateSecretRotated(ctx context.Context, currentVersion int32, lastRotatedAt time.Time) (string, error)
}

// ──────────────────────────────── Implementation ─────────────────────────────

type secretService struct {
	repo           repository.SecretRepository
	versionService SecretVersionService
}

// NewSecretService returns a SecretService that delegates persistence to repo.
func NewSecretService(repo repository.SecretRepository, versionService SecretVersionService) SecretService {
	return &secretService{repo: repo, versionService: versionService}
}

func (s *secretService) CreateSecretService(ctx context.Context, input model.CreateSecretInput) (model.Secret, error) {

	if err := validateCreateInput(input); err != nil {
		return model.Secret{}, err
	}

	// Apply business defaults
	secret := model.Secret{
		OrgID:                input.OrgID,
		Name:                 strings.TrimSpace(input.Name),
		Description:          strings.TrimSpace(input.Description),
		SecretType:           normalizeSecretType(input.SecretType),
		IsRotationEnabled:    input.IsRotationEnabled,
		RotationIntervalDays: input.RotationIntervalDays,
		CreatedBy:            input.CreatedBy,
		CurrentVersion:       1,
		Status:               model.SecretStatusActive,
		PayloadFormat:        derivePayloadFormat(input.Value),
	}

	if secret.IsRotationEnabled && secret.RotationIntervalDays <= 0 {
		secret.RotationIntervalDays = 30
	}

	createdSecret, err := s.repo.CreateSecretRepo(ctx, secret)

	if err != nil {
		return model.Secret{}, err
	}

	plainValue, err := extractPlainText(input.Value)
	if err != nil {
		return model.Secret{}, err
	}

	_, err = s.versionService.CreateSecretVersion(ctx, model.CreateSecretVersionInput{
		Secret_id:  createdSecret.ID,
		Created_by: input.CreatedBy,
		OrgId:      input.OrgID,
		Value:      plainValue,
	})

	if err != nil {
		fmt.Println("Error From Secret", err)
		return model.Secret{}, err
	}

	return createdSecret, nil
}

func (s *secretService) GetSecretByKey(ctx context.Context, orgID, key string) (model.Secret, error) {
	if strings.TrimSpace(orgID) == "" {
		return model.Secret{}, fmt.Errorf("%w: org_id is required", app_error.ErrInvalidInput)
	}
	if strings.TrimSpace(key) == "" {
		return model.Secret{}, fmt.Errorf("%w: key is required", app_error.ErrInvalidInput)
	}

	return s.repo.GetSecretByKey(ctx, orgID, key)
}

func (r *secretService) CheckSecretRotationScheduler(ctx context.Context, cs model.CheckSecretRotationInput) (model.CheckSecretRotationResponse, error) {
	return r.repo.CheckSecretRotationScheduler(ctx, cs)
}

func derivePayloadFormat(v model.SecretValue) string {

	if v.JSON != nil {
		return "json"
	}
	return "raw"
}

func extractPlainText(v model.SecretValue) (string, error) {
	if v.Raw != nil && *v.Raw != "" {
		return *v.Raw, nil
	}

	if v.JSON != nil && *v.JSON != "" {
		return *v.JSON, nil
	}

	return "", fmt.Errorf("secret payload is empty")
}

func validateCreateInput(input model.CreateSecretInput) error {
	if strings.TrimSpace(input.OrgID) == "" {
		return fmt.Errorf("%w: org_id is required", app_error.ErrInvalidInput)
	}
	if strings.TrimSpace(input.Name) == "" {
		return fmt.Errorf("%w: name is required", app_error.ErrInvalidInput)
	}
	if len(input.Name) > 255 {
		return fmt.Errorf("%w: name must be 255 characters or fewer", app_error.ErrInvalidInput)
	}
	if strings.TrimSpace(input.CreatedBy) == "" {
		return fmt.Errorf("%w: created_by is required", app_error.ErrInvalidInput)
	}
	return nil
}

func normalizeSecretType(t model.SecretType) model.SecretType {
	if t == model.SecretTypeGeneric {
		return model.SecretTypeGeneric
	}
	return t
}

func (r *secretService) UpdateSecretRotated(ctx context.Context, currentVersion int32, lastRotatedAt time.Time) (string, error) {
	return r.repo.UpdateSecretRotated(ctx, currentVersion, lastRotatedAt)
}
