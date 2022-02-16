package fiber

import (
	"coursebench-backend/internal/config"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"time"
)

type Config struct {
	ReadTimeout   time.Duration `mapstructure:"read_timeout"`
	WriteTimeout  time.Duration `mapstructure:"write_timeout"`
	IdleTimeout   time.Duration `mapstructure:"idle_timeout"`
	InDevelopment bool          `mapstructure:"in_development"`
	Listen        string        `mapstructure:"listen"`
}

var FiberConfig Config

func getFiberConfig() fiber.Config {
	cfg := viper.Sub("fiber")
	if cfg == nil {
		cfg = viper.New()
	}
	cfg.SetDefault("read_timeout", "10s")
	cfg.SetDefault("write_timeout", "10s")
	cfg.SetDefault("idle_timeout", "1m")
	cfg.SetDefault("in_development", config.GlobalConf.InDevelopment)
	cfg.SetDefault("listen", "0.0.0.0:10001")
	jsonErr := cfg.Unmarshal(&FiberConfig)
	if jsonErr != nil {
		panic(jsonErr)
	}

	return fiber.Config{
		ErrorHandler: errorHandler,
		ReadTimeout:  FiberConfig.ReadTimeout,
		WriteTimeout: FiberConfig.WriteTimeout,
		IdleTimeout:  FiberConfig.IdleTimeout,
	}
}
