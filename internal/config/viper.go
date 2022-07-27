package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type GlobalConfig struct {
	InDevelopment         bool   `mapstructure:"in_development"`
	ServerURL             string `mapstructure:"server_url"`
	DisableCaptchaAndMail bool   `mapstructure:"disable_captcha_and_mail"`
}

var GlobalConf GlobalConfig

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

}
