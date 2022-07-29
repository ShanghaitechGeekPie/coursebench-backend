package fiber

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
)

func New() (app *fiber.App) {
	app = fiber.New(getFiberConfig())
	app.Get("/metrics_monitor", monitor.New())
	return app
}
