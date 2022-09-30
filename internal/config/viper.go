package config

import (
	"fmt"
	"github.com/spf13/viper"
	syslog "log"
)

type GlobalConfig struct {
	InDevelopment   bool   `mapstructure:"in_development"`
	ServerURL       string `mapstructure:"server_url"`
	DisableCaptcha  bool   `mapstructure:"disable_captcha"`
	DisableMail     bool   `mapstructure:"disable_mail"`
	AvatarSizeLimit int64  `mapstructure:"avatar_size_limit"`
	MailSuffix      string `mapstructure:"mail_suffix"`
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

	config = viper.Sub("text")
	err = config.Unmarshal(&Text)
	if err != nil {
		syslog.Fatalf("Fatal error text config file: %v \n", err)
	}
}
