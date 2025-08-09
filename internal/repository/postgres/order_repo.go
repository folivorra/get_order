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

func (pg *PgOrderRepo) Get(ctx context.Context, uid uuid.UUID) (order *domain.Order, err error) {
	return nil, nil
}

func (pg *PgOrderRepo) Save(ctx context.Context, order *domain.Order) (err error) {
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
			item.OrderUID,
			item.ChrtID,
			item.TrackNumber,
			item.Price,
			item.RID,
			item.Name,
			item.Sale,
			item.Size,
			item.TotalPrice,
			item.NmID,
			item.Brand,
			item.Status,
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
