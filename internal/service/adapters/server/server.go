package server

import (
	"context"
	"fmt"

	_ "github.com/flexer2006/case-person-enrichment-go/docs/swagger"
	"github.com/flexer2006/case-person-enrichment-go/internal/service/domain"
	"github.com/flexer2006/case-person-enrichment-go/internal/service/ports"
	logger "github.com/flexer2006/case-person-enrichment-go/internal/utilities"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

// @title PES API
// @version 1.0
// @BasePath /api/v1

type Server struct {
	app    *fiber.App
	config domain.Config
}

func New(config domain.Config, api ports.API, repositories ports.Repositories) *Server {
	app := fiber.New(fiber.Config{
		ReadTimeout:  config.Server.ReadTimeout,
		WriteTimeout: config.Server.WriteTimeout,
		AppName:      "PES",
	})
	app.Get("/swagger", func(c fiber.Ctx) error {
		r := c.Redirect()
		r.Status(fiber.StatusFound)
		return r.To("/swagger/swagger.html")
	})
	app.Get("/swagger/swagger.html", func(c fiber.Ctx) error {
		return c.SendFile("./docs/swagger/swagger.html")
	})
	app.Get("/swagger/swagger.json", func(c fiber.Ctx) error {
		return c.SendFile("./docs/swagger/swagger.json")
	})
	Setup(app, api, repositories)
	return new(Server{
		app:    app,
		config: config,
	})
}

func Setup(app *fiber.App, api ports.API, repositories ports.Repositories) {
	personHandler, v1 := NewPersonHandler(api, repositories), app.Group("/api/v1")
	persons := v1.Group("/persons")
	persons.Get("/", personHandler.GetPersons)
	persons.Get("/:id", personHandler.GetPersonByID)
	persons.Post("/", personHandler.CreatePerson)
	persons.Put("/:id", personHandler.UpdatePerson)
	persons.Patch("/:id", personHandler.UpdatePerson)
	persons.Delete("/:id", personHandler.DeletePerson)
	persons.Post("/:id/enrich", personHandler.EnrichPerson)
}

func (s *Server) Start(ctx context.Context) error {
	address := fmt.Sprintf("%s:%d", s.config.Server.Host, s.config.Server.Port)
	logger.Info(ctx, "starting HTTP server", zap.String("address", address))
	go func() {
		if err := s.app.Listen(address); err != nil {
			logger.Error(ctx, "failed to start HTTP server", zap.Error(err))
		}
	}()
	<-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	logger.Info(ctx, "stopping HTTP server")
	if err := s.app.ShutdownWithContext(ctx); err != nil {
		logger.Error(ctx, "failed to shutdown HTTP server gracefully", zap.Error(err))
		return fmt.Errorf("failed to shutdown HTTP server gracefully: %w", err)
	}
	return nil
}
