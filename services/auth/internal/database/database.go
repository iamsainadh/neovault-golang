package database

import (
	"context"
	"fmt"
	"time"

	"github.com/iamsainadh/neovault/auth/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewConnectionPool(cfg *config.Config) (*pgxpool.Pool, error) {
	currCtx, cancel := context.WithTimeout(context.Background(), time.Second*10)

	defer cancel()

	pool, err := pgxpool.New(currCtx, cfg.DB_URI)

	if err != nil {
		return nil, fmt.Errorf("unable to create db pool: %w", err)

	}

	if err := pool.Ping(currCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	fmt.Println("Connected to PostgreSQL!")

	return pool, nil
}
