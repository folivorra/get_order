package usecase

import (
	"context"
	"errors"
	"log/slog"

	"github.com/folivorra/get_order/internal/config"
	"github.com/folivorra/get_order/internal/domain"
	"github.com/google/uuid"
)

type OrderRepo interface {
	Get(uid uuid.UUID) (order domain.Order, err error)
	Save(order domain.Order) (err error)
	Exists(ctx context.Context, uuid uuid.UUID) (exists bool, err error)
}

type OrderService struct {
	ctx    context.Context
	logger *slog.Logger
	cfg    config.Config
	repo   OrderRepo
}

func NewOrderService(ctx context.Context, logger *slog.Logger, cfg config.Config, repo OrderRepo) *OrderService {
	return &OrderService{
		ctx:    ctx,
		logger: logger,
		cfg:    cfg,
		repo:   repo,
	}
}

func (s *OrderService) ProcessIncomingOrder(order *domain.Order) error {
	ctx, cancel := context.WithTimeout(s.ctx, s.cfg.Timeout)
	defer cancel()

	exists, err := s.repo.Exists(ctx, order.OrderUID)
	if err != nil {
		return err
	}

	if exists {
		return errors.New("order already exists")
	}

	return s.repo.Save(*order)
}
