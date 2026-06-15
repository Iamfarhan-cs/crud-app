package main

import (
	"log"

	"github.com/Iamfarhan-cs/crud-app/internal/config"
)

func main() {
	cfg := config.Load()

	// The API command is responsible for application wiring only.
	// It should compose config, database, repositories, services, and handlers.
	// Business rules, SQL queries, and HTTP request handling must not live here.
	log.Printf("task management API scaffold loaded for %s", cfg.Environment)
}
