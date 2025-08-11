package postgres

import (
	"context"
	"database/sql"
	"github.com/folivorra/get_order/internal/config"
	"github.com/folivorra/get_order/internal/domain"
	"github.com/folivorra/get_order/internal/usecase"
	"github.com/google/uuid"
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
	ctx, cancel := context.WithTimeout(ctx, pg.cfg.GetTimeout)
	defer cancel()

	order := domain.Order{
		Delivery: &domain.Delivery{},
		Payment:  &domain.Payment{},
	}

	r, err := pg.db.QueryContext(ctx, orderGetQuery, uid)
	if err != nil {
		return nil, err
	}
	defer func(r *sql.Rows) {
		err = r.Close()
		if err != nil {
			return
		}
	}(r)

	for r.Next() {
		item := domain.OrderItem{
			Item: &domain.Item{},
		}

		if err = r.Scan(
			&order.OrderUID,
			&order.TrackNumber,
			&order.Entry,
			&order.Delivery.DeliveryUID,
			&order.Payment.PaymentUID,
			&order.Locale,
			&order.InternalSignature,
			&order.CustomerID,
			&order.DeliveryService,
			&order.ShardKey,
			&order.SmID,
			&order.DateCreated,
			&order.OofShard,

			&order.Delivery.DeliveryUID,
			&order.Delivery.Name,
			&order.Delivery.Phone,
			&order.Delivery.Zip,
			&order.Delivery.City,
			&order.Delivery.Address,
			&order.Delivery.Region,
			&order.Delivery.Email,

			&order.Payment.PaymentUID,
			&order.Payment.Transaction,
			&order.Payment.RequestID,
			&order.Payment.Currency,
			&order.Payment.Provider,
			&order.Payment.Amount,
			&order.Payment.PaymentDT,
			&order.Payment.Bank,
			&order.Payment.DeliveryCost,
			&order.Payment.GoodsTotal,
			&order.Payment.CustomFee,

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
			return nil, err
		}

		order.Items = append(order.Items, &item)
	}

	return &order, r.Err()
}

func (pg *PgOrderRepo) Save(ctx context.Context, order *domain.Order) error {
	ctx, cancel := context.WithTimeout(ctx, pg.cfg.SaveTimeout)
	defer cancel()

	tx, err := pg.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, deliverySaveQuery,
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
		_ = tx.Rollback()
		return err
	}

	_, err = tx.ExecContext(ctx, paymentSaveQuery,
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
		_ = tx.Rollback()
		return err
	}

	_, err = tx.ExecContext(ctx, orderSaveQuery,
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
		_ = tx.Rollback()
		return err
	}

	for _, item := range order.Items {
		_, err = tx.ExecContext(ctx, itemSaveQuery,
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
			_ = tx.Rollback()
			return err
		}

		_, err = tx.ExecContext(ctx, itemOrderSaveQuery,
			item.OrderItemUID,
			item.OrderUID,
			item.ItemUID,
			item.Price,
			item.Sale,
			item.TotalPrice,
			item.Quantity,
		)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (pg *PgOrderRepo) Exists(ctx context.Context, uuid uuid.UUID) (exists bool, err error) {
	ctx, cancel := context.WithTimeout(ctx, pg.cfg.ExistsTimeout)
	defer cancel()

	err = pg.db.QueryRowContext(ctx, existsCheckQuery, uuid).Scan(&exists)

	return exists, err
}
