package config

import (
	"context"
	"log"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/knadh/koanf/v2"
	"github.com/mcuadros/go-defaults"
	"github.com/redis/go-redis/v9"
)

var Rdb *redis.Client
var C Config

func init() {

	k := koanf.New(".")
	defaults.SetDefaults(&C)
	
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("failed to Read .env")
	}
	unmarshalerConfig := koanf.UnmarshalConf{Tag: "json"}
	if err := k.UnmarshalWithConf("", &C, unmarshalerConfig); err != nil {
		log.Fatal(err)
	}

	v := validator.New()
	if err := v.Struct(C); err != nil {
		log.Fatal(err)
	}
	Rdb = redis.NewClient(&redis.Options{
		Addr: 	C.Redis.Addr,
		Password: C.Redis.Password,
		DB:       C.Redis.DB,
	})
	if res := Rdb.Ping(ctx); res.Err() != nil {
		log.Fatal("cannot ping redis", res.Err())
	}

}
