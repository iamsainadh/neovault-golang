package repository

import (
	"context"
	"errors"
	"time"

	"github.com/iamsainadh/neovault/secret/internal/app_error"
	"github.com/iamsainadh/neovault/secret/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SecretVersionRepository interface {
	CreateSecretVersion(ctx context.Context, input model.CreateSecretVersionInput) (model.CreateSecretVersionResponse, error)
	GetSecretVersionBySecretId(ctx context.Context, secretId string) (model.GetSecretVersionBySecretIdResponse, error)
}

type secretVersionRepository struct {
	DB *pgxpool.Pool
}

func NewSecretVersionRepository(pool *pgxpool.Pool) SecretVersionRepository {
	return &secretVersionRepository{
		DB: pool,
	}
}

func (r *secretVersionRepository) CreateSecretVersion(ctx context.Context, input model.CreateSecretVersionInput) (model.CreateSecretVersionResponse, error) {
	currCtx, cancel := context.WithTimeout(ctx, time.Second*10)

	defer cancel()

	query := `INSERT INTO SECRET_VERSIONS(secret_id,version,encrypted_payload,encrypted_dek,vault_key_id,created_by) VALUES ($1,$2,$3,$4,$5,$6) RETURNING id,secret_id,version,created_by`

	var secretVersion model.CreateSecretVersionResponse

	err := r.DB.QueryRow(currCtx, query,
		input.Secret_id,
		input.Version,
		input.Encrypted_payload,
		input.Encrypted_dek,
		input.Vault_key_id,
		input.Created_by,
	).Scan(
		&secretVersion.Id,
		&secretVersion.Secret_id,
		&secretVersion.Version,
		&secretVersion.Created_by,
	)

	if err != nil {
		// var pgErr *pgconn.PgError

		if errors.Is(err, pgx.ErrNoRows) {
			return model.CreateSecretVersionResponse{}, app_error.ErrInvalidInput
		}
		return model.CreateSecretVersionResponse{}, err
	}

	return secretVersion, nil
}

func (r *secretVersionRepository) GetSecretVersionBySecretId(ctx context.Context, secretId string) (model.GetSecretVersionBySecretIdResponse, error) {
	currCtx, cancel := context.WithTimeout(ctx, time.Second*10)

	defer cancel()

	var secretVersion model.GetSecretVersionBySecretIdResponse

	const query = `SELECT id,secret_id,version,encrypted_payload,encrypted_dek,vault_key_id from SECRET_VERSIONS where secret_id = $1`

	err := r.DB.QueryRow(
		currCtx,
		query,
		secretId,
	).Scan(
		&secretVersion.Id,
		&secretVersion.SecretId,
		&secretVersion.Version,
		&secretVersion.EncryptedPayload,
		&secretVersion.EncryptedDek,
		&secretVersion.VaultKeyId,
	)

	if err != nil {
		return model.GetSecretVersionBySecretIdResponse{}, nil
	}

	return secretVersion, nil
}
