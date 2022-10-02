package fiber

import (
	"coursebench-backend/pkg/models"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"time"
)

func New() (app *fiber.App) {
	app = fiber.New(getFiberConfig())
	app.Get("/metrics_monitor", monitor.New())
	app.Use(limiter.New(limiter.Config{
		Max: 60,
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(models.ErrorResponse{
				Error:     true,
				Errno:     "TooManyRequests",
				Message:   "Too many requests, please try again a minute later",
				Timestamp: time.Now(),
			})
		},
	}))
	app.Use(LogMiddleware)
	return app
}
