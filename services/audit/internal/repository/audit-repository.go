package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuditRepository interface {
	// CreateAudit(ctx context.Context)
}

type auditRepository struct {
	pool *pgxpool.Pool
}

func NewAuditRepository(pool *pgxpool.Pool) AuditRepository {
	return &auditRepository{
		pool: pool,
	}
}
