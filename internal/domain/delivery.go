package domain

import "github.com/google/uuid"

type Delivery struct {
	DeliveryUID uuid.UUID `json:"delivery_uid"`
	Name        string    `json:"name"`
	Phone       string    `json:"phone"`
	Zip         string    `json:"zip"`
	City        string    `json:"city"`
	Address     string    `json:"address"`
	Region      string    `json:"region"`
	Email       string    `json:"email"`
}
