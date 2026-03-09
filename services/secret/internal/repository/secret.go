package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/bytedance/gopkg/util/logger"
	"github.com/iamsainadh/neovault/secret/internal/app_error"
	"github.com/iamsainadh/neovault/secret/internal/model"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SecretRepository interface {
	CreateSecretRepo(ctx context.Context, s model.Secret) (model.Secret, error)
	GetSecretByKey(ctx context.Context, orgID, key string) (model.Secret, error)
	CheckSecretRotationScheduler(ctx context.Context, cs model.CheckSecretRotationInput) (model.CheckSecretRotationResponse, error)
	UpdateSecretRotated(ctx context.Context, currentVersion int32, lastRotatedAt time.Time) (string, error)
}

type secretRepository struct {
	db *pgxpool.Pool
}

func NewSecretRepository(pool *pgxpool.Pool) SecretRepository {
	return &secretRepository{db: pool}
}

func (r *secretRepository) CreateSecretRepo(ctx context.Context, s model.Secret) (model.Secret, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	const query = `
		INSERT INTO secrets (
			org_id,
			name,
			description,
			secret_type,
			current_version,
			is_rotation_enabled,
			rotation_interval_days,
			status,
			created_by,
			payload_format
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9,$10)
		RETURNING
			id, org_id, name, description,
			secret_type, current_version,
			is_rotation_enabled, rotation_interval_days,
			status, created_by, created_at, updated_at,payload_format`

	var (
		secretTypeStr string
		statusStr     string
	)

	err := r.db.QueryRow(ctx, query,
		s.OrgID,
		s.Name,
		s.Description,
		secretTypeToString(s.SecretType),
		s.CurrentVersion,
		s.IsRotationEnabled,
		s.RotationIntervalDays,
		secretStatusToString(s.Status),
		s.CreatedBy,
		s.PayloadFormat,
	).Scan(
		&s.ID,
		&s.OrgID,
		&s.Name,
		&s.Description,
		&secretTypeStr,
		&s.CurrentVersion,
		&s.IsRotationEnabled,
		&s.RotationIntervalDays,
		&statusStr,
		&s.CreatedBy,
		&s.CreatedAt,
		&s.UpdatedAt,
		&s.PayloadFormat,
	)

	if err != nil {
		return model.Secret{}, mapDBError(err, "CreateSecretRepo")
	}

	s.SecretType = stringToSecretType(secretTypeStr)
	s.Status = stringToSecretStatus(statusStr)

	return s, nil
}

func (r *secretRepository) GetSecretByKey(ctx context.Context, orgID, key string) (model.Secret, error) {
	currCxt, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	const query = `
		SELECT
			id, org_id, name, description,
			secret_type, current_version,
			is_rotation_enabled, rotation_interval_days,
			status, created_by, created_at, updated_at
		FROM secrets
		WHERE org_id = $1 AND name = $2 AND status = 'active'
		LIMIT 1`

	var (
		s             model.Secret
		secretTypeStr string
		statusStr     string
	)

	err := r.db.QueryRow(currCxt, query, orgID, key).Scan(
		&s.ID,
		&s.OrgID,
		&s.Name,
		&s.Description,
		&secretTypeStr,
		&s.CurrentVersion,
		&s.IsRotationEnabled,
		&s.RotationIntervalDays,
		&statusStr,
		&s.CreatedBy,
		&s.CreatedAt,
		&s.UpdatedAt,
	)

	if err != nil {
		return model.Secret{}, mapDBError(err, "GetSecretByKey")
	}

	s.SecretType = stringToSecretType(secretTypeStr)
	s.Status = stringToSecretStatus(statusStr)

	return s, nil
}

// Check Secret Required Rotation

