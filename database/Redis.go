package database

import (
	"context"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
)

var Redis *redis.Client
var Ctx = context.Background()

func InitRedis() {
	Redis = redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})

	_, err := Redis.Ping(Ctx).Result()
	if err != nil {
		log.Fatal("فشل الاتصال بـ Redis")
	}
}