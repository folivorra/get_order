package main

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/folivorra/get_order/internal/domain"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

func main() {
	_ = gofakeit.Seed(0)
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

	wg := sync.WaitGroup{}
	wg.Add(5)
	for i := 0; i < 5; i++ {
		go func() {
			defer wg.Done()
			orderUID := uuid.New()
			order := domain.Order{
				OrderUID:    orderUID,
				TrackNumber: uuid.New().String(),
				Entry:       gofakeit.LetterN(5),
				Delivery: &domain.Delivery{
					Name:    gofakeit.Name(),
					Phone:   gofakeit.Phone(),
					Zip:     gofakeit.Zip(),
					City:    gofakeit.City(),
					Address: gofakeit.Email(),
					Region:  gofakeit.State(),
					Email:   gofakeit.Email(),
				},
				Payment: &domain.Payment{
					Transaction:  uuid.New().String(),
					RequestID:    gofakeit.DigitN(10),
					Currency:     gofakeit.CurrencyShort(),
					Provider:     gofakeit.CreditCardType(),
					Amount:       gofakeit.Number(1, 100),
					PaymentDT:    gofakeit.Number(1, 100),
					Bank:         gofakeit.BankType(),
					DeliveryCost: gofakeit.Number(1, 100),
					GoodsTotal:   gofakeit.Number(1, 100),
					CustomFee:    gofakeit.Number(1, 100),
				},
				Items: []domain.Item{
					{
						ChrtID:      gofakeit.Number(1, 100),
						TrackNumber: uuid.New().String(),
						Price:       gofakeit.Number(1, 100),
						RID:         uuid.New().String(),
						Name:        gofakeit.ProductName(),
						Sale:        gofakeit.Number(1, 100),
						Size:        gofakeit.DigitN(2),
						TotalPrice:  gofakeit.Number(1, 100),
						NmID:        gofakeit.Number(1, 100),
						Brand:       gofakeit.Company(),
						Status:      gofakeit.HTTPStatusCode(),
					},
				},
				Locale:            gofakeit.LanguageAbbreviation(),
				InternalSignature: gofakeit.LetterN(5),
				CustomerID:        uuid.New().String(),
				DeliveryService:   gofakeit.LetterN(5),
				ShardKey:          gofakeit.Digit(),
				SmID:              gofakeit.Number(1, 100),
				DateCreated:       gofakeit.PastDate().Format(time.RFC3339),
				OofShard:          gofakeit.Digit(),
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
		}()
	}
	wg.Wait()
}
