package db

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
)




func InitRedis() *redis.Client {
    redisAddr := os.Getenv("REDIS_ADDR")
    redisPassword := os.Getenv("REDIS_PASSWORD")

    redisDB, err := strconv.Atoi(os.Getenv("REDIS_DB"))
    if err != nil {
        log.Fatal("Invalid REDIS_DB:", err)
    }

    rdb := redis.NewClient(&redis.Options{
        Addr:     redisAddr,
        Password: redisPassword,
        DB:       redisDB,
    })

    ctx := context.Background()

    if err := rdb.Ping(ctx).Err(); err != nil {
        log.Fatal("Redis connection failed:", err)
    }

    log.Println("Redis connected successfully")

    return rdb
}