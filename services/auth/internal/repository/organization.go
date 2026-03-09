package repository

import (
	"context"
	"time"

	"github.com/iamsainadh/neovault/auth/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrganizationRepository interface {
	CreateOrganization(ctx context.Context, input model.CreateOrganizationInput, tx pgx.Tx) (*model.CreateOrganizationResponse, error)
}

type organizationRepository struct {
	Db *pgxpool.Pool
}

func NewOrganizationRepository(pool *pgxpool.Pool) OrganizationRepository {
	return &organizationRepository{
		Db: pool,
	}
}

func (or *organizationRepository) CreateOrganization(ctx context.Context, input model.CreateOrganizationInput, tx pgx.Tx) (*model.CreateOrganizationResponse, error) {
	currCtx, cancel := context.WithTimeout(ctx, time.Second*10)

	defer cancel()

	query := `INSERT INTO ORGANIZATIONS(name) VALUES ($1) RETURNING id,name`

	err := tx.QueryRow(currCtx, query, input.Name).Scan(
		&input.Id,
		&input.Name,
	)

	if err != nil {
		return &model.CreateOrganizationResponse{}, err
	}

	return &model.CreateOrganizationResponse{
		Id:   input.Id,
		Name: input.Name,
	}, nil
}
