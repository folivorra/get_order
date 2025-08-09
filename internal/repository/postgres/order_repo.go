package postgres

import (
	"context"
	"database/sql"
	"github.com/folivorra/get_order/internal/domain"
	"github.com/folivorra/get_order/internal/usecase"
	"github.com/google/uuid"
)

type PgOrderRepo struct {
	db *sql.DB
}

var _ usecase.OrderRepo = (*PgOrderRepo)(nil)

func NewPgOrderRepo(db *sql.DB) *PgOrderRepo {
	return &PgOrderRepo{
		db: db,
	}
}

func (pg *PgOrderRepo) Get(uid uuid.UUID) (order domain.Order, err error) {
	return domain.Order{}, nil
}

func (pg *PgOrderRepo) Save(order domain.Order) (err error) {
	// todo
	return nil
}

func (pg *PgOrderRepo) Exists(ctx context.Context, uuid uuid.UUID) (exists bool, err error) {
	const query string = `SELECT EXISTS (SELECT id FROM orders WHERE id = $1)`

	err = pg.db.QueryRowContext(ctx, query, uuid).Scan(&exists)

	return exists, err
}
