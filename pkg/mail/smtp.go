package mail

import (
	"coursebench-backend/pkg/errors"
	"coursebench-backend/pkg/models"
	"fmt"
	"gopkg.in/gomail.v2"
)

func sendMail(user *models.User, code string) (err error) {
	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress("geekpie@geekpie.club", "GeekPie Services"))
	m.SetHeader("To", user.Email)
	m.SetHeader("Subject", "用户注册验证")
	m.SetBody("text/html", fmt.Sprintf(`<html><body><h1>欢迎注册上海科技大学评教系统</h1> <p>请点击该链接完成注册:</p><a href=\"https://%s/activity/%s\n\">注册链接 </a> <br><br><br> <p>如果您没有注册过我们的服务，请无视该邮件</p> </body></html>`, "www.geekpie.club", code))
	d := gomail.NewDialer("smtp.qq.com", 465, "114514", "1919810")
	if err = d.DialAndSend(m); err != nil {
		return errors.Wrap(err, errors.SMTPError)
	}
	return nil
}
