package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/brolatief/fiqih-tahawulat/internal/cache"
	"github.com/brolatief/fiqih-tahawulat/internal/config"
	"github.com/brolatief/fiqih-tahawulat/internal/handler"
	"github.com/brolatief/fiqih-tahawulat/internal/middleware"
	"github.com/brolatief/fiqih-tahawulat/internal/repository"
	"github.com/brolatief/fiqih-tahawulat/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg := config.Load()

	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo})))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	pgRepo := repository.NewPostgres(pool)

	var c cache.Cache = cache.NoopCache{}
	slog.Info("cache disabled (no redis)")

	svc := service.NewService(pgRepo, c)
	h := handler.NewHandler(svc)

	r := chi.NewRouter()

	r.Use(middleware.Recover)
	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{cfg.AllowedOrigins},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Route("/v1", func(r chi.Router) {
		r.Get("/kitab", h.ListKitab)
		r.Get("/kitab/{id}", h.GetKitab)
		r.Get("/kitab/{id}/konten", h.GetKonten)
		r.Get("/kitab/{id}/konten/{hal}", h.GetKonten)
		r.Get("/search", h.Search)
		r.Get("/kategori", h.ListKategori)
	})

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	srv := &http.Server{
		Addr:         cfg.Addr,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		slog.Info("server starting", "addr", cfg.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	<-sigCh
	slog.Info("shutting down...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("shutdown error", "error", err)
	}
}
