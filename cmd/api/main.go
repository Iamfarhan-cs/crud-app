package main

import (
	"log"
	"net/http"
	"os"

	apphttp "github.com/Iamfarhan-cs/crud-app/internal/http"
	"github.com/Iamfarhan-cs/crud-app/internal/storage"
)

func main() {
	addr := ":" + envOrDefault("PORT", "8080")
	store := storage.NewMemoryUserStore()
	handler := apphttp.NewRouter(store)

	log.Printf("server listening on %s", addr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatal(err)
	}
}

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
