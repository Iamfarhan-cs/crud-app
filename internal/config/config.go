package config

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

// Load returns placeholder configuration for the architecture phase.
// Real environment parsing and validation will be added later.
func Load() Config {
	return Config{
		Environment: "development",
		HTTPPort:    "8080",
	}
}
