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
	"crypto/tls"
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
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	TLS      bool   `mapstructure:"tls"`
}

var redisConfig RedisConfig

func InitRedis() {
	config := viper.Sub("redis")
	if config == nil {
		return
	}

	config.SetDefault("password", "")
	config.SetDefault("tls", true)
	err := config.Unmarshal(&redisConfig)
	if err != nil {
		panic(err)
	}

	opts := &redis.Options{
		Addr:     fmt.Sprintf("%s:%d", redisConfig.Host, redisConfig.Port),
		Password: redisConfig.Password,
		DB:       0,
	}
	if redisConfig.TLS {
		opts.TLSConfig = &tls.Config{}
	}

	client := redis.NewClient(opts)
	if err := client.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}

	redisClient = make([]*redis.Client, 2)
	redisClient[0] = client
	redisClient[1] = client
}
