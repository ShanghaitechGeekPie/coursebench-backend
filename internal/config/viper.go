package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type GlobalConfig struct {
	InDevelopment bool `mapstructure:"in_development"`
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
	err = config.Unmarshal(&GlobalConf)
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}

}
