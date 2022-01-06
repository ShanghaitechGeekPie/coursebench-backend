package fiber

import (
	"github.com/gofiber/fiber/v2"
)

func New() (app *fiber.App) {
	app = fiber.New(getFiberConfig())
	return app
}
