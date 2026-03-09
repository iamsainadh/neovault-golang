package database

import (
	"context"
	"fmt"
	"time"

	"github.com/iamsainadh/neovault/secret/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	maxConns          = 20
	minConns          = 2
	maxConnLifetime   = 30 * time.Minute
	maxConnIdleTime   = 5 * time.Minute
	healthCheckPeriod = 1 * time.Minute
	connectTimeout    = 10 * time.Second
)

//! New Connection Pool

func PostgresConnectionPool(cfg *config.Config) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)

	defer cancel()

	pool, err := pgxpool.New(ctx, cfg.DB_URI)


	if err != nil {
		return nil, fmt.Errorf("unable to create db pool: %w", err)

	}

	//! Ping Db
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	fmt.Println("Connected to PostgreSQL!")

	return pool, nil
}
