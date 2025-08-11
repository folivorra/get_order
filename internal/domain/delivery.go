package domain

import "github.com/google/uuid"

type Delivery struct {
	DeliveryUID uuid.UUID
	Name        string
	Phone       string
	Zip         string
	City        string
	Address     string
	Region      string
	Email       string
}
