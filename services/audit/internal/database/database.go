package database

import (
	"context"

	"github.com/iamsainadh/neovault/audit/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresConnection(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, error) {

	pool, err := pgxpool.New(ctx, cfg.DbUri)

	if err != nil {
		pool.Close()
	}

	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
	}

	return pool, nil

}
