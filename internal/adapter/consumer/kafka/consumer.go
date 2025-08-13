package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/folivorra/get_order/internal/adapter/mapper"
	"github.com/folivorra/get_order/internal/config"
	"github.com/folivorra/get_order/internal/usecase"
	"github.com/segmentio/kafka-go"
	"log/slog"
	"time"
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
			c.Close()
			c.logger.Info("consumer stopped",
				slog.String("addr", c.cfg.KafkaBrokerAddr),
			)
			return
		default:
			msg, err := c.reader.FetchMessage(ctx)
			if err != nil {
				if errors.Is(err, context.Canceled) {
					c.Close()
					c.logger.Info("consumer stopped",
						slog.String("addr", c.cfg.KafkaBrokerAddr),
					)
					return
				}
				c.logger.Error("failed to read message, retrying",
					slog.String("error", err.Error()),
				)
				continue
			}

			var orderDTO mapper.OrderIntoDomainDTO
			commit := false

			if err = json.Unmarshal(msg.Value, &orderDTO); err != nil {
				commit = true
				c.logger.Error("failed to unmarshal message",
					slog.String("error", err.Error()),
				)
			} else {
				order := mapper.ConvertToDomain(&orderDTO)
				err = c.srv.ProcessIncomingOrder(ctx, order)

				switch {
				case errors.Is(err, usecase.ErrOrderAlreadyExists):
					commit = true
					c.logger.Warn("order already exists",
						slog.String("uuid", order.OrderUID.String()),
					)
				case err != nil:
					c.logger.Error("failed to process order, retrying",
						slog.String("uuid", order.OrderUID.String()),
						slog.String("error", err.Error()),
					)

					jitter := time.Duration(gofakeit.IntN(500)) * time.Millisecond

					time.Sleep(c.cfg.KafkaBackoff + jitter)
				default:
					commit = true
					c.logger.Info("order has been saved in db",
						slog.String("uuid", order.OrderUID.String()),
					)

					orderView, _ := json.MarshalIndent(orderDTO, "", "	")
					c.logger.Debug("saved order view",
						slog.String("order", string(orderView)),
					)
				}
			}

			if commit {
				if err = c.reader.CommitMessages(ctx, msg); err != nil {
					c.logger.Error("failed to commit message",
						slog.String("error", err.Error()),
					)
				}
			}
		}
	}
}

func (c *Consumer) Close() {
	if err := c.reader.Close(); err != nil {
		c.logger.Warn("failed to close kafka reader",
			slog.String("error", err.Error()),
		)
	} else {
		c.logger.Info("kafka reader has been closed")
	}
}
