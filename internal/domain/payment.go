package domain

import "github.com/google/uuid"

type Payment struct {
	PaymentUID   uuid.UUID
	Transaction  string
	RequestID    string
	Currency     string
	Provider     string
	Amount       int
	PaymentDT    int
	Bank         string
	DeliveryCost int
	GoodsTotal   int
	CustomFee    int
}
