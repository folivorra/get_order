package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/folivorra/get_order/internal/adapter/mapper"
	"log/slog"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/folivorra/get_order/internal/config"
	"github.com/folivorra/get_order/internal/domain"
	"github.com/folivorra/get_order/internal/usecase"
	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	logger slog.Logger
	cfg    config.Config
	reader *kafka.Reader
	srv    *usecase.OrderService
}

func NewConsumer(logger *slog.Logger, cfg config.Config, reader *kafka.Reader, srv *usecase.OrderService) *Consumer {
	return &Consumer{
		logger: *logger,
		cfg:    cfg,
		reader: reader,
		srv:    srv,
	}
}

func NewReader(cfg config.Config) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{cfg.Broker},
		Topic:       cfg.GetOrderTopic,
		GroupID:     cfg.ConsumerGroup,
		StartOffset: kafka.LastOffset,
	})
}

func (c *Consumer) Start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			c.Close()
			return
		default:
			msg, err := c.reader.FetchMessage(ctx)
			if err != nil {
				if errors.Is(err, context.Canceled) {
					return
				}
				c.logger.Error("failed to read message",
					slog.String("error", err.Error()),
				)
				continue
			}

			var orderDTO mapper.OrderIntoDomainDTO

			if err = json.Unmarshal(msg.Value, &orderDTO); err != nil {
				c.logger.Error("failed to unmarshal message",
					slog.String("error", err.Error()),
				)
				continue
			}

			order := orderDTO.ConvertToDomain()

			if err = c.retryAndBackoff(ctx, order, 5, 1*time.Second); err != nil {
				c.logger.Warn("message is duplicated",
					slog.String("uuid", order.OrderUID.String()),
					slog.String("error", err.Error()),
				)
				continue
			}

			if err = c.reader.CommitMessages(ctx, msg); err != nil {
				c.logger.Error("failed to commit message",
					slog.String("error", err.Error()),
				)
			}

			c.logger.Info("order has been saved in db",
				slog.String("uuid", order.OrderUID.String()),
			)
		}
	}
}

func (c *Consumer) retryAndBackoff(ctx context.Context, order *domain.Order, retries int, backoff time.Duration) error {
	_ = gofakeit.Seed(0)
	var err error

	for attempt := 0; attempt < retries; attempt++ {
		if err = c.srv.ProcessIncomingOrder(ctx, order); err == nil {
			return nil
		}

		time.Sleep(backoff + time.Duration(gofakeit.IntN(2000))*time.Millisecond)
	}

	return err
}

func (c *Consumer) Close() {
	if err := c.reader.Close(); err != nil {
		c.logger.Warn("failed to close kafka reader",
			slog.String("error", err.Error()),
		)
	}
}
