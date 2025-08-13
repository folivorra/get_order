package mapper

import (
	"github.com/folivorra/get_order/internal/domain"
	"github.com/google/uuid"
)

type DeliveryIntoDomainDTO struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email"`
}

type PaymentIntoDomainDTO struct {
	Transaction  string `json:"transaction"`
	RequestID    string `json:"request_id"`
	Currency     string `json:"currency"`
	Provider     string `json:"provider"`
	Amount       int    `json:"amount"`
	PaymentDT    int    `json:"payment_dt"`
	Bank         string `json:"bank"`
	DeliveryCost int    `json:"delivery_cost"`
	GoodsTotal   int    `json:"goods_total"`
	CustomFee    int    `json:"custom_fee"`
}

type ItemIntoDomainDTO struct {
	ItemUID     uuid.UUID `json:"item_uid"`
	ChrtID      int       `json:"chrt_id"`
	TrackNumber string    `json:"track_number"`
	Price       int       `json:"price"`
	Rid         string    `json:"rid"`
	Name        string    `json:"name"`
	Sale        int       `json:"sale"`
	Size        string    `json:"size"`
	TotalPrice  int       `json:"total_price"`
	NmID        int       `json:"nm_id"`
	Brand       string    `json:"brand"`
	Status      int       `json:"status"`
	Quantity    int       `json:"quantity"`
}

type OrderIntoDomainDTO struct {
	OrderUID          uuid.UUID             `json:"order_uid"`
	TrackNumber       string                `json:"track_number"`
	Entry             string                `json:"entry"`
	Delivery          DeliveryIntoDomainDTO `json:"delivery"`
	Payment           PaymentIntoDomainDTO  `json:"payment"`
	Items             []ItemIntoDomainDTO   `json:"items"`
	Locale            string                `json:"locale"`
	InternalSignature string                `json:"internal_signature"`
	CustomerID        string                `json:"customer_id"`
	DeliveryService   string                `json:"delivery_service"`
	Shardkey          string                `json:"shardkey"`
	SmID              int                   `json:"sm_id"`
	DateCreated       string                `json:"date_created"`
	OofShard          string                `json:"oof_shard"`
}

func ConvertToDomain(dto *OrderIntoDomainDTO) *domain.Order {
	items := make([]domain.OrderItem, len(dto.Items))
	for i, itemDTO := range dto.Items {
		items[i] = domain.OrderItem{
			OrderUID:   dto.OrderUID,
			ItemUID:    itemDTO.ItemUID,
			Price:      itemDTO.Price,
			Sale:       itemDTO.Sale,
			TotalPrice: itemDTO.TotalPrice,
			Quantity:   itemDTO.Quantity,
			Item: &domain.Item{
				ItemUID:     itemDTO.ItemUID,
				ChrtID:      itemDTO.ChrtID,
				TrackNumber: itemDTO.TrackNumber,
				RID:         itemDTO.Rid,
				Name:        itemDTO.Name,
				Size:        itemDTO.Size,
				NmID:        itemDTO.NmID,
				Brand:       itemDTO.Brand,
				Status:      itemDTO.Status,
			},
		}
	}

	return &domain.Order{
		OrderUID:    dto.OrderUID,
		TrackNumber: dto.TrackNumber,
		Entry:       dto.Entry,
		Delivery: domain.Delivery{
			Name:    dto.Delivery.Name,
			Phone:   dto.Delivery.Phone,
			Zip:     dto.Delivery.Zip,
			City:    dto.Delivery.City,
			Address: dto.Delivery.Address,
			Region:  dto.Delivery.Region,
			Email:   dto.Delivery.Email,
		},
		Payment: domain.Payment{
			Transaction:  dto.Payment.Transaction,
			RequestID:    dto.Payment.RequestID,
			Currency:     dto.Payment.Currency,
			Provider:     dto.Payment.Provider,
			Amount:       dto.Payment.Amount,
			PaymentDT:    dto.Payment.PaymentDT,
			Bank:         dto.Payment.Bank,
			DeliveryCost: dto.Payment.DeliveryCost,
			GoodsTotal:   dto.Payment.GoodsTotal,
			CustomFee:    dto.Payment.CustomFee,
		},
		Items:             items,
		Locale:            dto.Locale,
		InternalSignature: dto.InternalSignature,
		CustomerID:        dto.CustomerID,
		DeliveryService:   dto.DeliveryService,
		ShardKey:          dto.Shardkey,
		SmID:              dto.SmID,
		DateCreated:       dto.DateCreated,
		OofShard:          dto.OofShard,
	}
}
