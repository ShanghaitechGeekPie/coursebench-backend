package mail

import (
	"coursebench-backend/internal/config"
	"coursebench-backend/pkg/errors"
	"coursebench-backend/pkg/models"
	"fmt"
	"github.com/spf13/viper"
	"gopkg.in/gomail.v2"
)

type SMTPConfig struct {
	Address     string `mapstructure:"address"` // 发件邮箱地址
	Name        string `mapstructure:"name"`    // 发件人名称
	Host        string `mapstructure:"host"`
	Port        int    `mapstructure:"port"`
	Username    string `mapstructure:"username"`
	Password    string `mapstructure:"password"`
	ServiceName string `mapstructure:"service_name"`
}

var smtpConfig SMTPConfig

func sendMail(user *models.User, code string) (err error) {
	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress(smtpConfig.Address, smtpConfig.Name))
	m.SetHeader("To", user.Email)
	m.SetHeader("Subject", "用户注册验证")
	activeUrl := fmt.Sprintf("%s/active?id=%d&code=%s", config.GlobalConf.ServerURL, user.ID, code)
	m.SetBody("text/html", fmt.Sprintf(`<html><body><h1>欢迎注册%s</h1> <p>请点击该链接完成注册:</p><a href="%s">注册链接 </a> <br> <p>如果链接无法点击，请手动复制该链接并粘贴至浏览器：%s </p><br><br> <p>如果您没有注册过我们的服务，请无视该邮件</p> </body></html>`, smtpConfig.ServiceName, activeUrl, activeUrl))
	d := gomail.NewDialer(smtpConfig.Host, smtpConfig.Port, smtpConfig.Username, smtpConfig.Password)
	if err = d.DialAndSend(m); err != nil {
		return errors.Wrap(err, errors.SMTPError)
	}
	return nil
}

func InitSMTP() {
	config := viper.Sub("smtp")
	if config == nil {
		panic("SMTP config not found")
	}

	err := config.Unmarshal(&smtpConfig)
	if err != nil {
		panic(err)
	}

}
