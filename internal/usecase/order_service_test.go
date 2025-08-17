package usecase_test

import (
	"context"
	"errors"
	"github.com/folivorra/get_order/internal/config"
	"log/slog"
	"os"
	"testing"

	"github.com/folivorra/get_order/internal/domain"
	"github.com/folivorra/get_order/internal/usecase"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) Get(ctx context.Context, uid uuid.UUID) (*domain.Order, error) {
	args := m.Called(ctx, uid)
	return args.Get(0).(*domain.Order), args.Error(1)
}

func (m *MockRepo) Save(ctx context.Context, order *domain.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

func (m *MockRepo) GetLastN(ctx context.Context, n int) ([]*domain.Order, error) {
	args := m.Called(ctx, n)
	return args.Get(0).([]*domain.Order), args.Error(1)
}

type MockCache struct {
	mock.Mock
}

func (m *MockCache) Get(uid uuid.UUID) (*domain.Order, error) {
	args := m.Called(uid)
	return args.Get(0).(*domain.Order), args.Error(1)
}

func (m *MockCache) Set(order *domain.Order) {
	m.Called(order)
}

func TestProcessIncomingOrder_SavesOrder(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepo)
	cache := new(MockCache)
	logger := slog.New(
		slog.NewTextHandler(
			os.Stdout, &slog.HandlerOptions{
				Level:     slog.LevelDebug,
				AddSource: true,
			},
		),
	)
	cfg := config.NewConfig(logger)
	service := usecase.NewOrderService(logger, cfg, repo, cache)

	order := &domain.Order{
		Items:    []domain.OrderItem{{}},
		Delivery: domain.Delivery{},
		Payment:  domain.Payment{},
	}

	repo.On("Save", ctx, mock.AnythingOfType("*domain.Order")).Return(nil)

	err := service.ProcessIncomingOrder(ctx, order)

	assert.NoError(t, err)
	repo.AssertCalled(t, "Save", ctx, mock.AnythingOfType("*domain.Order"))
	assert.NotEqual(t, uuid.Nil, order.Delivery.DeliveryUID)
	assert.NotEqual(t, uuid.Nil, order.Payment.PaymentUID)
	assert.NotEqual(t, uuid.Nil, order.Items[0].OrderItemUID)
}

func TestGetOrder_CacheHit(t *testing.T) {
	uid := uuid.New()
	order := &domain.Order{OrderUID: uid}
	repo := new(MockRepo)
	cache := new(MockCache)
	logger := slog.New(
		slog.NewTextHandler(
			os.Stdout, &slog.HandlerOptions{
				Level:     slog.LevelDebug,
				AddSource: true,
			},
		),
	)
	cfg := config.NewConfig(logger)
	service := usecase.NewOrderService(logger, cfg, repo, cache)

	cache.On("Get", uid).Return(order, nil)

	got, err := service.GetOrder(context.Background(), uid)

	assert.NoError(t, err)
	assert.Equal(t, order, got)
	repo.AssertNotCalled(t, "Get", mock.Anything, mock.Anything)
}

func TestGetOrder_CacheMiss(t *testing.T) {
	ctx := context.Background()
	uid := uuid.New()
	order := &domain.Order{OrderUID: uid}

	repo := new(MockRepo)
	cache := new(MockCache)
	logger := slog.New(
		slog.NewTextHandler(
			os.Stdout, &slog.HandlerOptions{
				Level:     slog.LevelDebug,
				AddSource: true,
			},
		),
	)
	cfg := config.NewConfig(logger)
	service := usecase.NewOrderService(logger, cfg, repo, cache)

	cache.On("Get", uid).Return(&domain.Order{}, errors.New("not found"))
	repo.On("Get", ctx, uid).Return(order, nil)
	cache.On("Set", order).Return()

	got, err := service.GetOrder(ctx, uid)

	assert.NoError(t, err)
	assert.Equal(t, order, got)
	cache.AssertCalled(t, "Set", order)
}

func TestWarmUpCache_SetsOrders(t *testing.T) {
	ctx := context.Background()
	orders := []*domain.Order{
		{OrderUID: uuid.New()},
		{OrderUID: uuid.New()},
	}

	repo := new(MockRepo)
	cache := new(MockCache)
	logger := slog.New(
		slog.NewTextHandler(
			os.Stdout, &slog.HandlerOptions{
				Level:     slog.LevelDebug,
				AddSource: true,
			},
		),
	)
	cfg := config.NewConfig(logger)
	service := usecase.NewOrderService(logger, cfg, repo, cache)

	repo.On("GetLastN", ctx, 2).Return(orders, nil)
	cache.On("Set", orders[0]).Return()
	cache.On("Set", orders[1]).Return()

	err := service.WarmUpCache(ctx, 2)

	assert.NoError(t, err)
	cache.AssertCalled(t, "Set", orders[0])
	cache.AssertCalled(t, "Set", orders[1])
}
