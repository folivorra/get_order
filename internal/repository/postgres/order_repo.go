package postgres

import (
	"context"
	"database/sql"
	"errors"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/folivorra/get_order/internal/config"
	"github.com/folivorra/get_order/internal/domain"
	"github.com/folivorra/get_order/internal/usecase"
	"github.com/google/uuid"
	"time"
)

var (
	ErrMaxRetryAttemptsExceeded = errors.New("max retry attempts exceeded")
)

type PgOrderRepo struct {
	db  *sql.DB
	cfg config.Config
}

var _ usecase.OrderRepo = (*PgOrderRepo)(nil)

func NewPgOrderRepo(db *sql.DB, cfg config.Config) *PgOrderRepo {
	return &PgOrderRepo{
		db:  db,
		cfg: cfg,
	}
}

func (pg *PgOrderRepo) Get(ctx context.Context, uid uuid.UUID) (*domain.Order, error) {
	var order domain.Order

	err := retry(ctx, pg.cfg.PgMaxRetries, pg.cfg.PgBackoff, func() error {
		funcCtx, cancel := context.WithTimeout(ctx, pg.cfg.PgGetTimeout)
		defer cancel()
		funcOrder := domain.Order{
			Delivery: domain.Delivery{},
			Payment:  domain.Payment{},
		}

		r, err := pg.db.QueryContext(funcCtx, orderGetQuery, uid)
		if err != nil {
			return err
		}
		defer func() {
			_ = r.Close()
		}()

		for r.Next() {
			item := domain.OrderItem{
				Item: &domain.Item{},
			}

			if err = r.Scan(
				&funcOrder.OrderUID,
				&funcOrder.TrackNumber,
				&funcOrder.Entry,
				&funcOrder.Delivery.DeliveryUID,
				&funcOrder.Payment.PaymentUID,
				&funcOrder.Locale,
				&funcOrder.InternalSignature,
				&funcOrder.CustomerID,
				&funcOrder.DeliveryService,
				&funcOrder.ShardKey,
				&funcOrder.SmID,
				&funcOrder.DateCreated,
				&funcOrder.OofShard,

				&funcOrder.Delivery.DeliveryUID,
				&funcOrder.Delivery.Name,
				&funcOrder.Delivery.Phone,
				&funcOrder.Delivery.Zip,
				&funcOrder.Delivery.City,
				&funcOrder.Delivery.Address,
				&funcOrder.Delivery.Region,
				&funcOrder.Delivery.Email,

				&funcOrder.Payment.PaymentUID,
				&funcOrder.Payment.Transaction,
				&funcOrder.Payment.RequestID,
				&funcOrder.Payment.Currency,
				&funcOrder.Payment.Provider,
				&funcOrder.Payment.Amount,
				&funcOrder.Payment.PaymentDT,
				&funcOrder.Payment.Bank,
				&funcOrder.Payment.DeliveryCost,
				&funcOrder.Payment.GoodsTotal,
				&funcOrder.Payment.CustomFee,

				&item.OrderItemUID,
				&item.OrderUID,
				&item.ItemUID,
				&item.Price,
				&item.Sale,
				&item.TotalPrice,
				&item.Quantity,

				&item.Item.ItemUID,
				&item.Item.ChrtID,
				&item.Item.TrackNumber,
				&item.Item.RID,
				&item.Item.Name,
				&item.Item.Size,
				&item.Item.NmID,
				&item.Item.Brand,
				&item.Item.Status,
			); err != nil {
				return err
			}

			funcOrder.Items = append(funcOrder.Items, item)
		}

		if err = r.Err(); err != nil {
			return err
		}

		order = funcOrder

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (pg *PgOrderRepo) Save(ctx context.Context, order *domain.Order) error {
	err := retry(ctx, pg.cfg.PgMaxRetries, pg.cfg.PgBackoff, func() error {
		funcCtx, cancel := context.WithTimeout(ctx, pg.cfg.PgSaveTimeout)
		defer cancel()

		tx, err := pg.db.BeginTx(funcCtx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
		if err != nil {
			return err
		}

		defer func() {
			_ = tx.Rollback()
		}()

		_, err = tx.ExecContext(funcCtx, deliverySaveQuery,
			order.Delivery.DeliveryUID,
			order.Delivery.Name,
			order.Delivery.Phone,
			order.Delivery.Zip,
			order.Delivery.City,
			order.Delivery.Address,
			order.Delivery.Region,
			order.Delivery.Email,
		)
		if err != nil {
			return err
		}

		_, err = tx.ExecContext(funcCtx, paymentSaveQuery,
			order.Payment.PaymentUID,
			order.Payment.Transaction,
			order.Payment.RequestID,
			order.Payment.Currency,
			order.Payment.Provider,
			order.Payment.Amount,
			order.Payment.PaymentDT,
			order.Payment.Bank,
			order.Payment.DeliveryCost,
			order.Payment.GoodsTotal,
			order.Payment.CustomFee,
		)
		if err != nil {
			return err
		}

		_, err = tx.ExecContext(funcCtx, orderSaveQuery,
			order.OrderUID,
			order.TrackNumber,
			order.Entry,
			order.Delivery.DeliveryUID,
			order.Payment.PaymentUID,
			order.Locale,
			order.InternalSignature,
			order.CustomerID,
			order.DeliveryService,
			order.ShardKey,
			order.SmID,
			order.DateCreated,
			order.OofShard,
		)
		if err != nil {
			return err
		}

		for _, item := range order.Items {
			_, err = tx.ExecContext(funcCtx, itemSaveQuery,
				item.ItemUID,
				item.Item.ChrtID,
				item.Item.TrackNumber,
				item.Item.RID,
				item.Item.Name,
				item.Item.Size,
				item.Item.NmID,
				item.Item.Brand,
				item.Item.Status,
			)
			if err != nil {
				return err
			}

			_, err = tx.ExecContext(funcCtx, itemOrderSaveQuery,
				item.OrderItemUID,
				item.OrderUID,
				item.ItemUID,
				item.Price,
				item.Sale,
				item.TotalPrice,
				item.Quantity,
			)
			if err != nil {
				return err
			}
		}

		return tx.Commit()
	})

	return err
}

func (pg *PgOrderRepo) Exists(ctx context.Context, uuid uuid.UUID) (exists bool, err error) {
	err = retry(ctx, pg.cfg.PgMaxRetries, pg.cfg.PgBackoff, func() error {
		funcCtx, cancel := context.WithTimeout(ctx, pg.cfg.PgExistsTimeout)
		defer cancel()
		return pg.db.QueryRowContext(funcCtx, existsCheckQuery, uuid).Scan(&exists)
	})

	return exists, err
}

func retry(ctx context.Context, maxRetries int, backoff time.Duration, fn func() error) error {
	var err error

	for attempt := 0; attempt < maxRetries; attempt++ {
		err = fn()
		if err == nil || errors.Is(err, sql.ErrNoRows) {
			return err
		}

		jitter := time.Duration(gofakeit.IntN(500)) * time.Millisecond

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(backoff + jitter):
		}
	}

	return ErrMaxRetryAttemptsExceeded
}
