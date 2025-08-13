package mapper

import (
	"github.com/folivorra/get_order/internal/domain"
	"github.com/google/uuid"
)

type DeliveryFromDomainDTO struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email"`
}

type PaymentFromDomainDTO struct {
	Currency     string `json:"currency"`
	Amount       int    `json:"amount"`
	DeliveryCost int    `json:"delivery_cost"`
}

type ItemFromDomainDTO struct {
	Price      int    `json:"price"`
	Name       string `json:"name"`
	Sale       int    `json:"sale"`
	Size       string `json:"size"`
	TotalPrice int    `json:"total_price"`
	Brand      string `json:"brand"`
	Quantity   int    `json:"quantity"`
}

type OrderFromDomainDTO struct {
	OrderUID        uuid.UUID             `json:"order_uid"`
	TrackNumber     string                `json:"track_number"`
	Delivery        DeliveryFromDomainDTO `json:"delivery"`
	Payment         PaymentFromDomainDTO  `json:"payment"`
	Items           []ItemFromDomainDTO   `json:"items"`
	DeliveryService string                `json:"delivery_service"`
	DateCreated     string                `json:"date_created"`
}

func ConvertFromDomain(order *domain.Order) *OrderFromDomainDTO {
	orderDTO := OrderFromDomainDTO{
		OrderUID:    order.OrderUID,
		TrackNumber: order.TrackNumber,
		Delivery: DeliveryFromDomainDTO{
			Name:    order.Delivery.Name,
			Phone:   order.Delivery.Phone,
			Zip:     order.Delivery.Zip,
			City:    order.Delivery.City,
			Address: order.Delivery.Address,
			Region:  order.Delivery.Region,
			Email:   order.Delivery.Email,
		},
		Payment: PaymentFromDomainDTO{
			Currency:     order.Payment.Currency,
			Amount:       order.Payment.Amount,
			DeliveryCost: order.Payment.DeliveryCost,
		},
		Items:           make([]ItemFromDomainDTO, len(order.Items)),
		DeliveryService: order.DeliveryService,
		DateCreated:     order.DateCreated,
	}

	for i, item := range order.Items {
		orderDTO.Items[i] = ItemFromDomainDTO{
			Price:      item.Price,
			Name:       item.Item.Name,
			Sale:       item.Sale,
			Size:       item.Item.Size,
			TotalPrice: item.TotalPrice,
			Brand:      item.Item.Brand,
			Quantity:   item.Quantity,
		}
	}

	return &orderDTO
}
