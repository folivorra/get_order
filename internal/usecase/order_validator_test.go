package usecase_test

import (
	"testing"
	"time"

	"github.com/folivorra/get_order/internal/adapter/mapper"
	"github.com/folivorra/get_order/internal/usecase"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestValidateOrder_ValidOrder(t *testing.T) {
	order := &mapper.OrderIntoDomainDTO{
		OrderUID:    uuid.New(),
		TrackNumber: "TRACK123",
		Delivery: mapper.DeliveryIntoDomainDTO{
			Name: "Test User",
			City: "Test City",
		},
		Payment: mapper.PaymentIntoDomainDTO{
			Amount: 100,
		},
		Items: []mapper.ItemIntoDomainDTO{
			{
				NmID:       1,
				TotalPrice: 100,
			},
		},
		DateCreated: time.Now().Format(time.RFC3339),
	}

	err := usecase.ValidateOrder(order)
	assert.NoError(t, err)
}

func TestValidateOrder_OrderUIDEmpty(t *testing.T) {
	order := &mapper.OrderIntoDomainDTO{}
	err := usecase.ValidateOrder(order)
	assert.Equal(t, usecase.ErrOrderUIDIsEmpty, err)
}

func TestValidateOrder_TrackNumberEmpty(t *testing.T) {
	order := &mapper.OrderIntoDomainDTO{
		OrderUID: uuid.New(),
	}
	err := usecase.ValidateOrder(order)
	assert.Equal(t, usecase.ErrTrackNumberIsEmpty, err)
}

func TestValidateOrder_DeliveryIncomplete(t *testing.T) {
	order := &mapper.OrderIntoDomainDTO{
		OrderUID:    uuid.New(),
		TrackNumber: "TRACK123",
		Delivery:    mapper.DeliveryIntoDomainDTO{},
		Payment:     mapper.PaymentIntoDomainDTO{Amount: 10},
		Items:       []mapper.ItemIntoDomainDTO{{NmID: 1, TotalPrice: 10}},
		DateCreated: time.Now().Format(time.RFC3339),
	}
	err := usecase.ValidateOrder(order)
	assert.Equal(t, usecase.ErrDeliveryInfoIncomplete, err)
}

func TestValidateOrder_PaymentInvalid(t *testing.T) {
	order := &mapper.OrderIntoDomainDTO{
		OrderUID:    uuid.New(),
		TrackNumber: "TRACK123",
		Delivery:    mapper.DeliveryIntoDomainDTO{Name: "A", City: "B"},
		Payment:     mapper.PaymentIntoDomainDTO{Amount: 0},
		Items:       []mapper.ItemIntoDomainDTO{{NmID: 1, TotalPrice: 10}},
		DateCreated: time.Now().Format(time.RFC3339),
	}
	err := usecase.ValidateOrder(order)
	assert.Equal(t, usecase.ErrPaymentAmountInvalid, err)
}

func TestValidateOrder_ItemsInvalid(t *testing.T) {
	order := &mapper.OrderIntoDomainDTO{
		OrderUID:    uuid.New(),
		TrackNumber: "TRACK123",
		Delivery:    mapper.DeliveryIntoDomainDTO{Name: "A", City: "B"},
		Payment:     mapper.PaymentIntoDomainDTO{Amount: 10},
		Items:       []mapper.ItemIntoDomainDTO{},
		DateCreated: time.Now().Format(time.RFC3339),
	}
	err := usecase.ValidateOrder(order)
	assert.Equal(t, usecase.ErrItemsListIsEmpty, err)
}

func TestValidateOrder_ItemFieldsInvalid(t *testing.T) {
	order := &mapper.OrderIntoDomainDTO{
		OrderUID:    uuid.New(),
		TrackNumber: "TRACK123",
		Delivery:    mapper.DeliveryIntoDomainDTO{Name: "A", City: "B"},
		Payment:     mapper.PaymentIntoDomainDTO{Amount: 10},
		Items:       []mapper.ItemIntoDomainDTO{{NmID: 0, TotalPrice: 0}},
		DateCreated: time.Now().Format(time.RFC3339),
	}
	err := usecase.ValidateOrder(order)
	assert.Equal(t, usecase.ErrItemNmIdIsInvalid, err)
}

func TestValidateOrder_DateCreatedInvalid(t *testing.T) {
	order := &mapper.OrderIntoDomainDTO{
		OrderUID:    uuid.New(),
		TrackNumber: "TRACK123",
		Delivery:    mapper.DeliveryIntoDomainDTO{Name: "A", City: "B"},
		Payment:     mapper.PaymentIntoDomainDTO{Amount: 10},
		Items:       []mapper.ItemIntoDomainDTO{{NmID: 1, TotalPrice: 10}},
		DateCreated: "invalid-date",
	}
	err := usecase.ValidateOrder(order)
	assert.Equal(t, usecase.ErrDateCreatedIsInvalid, err)
}
