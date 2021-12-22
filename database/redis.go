package database

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"os"
)

var redisClient *redis.Client = nil

func GetRedis() *redis.Client {
	return redisClient
}

func init() {
	host := os.Getenv("REDIS_SERVER_HOST")
	port := os.Getenv("REDIS_SERVER_PORT")
	password := os.Getenv("REDIS_SERVER_PASSWORD")
	redisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: password,
		DB:       0,
	})
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}
}
