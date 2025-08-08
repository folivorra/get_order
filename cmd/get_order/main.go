package main

import (
	"context"
	broker "github.com/folivorra/get_order/internal/adapter/kafka"
	"github.com/folivorra/get_order/internal/config"
	"github.com/folivorra/get_order/internal/repository/postgres"
	"github.com/folivorra/get_order/internal/usecase"
	"log/slog"
	"os"
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

	repo := postgres.NewPgOrderRepo()

	service := usecase.NewOrderService(logger, cfg, repo)

	reader := broker.NewKafkaReader(cfg)
	consumer := broker.NewConsumer(logger, cfg, reader, service)

	consumer.Start(ctx)
}
