package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/brolatief/fiqih-tahawulat/internal/cache"
	"github.com/brolatief/fiqih-tahawulat/internal/config"
	mcpserver "github.com/brolatief/fiqih-tahawulat/internal/mcp"
	"github.com/brolatief/fiqih-tahawulat/internal/repository"
	"github.com/brolatief/fiqih-tahawulat/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	cfg := config.Load()

	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo})))

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	pgRepo := repository.NewPostgres(pool)
	c := cache.NoopCache{}
	svc := service.NewService(pgRepo, c)

	mcpSrv := mcpserver.NewServer(svc)
	mcpServer := server.NewMCPServer(
		"Fiqih Tahawulat MCP",
		"1.0.0",
	)

	mcpSrv.RegisterHandlers(mcpServer)

	slog.Info("MCP server starting (stdio)")
	server.ServeStdio(mcpServer)
}
