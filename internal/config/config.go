package config

import (
	"github.com/caarlos0/env"
	"log/slog"
)

type Config struct {
}

func NewConfig(logger *slog.Logger) Config {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		logger.Warn("fail to parse config, used default values",
			slog.String("error", err.Error()),
		)
	}

	return cfg
}
