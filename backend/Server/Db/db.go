package db

import (
    "context"
    "log"

    "github.com/jackc/pgx/v5/pgxpool"
)

func InitDB() *pgxpool.Pool {
    connStr := "postgres://postgres:password@localhost:5432/mydb"

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