package usecase

import (
	"context"
	"errors"
	"github.com/folivorra/get_order/internal/config"
	"github.com/folivorra/get_order/internal/domain"
	"github.com/google/uuid"
	"log/slog"
)

type OrderRepo interface {
	Get(ctx context.Context, uid uuid.UUID) (order *domain.Order, err error)
	Save(ctx context.Context, order *domain.Order) (err error)
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

func (s *OrderService) ProcessIncomingOrder(ctx context.Context, order *domain.Order) error {
	exists, err := s.repo.Exists(ctx, order.OrderUID)
	if err != nil {
		return err
	}

	if exists {
		return errors.New("order already exists")
	}

	// need to give uuid for objects before save in repo
	order.Delivery.DeliveryUID = uuid.New()
	order.Payment.PaymentUID = uuid.New()
	for i := range order.Items {
		order.Items[i].ItemUID = uuid.New()
		order.Items[i].OrderUID = order.OrderUID
	}

	return s.repo.Save(ctx, order)
}
