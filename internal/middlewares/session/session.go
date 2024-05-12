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

package session

import (
	"context"
	"coursebench-backend/internal/config"
	"coursebench-backend/pkg/database"
	"coursebench-backend/pkg/errors"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"time"
)

var store *session.Store

type RedisStore struct {
	db *redis.Client
}

func (r *RedisStore) Delete(key string) error {
	return r.db.Del(context.Background(), key).Err()
}

func (r *RedisStore) Reset() error {
	return r.db.FlushDB(context.Background()).Err()
}

func (r *RedisStore) Close() error {
	return r.db.Close()
}

func (r *RedisStore) Get(key string) ([]byte, error) {
	v, err := r.db.Get(context.Background(), key).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	return v, err
}

func (r *RedisStore) Set(key string, value []byte, ttl time.Duration) error {
	return r.db.Set(context.Background(), key, value, ttl).Err()
}

func Init() {
	redis := database.GetSessionRedis()
	store = session.New(session.Config{Expiration: time.Hour * 24 * 30, CookieHTTPOnly: !config.GlobalConf.InDevelopment, CookieSecure: !config.GlobalConf.InDevelopment, Storage: &RedisStore{db: redis}})
}

func GetStore() *session.Store {
	return store
}

func GetUserID(ctx *fiber.Ctx) (uint, error) {
	sess, err := store.Get(ctx)
	if err != nil {
		return 0, errors.Wrap(err, errors.InternalServerError)
	}
	t := sess.Get("user_id")
	id, ok := t.(uint)
	if !ok {
		return 0, errors.New(errors.UserNotLogin)
	}
	return id, nil
}

func GetSession(ctx *fiber.Ctx) (*session.Session, error) {
	return store.Get(ctx)
}
