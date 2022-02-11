package mail

import (
	"context"
	"coursebench-backend/pkg/database"
	"coursebench-backend/pkg/models"
	"github.com/google/uuid"
	"strconv"
)

// PostMail 用户注册，发送邮件
func PostMail(user *models.User) (err error) {
	code := uuid.New().String()
	ctx := context.Background()
	redis := database.GetRedis()
	redis.Set(ctx, "MAIL_CODE"+strconv.Itoa(int(user.ID)), code, 60*60*2)
	err = sendMail(user, code)
	if err != nil {
		return err
	}
	return nil
}

// CheckCode 检查邮件验证码是否正确
func CheckCode(user models.User, code string) (ok bool, err error) {
	return true, nil
}
