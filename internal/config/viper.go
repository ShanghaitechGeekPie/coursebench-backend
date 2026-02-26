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

package config

import (
	"fmt"
	syslog "log"
	"os"

	"github.com/spf13/viper"
)

type GlobalConfig struct {
	InDevelopment       bool   `mapstructure:"in_development"`
	ServerURL           string `mapstructure:"server_url"`
	DisableCaptcha      bool   `mapstructure:"disable_captcha"`
	DisableMail         bool   `mapstructure:"disable_mail"`
	AvatarSizeLimit     int64  `mapstructure:"avatar_size_limit"`
	MailSuffix          string `mapstructure:"mail_suffix"`
	GPTWorkerURL        string `mapstructure:"gpt_worker_url"`
	CasdoorEndpoint     string `mapstructure:"casdoor_endpoint"`
	CasdoorClientID     string `mapstructure:"casdoor_client_id"`
	CasdoorClientSecret string `mapstructure:"casdoor_client_secret"`
	CasdoorOrganization string `mapstructure:"casdoor_organization"`
	CasdoorApplication  string `mapstructure:"casdoor_application"`
	CasdoorRedirectURI  string `mapstructure:"casdoor_redirect_uri"`
	CasdoorFrontendURL  string `mapstructure:"casdoor_frontend_url"`
}
type TextConfig struct {
	ServiceName   string `mapstructure:"service_name"`
	ServiceNameEN string `mapstructure:"service_name_en"`
}

var GlobalConf GlobalConfig
var Text TextConfig

func SetupViper() {
	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/coursebench/")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}

	config := viper.Sub("global")
	config.SetDefault("in_development", false)
	config.SetDefault("disable_captcha_and_mail", false)
	err = config.Unmarshal(&GlobalConf)
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}

	if v := os.Getenv("CASDOOR_ENDPOINT"); v != "" {
		GlobalConf.CasdoorEndpoint = v
	}
	if v := os.Getenv("CASDOOR_CLIENT_ID"); v != "" {
		GlobalConf.CasdoorClientID = v
	}
	if v := os.Getenv("CASDOOR_CLIENT_SECRET"); v != "" {
		GlobalConf.CasdoorClientSecret = v
	}
	if v := os.Getenv("CASDOOR_ORGANIZATION"); v != "" {
		GlobalConf.CasdoorOrganization = v
	}
	if v := os.Getenv("CASDOOR_APPLICATION"); v != "" {
		GlobalConf.CasdoorApplication = v
	}
	if v := os.Getenv("CASDOOR_REDIRECT_URI"); v != "" {
		GlobalConf.CasdoorRedirectURI = v
	}
	if v := os.Getenv("CASDOOR_FRONTEND_URL"); v != "" {
		GlobalConf.CasdoorFrontendURL = v
	}

	config = viper.Sub("text")
	err = config.Unmarshal(&Text)
	if err != nil {
		syslog.Fatalf("Fatal error text config file: %v \n", err)
	}

	SetupFiberConfig()
}
