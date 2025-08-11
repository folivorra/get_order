package domain

import (
	"github.com/google/uuid"
)

type Order struct {
	OrderUID          uuid.UUID
	TrackNumber       string
	Entry             string
	Delivery          *Delivery
	Payment           *Payment
	Items             []*OrderItem
	Locale            string
	InternalSignature string
	CustomerID        string
	DeliveryService   string
	ShardKey          string
	SmID              int
	DateCreated       string
	OofShard          string
}
