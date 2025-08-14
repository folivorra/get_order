package config

import (
	"encoding/json"
	"github.com/caarlos0/env"
	"log/slog"
	"time"
)

type Config struct {
	KafkaBrokerAddr             string        `env:"KAFKA_BROKER" envDefault:"kafka:9092"`
	KafkaGetOrderTopic          string        `env:"KAFKA_GET_ORDER_TOPIC" envDefault:"get_orders"`
	KafkaConsumerGroup          string        `env:"KAFKA_CONSUMER_GROUP" envDefault:"default"`
	KafkaBackoff                time.Duration `env:"KAFKA_BACKOFF" envDefault:"500ms"`
	PgDsn                       string        `env:"PG_DSN" envDefault:"postgres://app:app@postgres:5432/?sslmode=disable"`
	PgSaveTimeout               time.Duration `env:"PG_SAVE_TIMEOUT" envDefault:"5s"`
	PgExistsTimeout             time.Duration `env:"PG_EXISTS_TIMEOUT" envDefault:"2s"`
	PgGetTimeout                time.Duration `env:"PG_GET_TIMEOUT" envDefault:"2s"`
	PgPoolMaxOpenConns          int           `env:"PG_POOL_MAX_OPEN_CONNS" envDefault:"10"`
	PgPoolMaxIdleConns          int           `env:"PG_POOL_MAX_IDLE_CONNS" envDefault:"5"`
	PgPoolConnMaxLifetime       time.Duration `env:"PG_POOL_CONN_MAX_LIFETIME" envDefault:"1h"`
	PgPingTimeout               time.Duration `env:"PG_PING_TIMEOUT" envDefault:"500ms"`
	PgMaxRetries                int           `env:"PG_MAX_RETRIES" envDefault:"3"`
	PgBackoff                   time.Duration `env:"PG_BACKOFF" envDefault:"500ms"`
	ServerHTTPPort              string        `env:"SERVER_HTTP_PORT" envDefault:"8080"`
	ServerHTTPShutdownTimeout   time.Duration `env:"SERVER_HTTP_SHUTDOWN_TIMEOUT" envDefault:"5s"`
	ServerHTTPReadHeaderTimeout time.Duration `env:"SERVER_HTTP_READ_HEADER_TIMEOUT" envDefault:"5s"`
	ServerHTTPReadTimeout       time.Duration `env:"SERVER_HTTP_READ_TIMEOUT" envDefault:"5s"`
	ServerHTTPWriteTimeout      time.Duration `env:"SERVER_HTTP_WRITE_TIMEOUT" envDefault:"10s"`
	ServerHTTPIdleTimeout       time.Duration `env:"SERVER_HTTP_IDLE_TIMEOUT" envDefault:"120s"`
	CacheCapacity               int           `env:"CACHE_CAPACITY" envDefault:"10"`
}

func NewConfig(logger *slog.Logger) Config {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		logger.Warn("fail to parse config, used default values",
			slog.String("error", err.Error()),
		)
	}

	cfgView, _ := json.MarshalIndent(cfg, "", "	")
	logger.Debug("env variables",
		slog.String("config", string(cfgView)),
	)

	return cfg
}
