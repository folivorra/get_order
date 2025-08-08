package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/folivorra/get_order/internal/config"
	"github.com/folivorra/get_order/internal/domain"
	"github.com/folivorra/get_order/internal/usecase"
	"github.com/segmentio/kafka-go"
	"log/slog"
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

func NewKafkaReader(cfg config.Config) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{cfg.Broker},
		Topic:       cfg.GetOrderTopic,
		GroupID:     "mygroup",
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
			msg, err := c.reader.ReadMessage(ctx)
			if err != nil {
				c.logger.Error("failed to read message",
					slog.String("error", err.Error()),
				)
			}

			var order domain.Order

			if err = json.Unmarshal(msg.Value, &order); err != nil {
				c.logger.Error("failed to unmarshal message",
					slog.String("error", err.Error()),
				)
			}

			if err = c.srv.ProcessIncomingOrder(&order); err != nil {
				c.logger.Error("failed to process order",
					slog.String("error", err.Error()),
				)
			}

			fmt.Println(order)
		}
	}
}

func (c *Consumer) Close() {
	if err := c.reader.Close(); err != nil {
		c.logger.Warn("failed to close kafka reader",
			slog.String("error", err.Error()),
		)
	}
}
