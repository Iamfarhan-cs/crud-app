package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/Iamfarhan-cs/crud-app/internal/config"
)

const (
	maxRequestBodyBytes = 1 << 20
	readTimeout         = 5 * time.Second
	writeTimeout        = 10 * time.Second
	idleTimeout         = 60 * time.Second
	shutdownTimeout     = 10 * time.Second
)

func main() {
	cfg := config.Load()

	// The API command is responsible for application wiring only.
	// It should compose config, database, repositories, services, and handlers.
	// Business rules, SQL queries, and HTTP request handling must not live here.
	router := http.NewServeMux()
	router.HandleFunc("/healthz", healthzHandler)

	server := &http.Server{
		Addr:         ":" + cfg.HTTPPort,
		Handler:      http.MaxBytesHandler(router, maxRequestBodyBytes),
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Printf("task management API listening on port %s in %s", cfg.HTTPPort, cfg.Environment)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server failed: %v", err)
		}
	}()

	<-ctx.Done()
	stop()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("server shutdown failed: %v", err)
	}

	log.Print("server stopped")
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok\n"))
}
