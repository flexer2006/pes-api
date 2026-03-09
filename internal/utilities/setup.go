package utilities

import (
	"context"
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"go.uber.org/zap"
)

type LoadOptions struct {
	ConfigPath string
}

func Load[T any](ctx context.Context, opts ...LoadOptions) (*T, error) {
	Info(ctx, "loading configuration")
	cfg := new(T)
	if len(opts) > 0 && opts[0].ConfigPath != "" {
		path := opts[0].ConfigPath
		if fi, err := os.Stat(path); err == nil && !fi.IsDir() {
			if err := cleanenv.ReadConfig(path, cfg); err != nil {
				Error(ctx, "config read failed", zap.Error(err), zap.String("path", path))
				return nil, fmt.Errorf("read config %s: %w", path, err)
			}
		}
	}
	if err := cleanenv.ReadEnv(cfg); err != nil {
		Error(ctx, "env read failed", zap.Error(err))
		return nil, fmt.Errorf("read env: %w", err)
	}
	if loggable, ok := any(cfg).(interface{ LogFields() []zap.Field }); ok {
		Info(ctx, "configuration loaded", loggable.LogFields()...)
	} else {
		Info(ctx, "configuration loaded")
	}
	return cfg, nil
}
