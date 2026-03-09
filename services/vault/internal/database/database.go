package database

import (
	"context"
	"fmt"
	"time"

	"github.com/iamsainadh/neovault/vault/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewConnection(cfg *config.Config) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)

	defer cancel()

	pool, err := pgxpool.New(ctx, cfg.DB_URI)

	if err != nil {
		return nil, fmt.Errorf("unable to create db pool: %w", err)

	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	fmt.Println("Connected to PostgreSQL!")

	return pool, nil

}
