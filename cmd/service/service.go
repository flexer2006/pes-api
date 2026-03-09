package main

import (
	"context"
	"os"
	"time"

	"github.com/flexer2006/case-person-enrichment-go/internal/service/app"
	"github.com/flexer2006/case-person-enrichment-go/internal/service/domain"
	. "github.com/flexer2006/case-person-enrichment-go/internal/utilities" //nolint:staticcheck
	"github.com/flexer2006/case-person-enrichment-go/internal/utilities/database"

	"go.uber.org/zap"
)

func main() {
	if err := run(); err != nil {
		os.Exit(1)
	}
}

func run() error {
	SetGlobal(NewConsole(InfoLevel, true))
	defer func() { _ = Global().Sync() }()
	ctx := context.Background()
	cfg, err := Load[domain.Config](ctx, LoadOptions{ConfigPath: "./deploy/.env"})
	if err != nil {
		Error(ctx, "load config", zap.Error(err))
		return err
	}
	var finalLogger *Logger
	if cfg.Logger.Model == "production" {
		finalLogger, err = NewProduction()
	} else {
		if cfg.Logger.Model != "development" {
			Warn(ctx, "unknown logger model, using development", zap.String("model", cfg.Logger.Model))
		}
		finalLogger, err = NewDevelopment()
	}
	if err != nil {
		Error(ctx, "init logger", zap.Error(err))
		return err
	}
	SetGlobal(finalLogger)
	shutdownTimeout, err := time.ParseDuration(cfg.Graceful.ShutdownTimeout)
	if err != nil {
		Error(ctx, "bad duration, defaulting", zap.Error(err))
		shutdownTimeout = 5 * time.Second
	}
	Info(ctx, "initializing database")
	dbCfg := database.Config{
		Postgres: database.PostgresConfig{
			Host:     cfg.Postgres.Host,
			Port:     cfg.Postgres.Port,
			User:     cfg.Postgres.User,
			Password: cfg.Postgres.Password,
			Database: cfg.Postgres.Database,
			SSLMode:  cfg.Postgres.SSLMode,
			MinConns: cfg.Postgres.PoolMinConns,
			MaxConns: cfg.Postgres.PoolMaxConns,
		},
		Migrate:         database.MigrateConfig{Path: cfg.Migrations.Path},
		ApplyMigrations: true,
	}
	data, err := database.New(ctx, dbCfg)
	if err != nil {
		Error(ctx, "init db", zap.Error(err))
		return err
	}
	if err := data.Ping(ctx); err != nil {
		Error(ctx, "db ping", zap.Error(err))
		return err
	}
	switch version, dirty, err := data.GetMigrationVersion(ctx); {
	case err != nil:
		Warn(ctx, "migration version", zap.Error(err))
	case dirty:
		Warn(ctx, "dirty migration", zap.Uint("version", version))
	default:
		Info(ctx, "migration", zap.Uint("version", version))
	}
	Info(ctx, "database ready")
	application, err := app.NewApplication(ctx, cfg, data.Provider(), nil)
	if err != nil {
		Error(ctx, "init app", zap.Error(err))
		return err
	}
	appCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	go func() {
		if err := application.Start(appCtx); err != nil {
			Error(ctx, "app stopped", zap.Error(err))
			cancel()
		}
	}()
	Info(ctx, "service started", zap.String("env", cfg.Logger.Model), zap.String("level", cfg.Logger.Level))
	err = Shutdown(ctx, shutdownTimeout,
		func(ctx context.Context) error { cancel(); return application.Stop(ctx) },
		func(ctx context.Context) error { Info(ctx, "closing db"); data.Close(ctx); return nil },
	)
	if err != nil {
		Error(ctx, "shutdown error", zap.Error(err))
	}
	Info(ctx, "shutdown complete")
	return nil
}
