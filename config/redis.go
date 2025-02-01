package config

import (
	"context"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client
var Ctx = context.Background()

func InitRedis() {
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		log.Fatal("REDIS_ADDR environment variable is not set")
	}

	RedisClient = redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	// Test Redis connection
	_, err := RedisClient.Ping(RedisClient.Context()).Result()
	if err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}

	log.Println("Redis connection established successfully!")
}
