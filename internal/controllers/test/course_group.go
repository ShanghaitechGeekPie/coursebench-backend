package test

import (
	"coursebench-backend/internal/controllers/test/course_group"
	"github.com/gofiber/fiber/v2"
)

func CourseGroupRoutes(r fiber.Router) {
	route := r.Group("/course_group")
	route.Post("/add", course_group.Add)
}
