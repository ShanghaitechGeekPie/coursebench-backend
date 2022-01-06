package fiber

import (
	"github.com/gofiber/fiber/v2"
	"os"
	"time"
)

func getFiberConfig() fiber.Config {
	readTimeout, err := time.ParseDuration(os.Getenv("FIBER_READ_TIMEOUT"))
	if err != nil {
		panic(err)
	}
	writeTimeout, err := time.ParseDuration(os.Getenv("FIBER_WRITE_TIMEOUT"))
	if err != nil {
		panic(err)
	}
	idleTimeout, err := time.ParseDuration(os.Getenv("FIBER_IDLE_TIMEOUT"))
	if err != nil {
		panic(err)
	}
	return fiber.Config{
		ErrorHandler: errorHandler,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	}
}
