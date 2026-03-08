package domain

import "time"

type Config struct {
	Postgres struct {
		Host                string `env:"POSTGRES_HOST" env-default:"localhost"`
		User                string `env:"POSTGRES_USER" env-default:"postgres"`
		Password            string `env:"POSTGRES_PASSWORD" env-default:"postgres"`
		Database            string `env:"POSTGRES_DB" env-default:"postgres"`
		SSLMode             string `env:"POSTGRES_SSLMODE" env-default:"disable"`
		ConnectTimeout      string `env:"PGX_CONNECT_TIMEOUT" env-default:"10s"`
		AcquireTimeout      string `env:"PGX_POOL_ACQUIRE_TIMEOUT" env-default:"30s"`
		Port                int    `env:"POSTGRES_PORT" env-default:"5432"`
		PoolMinConns        int    `env:"PGX_POOL_MIN_CONNS" env-default:"2"`
		PoolMaxConns        int    `env:"PGX_POOL_MAX_CONNS" env-default:"10"`
		PoolMaxConnLifetime int    `env:"PGX_POOL_MAX_CONN_LIFETIME" env-default:"3600"`
		PoolMaxConnIdleTime int    `env:"PGX_POOL_MAX_CONN_IDLE_TIME" env-default:"600"`
	}

	Logger struct {
		Level        string `env:"LOGGER_LEVEL" env-default:"info"`
		Format       string `env:"LOGGER_FORMAT" env-default:"json"`
		Output       string `env:"LOGGER_OUTPUT" env-default:"stdout"`
		TimeEncoding string `env:"LOGGER_TIME_ENCODING" env-default:"iso8601"`
		Model        string `env:"LOGGER_MODEL" env-default:"development"`
		Caller       bool   `env:"LOGGER_CALLER" env-default:"true"`
		Stacktrace   bool   `env:"LOGGER_STACKTRACE" env-default:"true"`
	}

	Server struct {
		Host         string        `env:"HTTP_HOST" env-default:"0.0.0.0"`
		Port         int           `env:"HTTP_PORT" env-default:"8080"`
		ReadTimeout  time.Duration `env:"HTTP_READ_TIMEOUT" env-default:"5s"`
		WriteTimeout time.Duration `env:"HTTP_WRITE_TIMEOUT" env-default:"10s"`
	}

	Migrations struct {
		Path string `env:"MIGRATIONS_DIR" env-default:"./migrations"`
	}

	Graceful struct {
		ShutdownTimeout string `env:"GRACEFUL_SHUTDOWN_TIMEOUT" env-default:"5s"`
	}
}
