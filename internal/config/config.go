package config

import (
	"github.com/caarlos0/env"
	"log/slog"
	"time"
)

type Config struct {
	Port          string        `env:"PORT" envDefault:"8080"`
	Timeout       time.Duration `env:"TIMEOUT" envDefault:"5s"`
	Broker        string        `env:"BROKER" envDefault:"kafka:9092"`
	GetOrderTopic string        `env:"GET_ORDER_TOPIC" envDefault:"get_orders"`
	ConsumerGroup string        `env:"CONSUMER_GROUP" envDefault:"default"`
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
