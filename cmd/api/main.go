package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Iamfarhan-cs/crud-app/internal/config"
	"github.com/Iamfarhan-cs/crud-app/internal/database"
	"github.com/Iamfarhan-cs/crud-app/internal/task"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	startupCtx, cancelStartup := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelStartup()

	db, err := database.OpenPostgres(startupCtx, database.PostgresConfig{
		DatabaseURL:           cfg.DatabaseURL,
		MaxOpenConnections:    cfg.DBMaxOpenConnections,
		MaxIdleConnections:    cfg.DBMaxIdleConnections,
		ConnectionMaxLifetime: cfg.DBConnectionMaxLife,
	})
	if err != nil {
		log.Fatalf("open postgres: %v", err)
	}
	defer db.Close()

	taskRepository := task.NewPostgresRepository(db)
	taskService := task.NewService(taskRepository)
	taskHandler := task.NewHandler(taskService)

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", healthzHandler)
	taskHandler.RegisterRoutes(mux)

	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      mux,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	go func() {
		log.Printf("server starting on port %s", cfg.Port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server failed: %v", err)
		}
	}()

	shutdownSignal := make(chan os.Signal, 1)
	signal.Notify(shutdownSignal, os.Interrupt, syscall.SIGTERM)
	<-shutdownSignal
	signal.Stop(shutdownSignal)

	log.Print("server shutdown started")
	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancelShutdown()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("server shutdown failed: %v", err)
	}

	log.Print("server shutdown complete")
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
