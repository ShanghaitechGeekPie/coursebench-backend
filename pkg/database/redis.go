package database

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

var redisClient *redis.Client = nil

func GetRedis() *redis.Client {
	return redisClient
}

type RedisConfig struct {
	Host            string `mapstructure:"host"`
	Port            int    `mapstructure:"port"`
	Password        string `mapstructure:"password"`
	MaxMemory       string `mapstructure:"max_memory"`
	MaxMemoryPolicy string `mapstructure:"max_memory_policy"`
}

var redisConfig RedisConfig

func InitRedis() {
	config := viper.Sub("redis")
	if config == nil {
		return
	}

	config.SetDefault("password", "")
	config.SetDefault("max_memory", "32GB")
	config.SetDefault("max_memory_policy", "volatile-lru")
	err := config.Unmarshal(&redisConfig)
	if err != nil {
		panic(err)
	}
	redisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", redisConfig.Host, redisConfig.Port),
		Password: redisConfig.Password,
		DB:       0,
	})
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}
	if err := redisClient.ConfigSet(context.Background(), "maxmemory", redisConfig.MaxMemory).Err(); err != nil {
		panic(err)
	}
	if err := redisClient.ConfigSet(context.Background(), "maxmemory-policy", redisConfig.MaxMemoryPolicy).Err(); err != nil {
		panic(err)
	}
}
