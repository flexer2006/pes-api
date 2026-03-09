package main

import (
	"context"
	"os"
	"time"

	"github.com/flexer2006/case-person-enrichment-go/internal/service/adapters/postgres"
	"github.com/flexer2006/case-person-enrichment-go/internal/service/app"
	"github.com/flexer2006/case-person-enrichment-go/internal/service/domain"
	"github.com/flexer2006/case-person-enrichment-go/internal/service/logger"
	"go.uber.org/zap"
)

func main() {
	if err := run(); err != nil {
		os.Exit(1)
	}
}

func run() error {
	logger.SetGlobal(logger.NewConsole(logger.InfoLevel, true))
	defer func() { _ = logger.Global().Sync() }()
	ctx := context.Background()
	cfg, err := app.Load[domain.Config](ctx, app.LoadOptions{ConfigPath: "./deploy/.env"})
	if err != nil {
		logger.Error(ctx, "load config", zap.Error(err))
		return err
	}
	var finalLogger *logger.Logger
	if cfg.Logger.Model == "production" {
		finalLogger, err = logger.NewProduction()
	} else {
		if cfg.Logger.Model != "development" {
			logger.Warn(ctx, "unknown logger model, using development", zap.String("model", cfg.Logger.Model))
		}
		finalLogger, err = logger.NewDevelopment()
	}
	if err != nil {
		logger.Error(ctx, "init logger", zap.Error(err))
		return err
	}
	logger.SetGlobal(finalLogger)
	shutdownTimeout, err := time.ParseDuration(cfg.Graceful.ShutdownTimeout)
	if err != nil {
		logger.Error(ctx, "bad duration, defaulting", zap.Error(err))
		shutdownTimeout = 5 * time.Second
	}
	logger.Info(ctx, "initializing database")
	dbCfg := postgres.Config{
		Postgres: postgres.PostgresConfig{
			Host:     cfg.Postgres.Host,
			Port:     cfg.Postgres.Port,
			User:     cfg.Postgres.User,
			Password: cfg.Postgres.Password,
			Database: cfg.Postgres.Database,
			SSLMode:  cfg.Postgres.SSLMode,
			MinConns: cfg.Postgres.PoolMinConns,
			MaxConns: cfg.Postgres.PoolMaxConns,
		},
		MigrationsPath:  cfg.Migrations.Path,
		ApplyMigrations: true,
	}
	data, err := postgres.NewDatabase(ctx, dbCfg)
	if err != nil {
		logger.Error(ctx, "init db", zap.Error(err))
		return err
	}
	if err := data.Ping(ctx); err != nil {
		logger.Error(ctx, "db ping", zap.Error(err))
		return err
	}
	logger.Info(ctx, "database ready")
	application, err := app.NewApplication(ctx, cfg, data, nil)
	if err != nil {
		logger.Error(ctx, "init app", zap.Error(err))
		return err
	}
	appCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	go func() {
		if err := application.Start(appCtx); err != nil {
			logger.Error(ctx, "app stopped", zap.Error(err))
			cancel()
		}
	}()
	logger.Info(ctx, "service started", zap.String("env", cfg.Logger.Model), zap.String("level", cfg.Logger.Level))
	err = app.Shutdown(ctx, shutdownTimeout,
		func(ctx context.Context) error { cancel(); return application.Stop(ctx) },
		func(ctx context.Context) error { logger.Info(ctx, "closing db"); data.Close(ctx); return nil },
	)
	if err != nil {
		logger.Error(ctx, "shutdown error", zap.Error(err))
	}
	logger.Info(ctx, "shutdown complete")
	return nil
}
