package test

import (
	"coursebench-backend/internal/controllers/test/teacher"
	"github.com/gofiber/fiber/v2"
)

func TeacherRoutes(r fiber.Router) {
	route := r.Group("/teacher")
	route.Post("/add", teacher.Add)
}
