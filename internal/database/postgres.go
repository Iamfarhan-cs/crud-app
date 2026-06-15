package database

// PostgresConfig contains the values needed to configure a PostgreSQL connection.
// This file should own database connection setup concerns.
// Task business rules, HTTP routing, and environment parsing must not live here.
type PostgresConfig struct {
	DSN string
}

// Postgres is a placeholder for the future PostgreSQL connection wrapper.
// A real database connection will be introduced in a later phase.
type Postgres struct {
	DSN string
}

// NewPostgres records the intended database configuration without opening a connection yet.
func NewPostgres(cfg PostgresConfig) *Postgres {
	return &Postgres{DSN: cfg.DSN}
}
