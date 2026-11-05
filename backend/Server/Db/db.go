package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func InitDB() *pgxpool.Pool {
    connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
        os.Getenv("DB_USER"),
        os.Getenv("DB_PASSWORD"),
        os.Getenv("DB_HOST"),
        os.Getenv("DB_PORT"),
        os.Getenv("DB_NAME"),
    )

    config, err := pgxpool.ParseConfig(connStr)
    if err != nil {
        log.Fatal(err)
    }

    db, err := pgxpool.NewWithConfig(context.Background(), config)
    if err != nil {
        log.Fatal(err)
    }

    if err := db.Ping(context.Background()); err != nil {
        db.Close()
        log.Fatal(err)
    }

    log.Println("Connected to PostgreSQL")

    return db
}