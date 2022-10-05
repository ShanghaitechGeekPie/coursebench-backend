package config

import (
	"github.com/spf13/viper"
	"time"
)

type FiberConfigType struct {
	ReadTimeout    time.Duration `mapstructure:"read_timeout"`
	WriteTimeout   time.Duration `mapstructure:"write_timeout"`
	IdleTimeout    time.Duration `mapstructure:"idle_timeout"`
	InDevelopment  bool          `mapstructure:"in_development"`
	Listen         string        `mapstructure:"listen"`
	UseXForwardFor bool          `mapstructure:"use_x_forward_for"`
}

var FiberConfig FiberConfigType

func SetupFiberConfig() {
	cfg := viper.Sub("fiber")
	if cfg == nil {
		cfg = viper.New()
	}
	cfg.SetDefault("read_timeout", "10s")
	cfg.SetDefault("write_timeout", "10s")
	cfg.SetDefault("idle_timeout", "1m")
	cfg.SetDefault("in_development", GlobalConf.InDevelopment)
	cfg.SetDefault("listen", "0.0.0.0:10001")
	jsonErr := cfg.Unmarshal(&FiberConfig)
	if jsonErr != nil {
		panic(jsonErr)
	}
}
