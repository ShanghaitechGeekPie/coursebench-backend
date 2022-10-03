package fiber

import (
	"coursebench-backend/internal/middlewares/logger"
	"coursebench-backend/internal/middlewares/session"
	"coursebench-backend/pkg/models"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"time"
)

func New() (app *fiber.App) {
	session.Init()
	app = fiber.New(getFiberConfig())
	app.Get("/metrics_monitor", monitor.New())
	app.Use(limiter.New(limiter.Config{
		KeyGenerator: func(c *fiber.Ctx) string {
			if FiberConfig.UseXForwardFor {
				if len(c.IPs()) > 0 {
					return c.IPs()[0]
				} else {
					return c.IP()
				}
			} else {
				return c.IP()
			}
		},
		Max: 300,
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(models.ErrorResponse{
				Error:     true,
				Errno:     "TooManyRequests",
				Message:   "Too many requests, please try again a minute later",
				Timestamp: time.Now(),
			})
		},
	}))
	app.Use(logger.LogMiddleware)
	return app
}
