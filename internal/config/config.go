package config

import "os"

// Config contains application settings loaded from the runtime environment.
// This package should own environment and configuration loading.
// Secrets must come from environment variables or a managed secret provider,
// not from committed files. Local .env files must stay ignored by Git.
// HTTP handlers, business rules, and database queries must not live here.
type Config struct {
	Environment string
	HTTPPort    string
	PostgresDSN string
}

// Load returns minimal runtime configuration.
// Full validation can be added when configuration rules are finalized.
func Load() Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return Config{
		Environment: "development",
		HTTPPort:    port,
	}
}
