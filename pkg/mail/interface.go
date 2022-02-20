package mail

import (
	"context"
	"coursebench-backend/pkg/database"
	"coursebench-backend/pkg/models"
	"fmt"
	"github.com/google/uuid"
)

// PostMail 用户注册，发送邮件
func PostMail(user *models.User) (err error) {
	code := uuid.New().String()
	ctx := context.Background()
	redis := database.GetRedis()
	redis.Set(ctx, fmt.Sprintf("mail_code:%d", user.ID), code, 60*60*2)
	err = sendMail(user, code)
	if err != nil {
		return err
	}
	return nil
}

// CheckCode 检查邮件验证码是否正确
func CheckCode(user models.User, code string) (ok bool, err error) {
	ctx := context.Background()
	rds := database.GetRedis()
	result := rds.Get(ctx, fmt.Sprintf("mail_code:%d", user.ID))
	if err := result.Err(); err != nil {
		return false, err
	}
	if result.Val() == code {
		return true, nil
	}
	return true, nil
}
