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
	if config.GlobalConf.DisableCaptchaAndMail {
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
	if config.GlobalConf.DisableCaptchaAndMail {
		return true, nil
	}
	ctx := context.Background()
	rds := database.GetRedis()
	result := rds.Get(ctx, fmt.Sprintf("%s:%d", service, user.ID))
	if err := result.Err(); err != nil {
		return false, errors.Wrap(err, errors.MailCodeInvalid)
	}
	if result.Val() != code {
		return false, nil
	}
	return true, nil
}
