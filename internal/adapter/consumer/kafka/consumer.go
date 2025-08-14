package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/folivorra/get_order/internal/adapter/mapper"
	"github.com/folivorra/get_order/internal/config"
	"github.com/folivorra/get_order/internal/repository/postgres"
	"github.com/folivorra/get_order/internal/usecase"
	"github.com/segmentio/kafka-go"
	"log/slog"
	"time"
)

type Consumer struct {
	logger *slog.Logger
	cfg    config.Config
	reader *kafka.Reader
	srv    *usecase.OrderService
}

func NewConsumer(logger *slog.Logger, cfg config.Config, reader *kafka.Reader, srv *usecase.OrderService) *Consumer {
	return &Consumer{
		logger: logger,
		cfg:    cfg,
		reader: reader,
		srv:    srv,
	}
}

func NewReader(cfg config.Config) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{cfg.KafkaBrokerAddr},
		Topic:       cfg.KafkaGetOrderTopic,
		GroupID:     cfg.KafkaConsumerGroup,
		StartOffset: kafka.LastOffset,
	})
}

func (c *Consumer) Start(ctx context.Context) {
	c.logger.Info("consumer started",
		slog.String("addr", c.cfg.KafkaBrokerAddr),
	)

	for {
		select {
		case <-ctx.Done():
			c.logger.Info("consumer stopped",
				slog.String("addr", c.cfg.KafkaBrokerAddr),
			)
			return
		default:
			msg, err := c.reader.FetchMessage(ctx)
			if err != nil {
				if errors.Is(err, context.Canceled) {
					c.logger.Info("consumer stopped",
						slog.String("addr", c.cfg.KafkaBrokerAddr),
					)
					return
				}
				c.logger.Error("failed to read message, retrying",
					slog.String("error", err.Error()),
				)
				jitter := time.Duration(gofakeit.IntN(500)) * time.Millisecond
				time.Sleep(c.cfg.KafkaBackoff + jitter)
				continue
			}

			var orderDTO mapper.OrderIntoDomainDTO

			if err = json.Unmarshal(msg.Value, &orderDTO); err != nil {
				c.logger.Error("failed to unmarshal message",
					slog.String("error", err.Error()),
				)
			}

			order := mapper.ConvertToDomain(&orderDTO)
			err = c.srv.ProcessIncomingOrder(ctx, order)

			switch {
			case errors.Is(err, postgres.ErrOrderAlreadyExists):
				c.logger.Warn("order already exists",
					slog.String("uuid", order.OrderUID.String()),
				)
			case err != nil:
				c.logger.Error("failed to process order",
					slog.String("uuid", order.OrderUID.String()),
					slog.String("error", err.Error()),
				)
			default:
				c.logger.Info("order has been saved in db",
					slog.String("uuid", order.OrderUID.String()),
				)

				orderView, _ := json.MarshalIndent(orderDTO, "", "	")
				c.logger.Debug("saved order view",
					slog.String("order", string(orderView)),
				)
			}

			if err = c.reader.CommitMessages(ctx, msg); err != nil {
				c.logger.Error("failed to commit message",
					slog.String("error", err.Error()),
				)
			}
		}
	}
}
