package domain

import "github.com/google/uuid"

type OrderItem struct {
	OrderItemUID uuid.UUID
	OrderUID     uuid.UUID
	ItemUID      uuid.UUID
	Item         *Item
	Price        int
	Sale         int
	TotalPrice   int
	Quantity     int
}
