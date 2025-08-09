package storage

import (
	"context"
	"database/sql"
	"log"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func NewPgClient(ctx context.Context, dsn string) *sql.DB {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal(err)
	}

	db.SetMaxOpenConns(10)           // максимум 10 одновременных открытых соединений
	db.SetMaxIdleConns(5)            // максимум 5 соединений в простое
	db.SetConnMaxLifetime(time.Hour) // таймаут

	timeout, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
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
