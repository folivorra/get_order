package main

import (
	"context"
	"encoding/json"
	"github.com/folivorra/get_order/internal/domain"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"log"
)

func main() {
	ctx := context.Background()

	writer := kafka.Writer{
		Addr:         kafka.TCP("localhost:29092"),
		Topic:        "get_orders",
		RequiredAcks: kafka.RequireOne,
	}
	defer func() {
		e := writer.Close()
		if e != nil {
			log.Println(e)
		}
	}()

	for i := 0; i < 5; i++ {
		order := domain.Order{
			OrderUID:    uuid.New(),
			TrackNumber: "WBILMTESTTRACK",
			Entry:       "WBIL",
			Delivery: &domain.Delivery{
				Name:    "Test Testov",
				Phone:   "+9720000000",
				Zip:     "2639809",
				City:    "Kiryat Mozkin",
				Address: "Ploshad Mira 15",
				Region:  "Kraiot",
				Email:   "test@gmail.com",
			},
			Payment: &domain.Payment{
				Transaction:  "b563feb7b2b84b6test",
				RequestID:    "",
				Currency:     "USD",
				Provider:     "wbpay",
				Amount:       1817,
				PaymentDT:    1637907727,
				Bank:         "alpha",
				DeliveryCost: 1500,
				GoodsTotal:   317,
				CustomFee:    0,
			},
			Items: []domain.Item{
				{
					ChrtID:      9934930,
					TrackNumber: "WBILMTESTTRACK",
					Price:       453,
					RID:         "ab4219087a764ae0btest",
					Name:        "Mascaras",
					Sale:        30,
					Size:        "0",
					TotalPrice:  317,
					NmID:        2389212,
					Brand:       "Vivienne Sabo",
					Status:      202,
				},
			},
			Locale:            "en",
			InternalSignature: "",
			CustomerID:        "test",
			DeliveryService:   "meest",
			ShardKey:          "9",
			SmID:              99,
			DateCreated:       "2021-11-26T06:22:19Z",
			OofShard:          "1",
		}

		b, err := json.Marshal(order)
		if err != nil {
			log.Println("fail to marshal order")
		}

		err = writer.WriteMessages(ctx, kafka.Message{
			Value: b,
		})
		if err != nil {
			log.Println("fail to write messages", err)
		}
	}
}
