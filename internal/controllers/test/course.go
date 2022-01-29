package test

import (
	"coursebench-backend/internal/controllers/test/course"
	"github.com/gofiber/fiber/v2"
)

func CourseRoutes(r fiber.Router) {
	route := r.Group("/course")
	route.Post("/add", course.Add)
}
