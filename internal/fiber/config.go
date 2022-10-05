package fiber

import (
	"coursebench-backend/internal/config"
	"github.com/gofiber/fiber/v2"
)

func getFiberConfig() fiber.Config {
	return fiber.Config{
		ErrorHandler: errorHandler,
		ReadTimeout:  config.FiberConfig.ReadTimeout,
		WriteTimeout: config.FiberConfig.WriteTimeout,
		IdleTimeout:  config.FiberConfig.IdleTimeout,
	}
}
