package main

import (
	"context"
	"github.com/folivorra/get_order/internal/storage"
	"log/slog"
	"os"

	broker "github.com/folivorra/get_order/internal/adapter/kafka"
	"github.com/folivorra/get_order/internal/config"
	"github.com/folivorra/get_order/internal/repository/postgres"
	"github.com/folivorra/get_order/internal/usecase"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := slog.New(
		slog.NewTextHandler(
			os.Stdout, &slog.HandlerOptions{
				Level:     slog.LevelDebug,
				AddSource: true,
			},
		),
	)

	cfg := config.NewConfig(logger)

	pgClient := storage.NewPgClient(ctx, cfg.PgDsn)
	defer func() {
		if err := pgClient.Close(); err != nil {
			logger.Warn("failed to close pgClient")
		}
	}()
	pgRepo := postgres.NewPgOrderRepo(pgClient)

	service := usecase.NewOrderService(logger, cfg, pgRepo)

	reader := broker.NewKafkaReader(cfg)
	consumer := broker.NewConsumer(logger, cfg, reader, service)

	consumer.Start(ctx)
}
