package app

import (
	"context"
	"fmt"
	"time"

	"github.com/flexer2006/case-person-enrichment-go/internal/service/adapters/enrichment"
	"github.com/flexer2006/case-person-enrichment-go/internal/service/adapters/postgres"
	"github.com/flexer2006/case-person-enrichment-go/internal/service/adapters/server"
	"github.com/flexer2006/case-person-enrichment-go/internal/service/domain"
	"github.com/flexer2006/case-person-enrichment-go/internal/service/ports"
	logger "github.com/flexer2006/case-person-enrichment-go/internal/utilities"
	dbpkg "github.com/flexer2006/case-person-enrichment-go/internal/utilities/database"

	"go.uber.org/zap"
)

type Application struct {
	db         dbpkg.PostgresProvider
	config     *domain.Config
	httpServer *server.Server
}

func NewApplication(ctx context.Context, config *domain.Config, database dbpkg.PostgresProvider, apiAdapter ports.API) (*Application, error) {
	logger.Info(ctx, "initializing application")
	repos := postgres.New(database)
	if apiAdapter == nil {
		apiAdapter = enrichment.NewAPI()
	}
	app := new(Application{
		config:     config,
		db:         database,
		httpServer: server.New(*config, apiAdapter, repos),
	})
	logger.Info(ctx, "application initialized successfully")
	return app, nil
}

func (a *Application) Start(ctx context.Context) error {
	logger.Info(ctx, "starting application")
	if err := a.httpServer.Start(ctx); err != nil {
		return fmt.Errorf("failed to start HTTP server: %w", err)
	}
	return nil
}

func (a *Application) Stop(ctx context.Context) error {
	logger.Info(ctx, "stopping application")
	shutdownTimeout, err := time.ParseDuration(a.config.Graceful.ShutdownTimeout)
	if err != nil {
		shutdownTimeout = 5 * time.Second
		logger.Warn(ctx, "invalid graceful shutdown timeout, using default", zap.String("default", shutdownTimeout.String()))
	}
	ctx, cancel := context.WithTimeout(ctx, shutdownTimeout)
	defer cancel()
	if err := a.httpServer.Stop(ctx); err != nil {
		logger.Error(ctx, "error stopping HTTP server", zap.Error(err))
	}
	if a.db != nil {
		a.db.Close(ctx)
	}
	logger.Info(ctx, "application stopped")
	return nil
}
