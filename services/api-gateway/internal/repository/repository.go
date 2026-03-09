package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	GetUserRole(ctx context.Context, userID string, orgId string) (string, error)
}

type repository struct {
	Db *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) Repository {
	return &repository{
		Db: pool,
	}
}

func (r *repository) GetUserRole(ctx context.Context, userID string, orgId string) (string, error) {
	currCtx, cancel := context.WithTimeout(ctx, time.Second*10)

	defer cancel()

	query := `SELECT role from organization_members where user_id=$1 and org_id=$2`

	var role string

	if err := r.Db.QueryRow(currCtx, query, userID, orgId).Scan(&role); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", err
		}
		return "", fmt.Errorf("get user role: %w", err)
	}

	return role, nil
}
