package database

import (
	"context"
	"fmt"
	"time"

	"github.com/iamsainadh/neovault/scheduler/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPoolConnection(cfg *config.Config) (*pgxpool.Pool, error) {
	currCtx, cancel := context.WithTimeout(context.Background(), time.Second*10)

	defer cancel()

	pool, err := pgxpool.New(currCtx, cfg.DbUri)

	if err != nil {
		return nil, fmt.Errorf("Unable to create db pool: %w", err)
	}

	if err := pool.Ping(currCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("database ping failed: %w", err)

	}
	fmt.Println("Connected to postgreSQL!")

	return pool, nil
}
