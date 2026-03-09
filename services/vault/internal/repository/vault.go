package repository

import (
	"context"
	"errors"
	"time"

	"github.com/iamsainadh/neovault/vault/internal/app_error"
	"github.com/iamsainadh/neovault/vault/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type VaultRepository interface {
	CreateVaultKeyRepo(ctx context.Context, v model.VaultKey) (model.VaultKey, error)
	GetVaultByOrgId(ctx context.Context, input model.GetVaultByOrgInputs) (model.GetVaultByOrgResponse, error)
}

type vaultRepository struct {
	db *pgxpool.Pool
}

func NewVaultRepository(pool *pgxpool.Pool) VaultRepository {
	return &vaultRepository{db: pool}
}

func (r *vaultRepository) CreateVaultKeyRepo(ctx context.Context, v model.VaultKey) (model.VaultKey, error) {
	currCxt, cancel := context.WithTimeout(ctx, 10*time.Second)

	defer cancel()

	const query = `INSERT INTO VAULT_KEYS (org_id,encrypted_key,algorithm) VALUES($1, $2, $3) RETURNING id,org_id,key_version`

	var vaultStatus = ""

	err := r.db.QueryRow(currCxt, query, v.OrgId, v.EncryptedKey, v.Algorithm).Scan(
		&v.ID,
		&v.OrgId,
		&v.KeyVersion,
	)

	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" &&
				pgErr.ConstraintName == "uq_vault_org_id_key_version" {

				return model.VaultKey{}, app_error.ErrDuplicateVaultKey
			}
		}
		return model.VaultKey{}, err
	}

	v.Status = vaultStringToStatus(vaultStatus)

	return v, nil
}

func (r *vaultRepository) GetVaultByOrgId(ctx context.Context, input model.GetVaultByOrgInputs) (model.GetVaultByOrgResponse, error) {
	currCxt, cancel := context.WithTimeout(ctx, 10*time.Second)

	defer cancel()

	query := `SELECT id,key_version,encrypted_key from vault_keys where org_id =$1  AND status = 'active' AND  deleted_at IS NULL ORDER BY key_version DESC LIMIT 1;`

	var res model.GetVaultByOrgResponse

	err := r.db.QueryRow(currCxt, query, input.OrgId).Scan(
		&res.Id,
		&res.Version, &res.Encrypted_key,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.GetVaultByOrgResponse{}, app_error.ErrOrgNotFound
		}
		return model.GetVaultByOrgResponse{}, err
	}

	return res, nil
}

// func vaultStatusToString(s model.VaultKeyStatus) string {
// 	switch s {
// 	case model.Active:
// 		return "active"
// 	case model.Rotated:
// 		return "archived"
// 	case model.Retired:
// 		return "retired"
// 	case model.Compromised:
// 		return "compromised"
// 	default:
// 		return "active"
// 	}
// }

func vaultStringToStatus(val string) model.VaultKeyStatus {
	switch val {
	case "active":
		return model.Active
	case "archived":
		return model.Rotated
	case "deleted":
		return model.Retired
	case "compromised":
		return model.Compromised
	default:
		return "active"
	}
}
