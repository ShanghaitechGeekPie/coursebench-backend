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
	"coursebench-backend/internal/config"
	"coursebench-backend/pkg/errors"
	"coursebench-backend/pkg/models"
	"fmt"
	"github.com/spf13/viper"
	"gopkg.in/gomail.v2"
	"strings"
)

type SMTPConfig struct {
	Address  string `mapstructure:"address"` // 发件邮箱地址
	Name     string `mapstructure:"name"`    // 发件人名称
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	//ServiceName string `mapstructure:"service_name"`
}

var smtpConfig SMTPConfig

// sendMail 发送邮件
// code: 验证码
// subject: 邮件主题
// url: 激活地址链接
// body: 邮件正文(可使用 {activeURL} 作为激活链接的占位符)
func sendMail(user *models.User, code string, subject string, url string, body string) (err error) {
	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress(smtpConfig.Address, smtpConfig.Name))
	m.SetHeader("To", user.Email)
	m.SetHeader("Subject", subject)
	activeUrl := fmt.Sprintf("%s/%s?id=%d&code=%s", config.GlobalConf.ServerURL, url, user.ID, code)
	body = strings.Replace(body, "{activeURL}", activeUrl, -1)
	m.SetBody("text/html", body)
	d := gomail.NewDialer(smtpConfig.Host, smtpConfig.Port, smtpConfig.Username, smtpConfig.Password)
	if err = d.DialAndSend(m); err != nil {
		return errors.Wrap(err, errors.SMTPError)
	}
	return nil
}

func InitSMTP() {
	cfg := viper.Sub("smtp")
	if cfg == nil {
		panic("SMTP config not found")
	}

	err := cfg.Unmarshal(&smtpConfig)
	if err != nil {
		panic(err)
	}

}
