package database

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

var redisClient []*redis.Client

func GetRedis() *redis.Client {
	return redisClient[0]
}

func GetSessionRedis() *redis.Client {
	return redisClient[1]
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
	redisClient = make([]*redis.Client, 2)
	for i := 0; i < 2; i++ {
		redisClient[i] = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", redisConfig.Host, redisConfig.Port),
			Password: redisConfig.Password,
			DB:       i,
		})
		if err := redisClient[i].Ping(context.Background()).Err(); err != nil {
			panic(err)
		}
		if err := redisClient[i].ConfigSet(context.Background(), "maxmemory", redisConfig.MaxMemory).Err(); err != nil {
			panic(err)
		}
		if err := redisClient[i].ConfigSet(context.Background(), "maxmemory-policy", redisConfig.MaxMemoryPolicy).Err(); err != nil {
			panic(err)
		}
	}
}
