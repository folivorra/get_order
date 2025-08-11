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

type OrderCache interface {
	Get(ctx context.Context, uid uuid.UUID) (order *domain.Order, err error)
	Set(ctx context.Context, order *domain.Order) (err error)
}

type OrderService struct {
	logger *slog.Logger
	cfg    config.Config
	repo   OrderRepo
	//cache  OrderCache
}

func NewOrderService(logger *slog.Logger, cfg config.Config, repo OrderRepo) *OrderService {
	return &OrderService{
		logger: logger,
		cfg:    cfg,
		repo:   repo,
		//cache:  cache,
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
		order.Items[i].OrderItemUID = uuid.New()
	}

	return s.repo.Save(ctx, order)
}

func (s *OrderService) GetOrder(ctx context.Context, uuid uuid.UUID) (*domain.Order, error) {
	//order, err := s.cache.Get(ctx, uuid)
	//if err == nil {
	//	return order, nil
	//}

	order, err := s.repo.Get(ctx, uuid)
	if err != nil {
		return nil, err
	}

	//err = s.cache.Set(ctx, order) // todo mb ignore err

	return order, err
}
