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

package mail

import (
	"context"
	"coursebench-backend/internal/config"
	"coursebench-backend/pkg/database"
	"coursebench-backend/pkg/errors"
	"coursebench-backend/pkg/models"
	"fmt"
	"github.com/google/uuid"
	"time"
)

// PostMail 用户注册或找回密码，发送邮件
func PostMail(user *models.User, service string, subject string, url string, body string) (err error) {
	if config.GlobalConf.DisableMail {
		return nil
	}
	code := uuid.New().String()
	ctx := context.Background()
	redis := database.GetRedis()
	redis.Set(ctx, fmt.Sprintf("%s:%d", service, user.ID), code, time.Hour*2)
	err = sendMail(user, code, subject, url, body)
	if err != nil {
		return err
	}
	return nil
}

// CheckCode 检查邮件验证码是否正确
func CheckCode(user *models.User, code string, service string) (ok bool, err error) {
	if config.GlobalConf.DisableMail {
		return true, nil
	}
	ctx := context.Background()
	rds := database.GetRedis()
	key := fmt.Sprintf("%s:%d", service, user.ID)
	result := rds.Get(ctx, key)
	if err := result.Err(); err != nil {
		return false, errors.Wrap(err, errors.MailCodeInvalid)
	}
	if result.Val() != code {
		return false, nil
	}
	rds.Del(ctx, key)
	return true, nil
}
