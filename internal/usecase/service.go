package usecase

import (
	"errors"
	"log/slog"

	"github.com/folivorra/get_order/internal/config"
	"github.com/folivorra/get_order/internal/domain"
	"github.com/google/uuid"
)

type OrderRepo interface {
	Get(uid uuid.UUID) (order domain.Order, err error)
	Save(order domain.Order) (err error)
	Exists(uuid uuid.UUID) (exists bool, err error)
}

type OrderService struct {
	logger *slog.Logger
	cfg    config.Config
	repo   OrderRepo
}

func NewOrderService(logger *slog.Logger, cfg config.Config, repo OrderRepo) *OrderService {
	return &OrderService{
		logger: logger,
		cfg:    cfg,
		repo:   repo,
	}
}

func (s *OrderService) ProcessIncomingOrder(order *domain.Order) error {
	exist, err := s.repo.Exists(order.OrderUID)
	if err != nil {
		return err
	}

	if exist {
		return errors.New("order already exists")
	}

	return s.repo.Save(*order)
}
