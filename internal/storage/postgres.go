package storage

import (
	"context"
	"database/sql"
	"github.com/folivorra/get_order/internal/config"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
)

func NewPgClient(ctx context.Context, cfg config.Config) *sql.DB {
	db, err := sql.Open("pgx", cfg.PgDsn)
	if err != nil {
		log.Fatal(err)
	}

	db.SetMaxOpenConns(cfg.PgPoolMaxOpenConns)       // максимум 10 одновременных открытых соединений
	db.SetMaxIdleConns(cfg.PgPoolMaxIdleConns)       // максимум 5 соединений в простое
	db.SetConnMaxLifetime(cfg.PgPoolConnMaxLifetime) // таймаут

	timeout, cancel := context.WithTimeout(ctx, cfg.PgPingTimeout)
	defer cancel()

	if err = db.PingContext(timeout); err != nil {
		log.Fatal(err)
	}

	//app.RegisterCleanup(func(ctx context.Context) {
	//	timeout, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	//	defer cancel()
	//
	//	if err := db.PingContext(timeout); err != nil {
	//		logger.ErrorLogger.Println("postgres connection error: %v", err)
	//	} else if err := db.Close(); err != nil {
	//		logger.ErrorLogger.Println("postgres close error: %v", err)
	//	}
	//}) todo

	return db
}
