package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/alysonryan/message-service/backgrounds"
	"github.com/alysonryan/message-service/config"
	"github.com/alysonryan/message-service/db"
	"github.com/alysonryan/message-service/models"
	"github.com/alysonryan/message-service/repositories"
	"github.com/alysonryan/message-service/services"
	"github.com/alysonryan/message-service/workers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	if err := run(); err != nil {
		slog.Error("service stopped with error", "error", err)
		os.Exit(1)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	conn, err := db.Init(cfg)
	if err != nil {
		return fmt.Errorf("failed to init database: %w", err)
	}
	defer conn.Close()

	repos := repositories.NewRepositoriesInstance(conn)
	svcs := services.NewServicesInstance(cfg, repos)

	queue := make(chan models.PendingMessage, 100)

	wkrs := workers.NewWorkersInstance(queue, svcs, repos)
	bgs := backgrounds.NewBackgroundsInstance(repos, queue)

	var wg sync.WaitGroup

	wg.Add(2)
	go func() {
		defer wg.Done()
		wkrs.WhatsApp.Run(ctx)
	}()
	go func() {
		defer wg.Done()
		bgs.PollingDB.Run(ctx)
	}()

	slog.Info("service started")

	<-ctx.Done()
	slog.Info("shutting down...")

	wg.Wait()
	slog.Info("shutdown complete")

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	return nil
}