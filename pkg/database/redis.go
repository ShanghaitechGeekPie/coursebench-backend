// Copyright (C) 2021-2024 ShanghaiTech GeekPie
// This file is part of CourseBench Backend.
//
// CourseBench Backend is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// CourseBench Backend is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with CourseBench Backend.  If not, see <http://www.gnu.org/licenses/>.

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
