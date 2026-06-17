package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

type PostgresConfig struct {
	DatabaseURL           string
	MaxOpenConnections    int
	MaxIdleConnections    int
	ConnectionMaxLifetime time.Duration
}

func OpenPostgres(ctx context.Context, cfg PostgresConfig) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("open postgres: %w", err)
	}

	db.SetMaxOpenConns(cfg.MaxOpenConnections)
	db.SetMaxIdleConns(cfg.MaxIdleConnections)
	db.SetConnMaxLifetime(cfg.ConnectionMaxLifetime)

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping postgres: %w", err)
	}

	return db, nil
}