func (r *secretRepository) CheckSecretRotationScheduler(ctx context.Context, cs model.CheckSecretRotationInput) (model.CheckSecretRotationResponse, error) {
	currCtx, cancel := context.WithTimeout(ctx, time.Second*10)

	defer cancel()

	limit := 10
	if cs.Limit > 0 {
		limit = int(cs.Limit)
	}

	query := fmt.Sprintf(`
        SELECT 
            id,
            org_id,
            name,
            current_version
        FROM secrets
        WHERE 
            is_rotation_enabled = TRUE
            AND status = 'active'
            AND deleted_at IS NULL
            AND (
                COALESCE(last_rotated_at, created_at)
                + (rotation_interval_days || ' days')::interval
                <= NOW()
            )
        ORDER BY 
            (
                NOW() - (
                    COALESCE(last_rotated_at, created_at)
                    + (rotation_interval_days || ' days')::interval
                )
            ) DESC
        LIMIT %d
    `, limit)

	rows, err := r.db.Query(currCtx, query)

	logger.Info(rows)
	if err != nil {
		return model.CheckSecretRotationResponse{}, err
	}

	defer rows.Close()

	var secrets []*model.SecretRotationItem

	for rows.Next() {
		item := &model.SecretRotationItem{}

		err := rows.Scan(
			&item.Id,
			&item.OrgID,
			&item.Name,
			&item.CurrentVersion,
		)
		if err != nil {
			return model.CheckSecretRotationResponse{}, err
		}

		secrets = append(secrets, item)
	}

	return model.CheckSecretRotationResponse{
		Secrets: secrets,
	}, nil

}

func (r *secretRepository) UpdateSecretRotated(ctx context.Context, currentVersion int32, lastRotatedAt time.Time) (string, error) {
	currCtx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	query := `UPDATE SECRETS SET current_version = $1 , last_rotated_at = $2  where id = $3`
	
	status := "Successfully Updated"

	err := r.db.QueryRow(currCtx, query, currentVersion, lastRotatedAt)

	if err != nil {
		return "", nil
	}

	return status, nil
}

func mapDBError(err error, op string) error {

	if errors.Is(err, errNoRows()) {
		return fmt.Errorf("%s: %w", op, app_error.ErrSecretNotFound)
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505": // unique_violation
			if pgErr.ConstraintName == "uq_secrets_org_name" {
				return fmt.Errorf("%s: %w", op, app_error.ErrSecretAlreadyExists)
			}
		case "23503": // foreign_key_violation
			return fmt.Errorf("%s: %w", op, app_error.ErrForeignKeyViolation)
		}
	}

	return fmt.Errorf("%s: %w", op, err)
}

func errNoRows() error {

	return pgxNoRowsErr{}
}

type pgxNoRowsErr struct{}

func (pgxNoRowsErr) Error() string { return "no rows in result set" }

// Helper

func secretTypeToString(t model.SecretType) string {
	switch t {
	case model.SecretTypeGeneric:
		return "generic"
	case model.SecretTypeAPIKey:
		return "api_key"
	case model.SecretTypePassword:
		return "password"
	case model.SecretTypeToken:
		return "token"
	default:
		return "generic"
	}
}

func secretStatusToString(s model.SecretStatus) string {
	switch s {
	case model.SecretStatusActive:
		return "active"
	case model.SecretStatusArchived:
		return "archived"
	case model.SecretStatusDeleted:
		return "deleted"
	default:
		return "active"
	}
}

func stringToSecretType(v string) model.SecretType {
	switch v {
	case "generic":
		return model.SecretTypeGeneric
	case "api_key":
		return model.SecretTypeAPIKey
	case "password":
		return model.SecretTypePassword
	case "token":
		return model.SecretTypeToken
	default:
		return model.SecretTypeGeneric
	}
}

func stringToSecretStatus(v string) model.SecretStatus {
	switch v {
	case "active":
		return model.SecretStatusActive
	case "archived":
		return model.SecretStatusArchived
	case "deleted":
		return model.SecretStatusDeleted
	default:
		return model.SecretStatusActive
	}
}
