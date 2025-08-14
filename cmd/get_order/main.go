package main

import (
	"context"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/folivorra/get_order/internal/adapter/cache/inmemory"
	"github.com/folivorra/get_order/internal/adapter/consumer/kafka"
	"github.com/folivorra/get_order/internal/adapter/controller/rest"
	"github.com/folivorra/get_order/internal/adapter/middleware"
	"github.com/folivorra/get_order/internal/config"
	"github.com/folivorra/get_order/internal/repository/postgres"
	"github.com/folivorra/get_order/internal/storage"
	"github.com/folivorra/get_order/internal/usecase"
	"github.com/gorilla/mux"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_ = gofakeit.Seed(0)

	// logger
	logger := slog.New(
		slog.NewTextHandler(
			os.Stdout, &slog.HandlerOptions{
				Level:     slog.LevelDebug,
				AddSource: true,
			},
		),
	)

	// cfg
	cfg := config.NewConfig(logger)

	// postgres | repo
	pgClient := storage.NewPgClient(ctx, cfg)
	defer func() {
		_ = pgClient.Close()
	}()
	pgRepo := postgres.NewPgOrderRepo(pgClient, cfg)

	// inmemory | cache
	inMemCache := inmemory.NewInMemOrderCache(logger, cfg.CacheCapacity)

	// service layer
	service := usecase.NewOrderService(logger, cfg, pgRepo, inMemCache)

	// kafka
	kafkaReader := kafka.NewReader(cfg)
	kafkaConsumer := kafka.NewConsumer(logger, cfg, kafkaReader, service)
	go kafkaConsumer.Start(ctx)
	defer func() {
		_ = kafkaReader.Close()
	}()

	// router mux
	router := mux.NewRouter()
	router.Use(middleware.LoggingMiddleware(logger))

	// http | controller
	controller := rest.NewController(service, cfg, logger)
	controller.RegisterRoutes(router)

	// http | server
	server := rest.NewServer(
		&http.Server{
			Addr:              ":" + cfg.ServerHTTPPort,
			Handler:           router,
			ReadHeaderTimeout: cfg.ServerHTTPReadHeaderTimeout,
			ReadTimeout:       cfg.ServerHTTPReadTimeout,
			WriteTimeout:      cfg.ServerHTTPWriteTimeout,
			IdleTimeout:       cfg.ServerHTTPIdleTimeout,
		},
		cfg,
		logger,
	)

	go func() {
		if err := server.Run(); err != nil {
			logger.Error("failed to start server",
				slog.String("addr", cfg.ServerHTTPPort),
				slog.String("err", err.Error()),
			)
		}
	}()
	defer server.Stop(ctx)

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	<-shutdown
	cancel()
}
