package config

import (
	"github.com/caarlos0/env"
	"log/slog"
	"time"
)

type Config struct {
	Port          string        `env:"PORT" envDefault:"8080"`
	ServerTimeout time.Duration `env:"SERVER_TIMEOUT" envDefault:"5s"`
	Broker        string        `env:"BROKER" envDefault:"kafka:9092"`
	GetOrderTopic string        `env:"GET_ORDER_TOPIC" envDefault:"get_orders"`
	ConsumerGroup string        `env:"CONSUMER_GROUP" envDefault:"default"`
	PgDsn         string        `env:"PG_DSN" envDefault:"postgres://app:app@postgres:5432/?sslmode=disable"`
	SaveTimeout   time.Duration `env:"SAVE_TIMEOUT" envDefault:"5s"`
	ExistsTimeout time.Duration `env:"EXISTS_TIMEOUT" envDefault:"2s"`
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
