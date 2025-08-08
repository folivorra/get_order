package usecase

import (
	"github.com/folivorra/get_order/internal/config"
	"github.com/folivorra/get_order/internal/domain"
	"github.com/google/uuid"
	"log/slog"
)

type OrderRepo interface {
	Get(uid uuid.UUID) (order domain.Order, err error)
	Save(order domain.Order) (err error)
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
	return s.repo.Save(*order)
}
