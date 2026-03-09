package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/flexer2006/case-person-enrichment-go/internal/service/domain"
	"github.com/flexer2006/case-person-enrichment-go/internal/service/logger"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Config struct {
	Postgres        PostgresConfig
	MigrationsPath  string
	ApplyMigrations bool
}

type PostgresConfig struct {
	Host     string
	Port     int
	User     string
	Password string //nolint:gosec
	Database string
	SSLMode  string
	MinConns int
	MaxConns int
}

func (c PostgresConfig) Validate() error {
	if c.Host == "" || c.Port == 0 || c.User == "" || c.Database == "" {
		return domain.ErrInvalidConfiguration
	}
	return nil
}

func (c PostgresConfig) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s", c.User, c.Password, c.Host, c.Port, c.Database, c.SSLMode)
}

type Database struct {
	pool *pgxpool.Pool
	dsn  string
}

func NewDatabase(ctx context.Context, cfg Config) (*Database, error) {
	if err := cfg.Postgres.Validate(); err != nil {
		logger.Error(ctx, "invalid database configuration", zap.Error(err))
		return nil, err
	}
	dsn := cfg.Postgres.DSN()
	pool, err := connect(ctx, dsn, cfg.Postgres.MinConns, cfg.Postgres.MaxConns)
	if err != nil {
		return nil, fmt.Errorf("setup db: %w", err)
	}
	db := new(Database{pool: pool, dsn: dsn})
	if cfg.ApplyMigrations && cfg.MigrationsPath != "" {
		if err := runMigrations(ctx, cfg.MigrationsPath, dsn); err != nil {
			db.Close(ctx)
			return nil, err
		}
	}
	return db, nil
}

func (d *Database) Pool() *pgxpool.Pool { return d.pool }

func (d *Database) Close(ctx context.Context) {
	if d.pool != nil {
		logger.Info(ctx, "closing postgres database connection")
		d.pool.Close()
	}
}

func (d *Database) Ping(ctx context.Context) error {
	if d.pool == nil {
		return fmt.Errorf("ping: pool nil")
	}
	return d.pool.Ping(ctx)
}

func (d *Database) GetDSN() string { return d.dsn }

func connect(ctx context.Context, dsn string, minConn, maxConn int) (*pgxpool.Pool, error) {
	logger.Info(ctx, "connecting to postgres database")
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		logger.Error(ctx, "parse config failed", zap.Error(err))
		return nil, fmt.Errorf("parse config: %w", err)
	}
	cfg.MinConns, cfg.MaxConns = clamp32(minConn), clamp32(maxConn)
	cfg.ConnConfig.ConnectTimeout = 5 * time.Second
	cfg.HealthCheckPeriod = 1 * time.Minute
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		logger.Error(ctx, "create pool failed", zap.Error(err))
		return nil, fmt.Errorf("create pool: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		logger.Error(ctx, "ping failed", zap.Error(err))
		return nil, fmt.Errorf("ping: %w", err)
	}
	logger.Info(ctx, "connected to postgres database")
	return pool, nil
}

func runMigrations(ctx context.Context, path, dsn string) error {
	mig, err := migrate.New("file://"+path, dsn)
	if err != nil {
		logger.Error(ctx, "failed to create migration instance", zap.Error(err), zap.String("path", path))
		return fmt.Errorf("migrations: %w", err)
	}
	srcErr, dbErr := mig.Close()
	if srcErr != nil {
		logger.Error(ctx, "failed to close migration source", zap.Error(srcErr))
	}
	if dbErr != nil {
		logger.Error(ctx, "failed to close migration database", zap.Error(dbErr))
	}
	if err := mig.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		logger.Error(ctx, "migration apply failed", zap.Error(err))
		return fmt.Errorf("migrations: %w", err)
	}
	logger.Info(ctx, "database migrations applied")
	return nil
}
