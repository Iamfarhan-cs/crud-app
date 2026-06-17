package config

import (
	"errors"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port                 string
	DatabaseURL          string
	ReadTimeout          time.Duration
	WriteTimeout         time.Duration
	IdleTimeout          time.Duration
	ShutdownTimeout      time.Duration
	DBMaxOpenConnections int
	DBMaxIdleConnections int
	DBConnectionMaxLife  time.Duration
}

func Load() (Config, error) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return Config{}, errors.New("DATABASE_URL is required")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	readTimeout, err := durationFromEnv("READ_TIMEOUT", 5*time.Second)
	if err != nil {
		return Config{}, err
	}
	writeTimeout, err := durationFromEnv("WRITE_TIMEOUT", 10*time.Second)
	if err != nil {
		return Config{}, err
	}
	idleTimeout, err := durationFromEnv("IDLE_TIMEOUT", 60*time.Second)
	if err != nil {
		return Config{}, err
	}
	shutdownTimeout, err := durationFromEnv("SHUTDOWN_TIMEOUT", 10*time.Second)
	if err != nil {
		return Config{}, err
	}
	dbConnectionMaxLife, err := durationFromEnv("DB_CONNECTION_MAX_LIFE", 30*time.Minute)
	if err != nil {
		return Config{}, err
	}

	dbMaxOpenConnections, err := intFromEnv("DB_MAX_OPEN_CONNECTIONS", 10)
	if err != nil {
		return Config{}, err
	}
	dbMaxIdleConnections, err := intFromEnv("DB_MAX_IDLE_CONNECTIONS", 5)
	if err != nil {
		return Config{}, err
	}

	return Config{
		Port:                 port,
		DatabaseURL:          databaseURL,
		ReadTimeout:          readTimeout,
		WriteTimeout:         writeTimeout,
		IdleTimeout:          idleTimeout,
		ShutdownTimeout:      shutdownTimeout,
		DBMaxOpenConnections: dbMaxOpenConnections,
		DBMaxIdleConnections: dbMaxIdleConnections,
		DBConnectionMaxLife:  dbConnectionMaxLife,
	}, nil
}

func durationFromEnv(name string, defaultValue time.Duration) (time.Duration, error) {
	value := os.Getenv(name)
	if value == "" {
		return defaultValue, nil
	}

	return time.ParseDuration(value)
}

func intFromEnv(name string, defaultValue int) (int, error) {
	value := os.Getenv(name)
	if value == "" {
		return defaultValue, nil
	}

	return strconv.Atoi(value)
}
