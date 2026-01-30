package database

import (
	"context"
	"database/sql"
	"log"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func InitDB(conString string) (*sql.DB, error) {
	db, err := sql.Open("pgx", conString)
	if err != nil {
		return nil, err
	}

	if err := db.PingContext(context.Background()); err != nil {
		_ = db.Close()
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetConnMaxIdleTime(5 * time.Minute)

	log.Println("Database connection established successfully")

	return db, nil
}
