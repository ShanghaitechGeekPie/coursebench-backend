package test

import (
	"coursebench-backend/internal/controllers/test/captcha"
	"github.com/gofiber/fiber/v2"
)

func CaptchaRoutes(r fiber.Router) {
	route := r.Group("/captcha")
	route.Post("/get", captcha.Get)
	route.Post("/verify", captcha.Verify)
}
