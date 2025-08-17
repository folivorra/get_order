package usecase

import (
	"errors"
	"github.com/folivorra/get_order/internal/adapter/mapper"
	"github.com/google/uuid"
	"time"
)

var (
	ErrOrderUIDIsEmpty         = errors.New("order_uid is empty")
	ErrTrackNumberIsEmpty      = errors.New("track_number is empty")
	ErrDeliveryInfoIncomplete  = errors.New("delivery info is incomplete")
	ErrPaymentAmountInvalid    = errors.New("payment amount is invalid")
	ErrItemsListIsEmpty        = errors.New("items list is empty")
	ErrItemNmIdIsInvalid       = errors.New("nm_id is invalid")
	ErrItemTotalPriceIsInvalid = errors.New("total_price is invalid")
	ErrDateCreatedIsInvalid    = errors.New("date_created is invalid")
)

func ValidateOrder(order *mapper.OrderIntoDomainDTO) error {
	if order.OrderUID == uuid.Nil {
		return ErrOrderUIDIsEmpty
	}
	if order.TrackNumber == "" {
		return ErrTrackNumberIsEmpty
	}
	if order.Delivery.Name == "" || order.Delivery.City == "" {
		return ErrDeliveryInfoIncomplete
	}
	if order.Payment.Amount <= 0 {
		return ErrPaymentAmountInvalid
	}

	if len(order.Items) == 0 {
		return ErrItemsListIsEmpty
	}
	for _, item := range order.Items {
		if item.NmID <= 0 {
			return ErrItemNmIdIsInvalid
		}
		if item.TotalPrice <= 0 {
			return ErrItemTotalPriceIsInvalid
		}
	}

	if _, err := time.Parse(time.RFC3339, order.DateCreated); err != nil {
		return ErrDateCreatedIsInvalid
	}

	return nil
}
