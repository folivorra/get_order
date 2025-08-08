package postgres

import (
	"github.com/folivorra/get_order/internal/domain"
	"github.com/google/uuid"
)

type PgOrderRepo struct {
}

func NewPgOrderRepo() *PgOrderRepo {
	return &PgOrderRepo{}
}

func (pg *PgOrderRepo) Get(uid uuid.UUID) (order domain.Order, err error) {
	return domain.Order{}, nil
}

func (pg *PgOrderRepo) Save(order domain.Order) (err error) {
	return nil
}
