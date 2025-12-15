package config

import (
	"context"
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

var Rdb *redis.Client

func init() {

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("failed to Read .env")
	}

	Rdb = redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})
	if res := Rdb.Ping(ctx); res.Err() != nil {
		log.Fatal("cannot ping redis")
	}

}
