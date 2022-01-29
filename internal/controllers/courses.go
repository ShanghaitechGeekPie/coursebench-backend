package controllers

import (
	"coursebench-backend/internal/controllers/courses"
	"github.com/gofiber/fiber/v2"
)

func CourseRoutes(r fiber.Router) {
	route := r.Group("/course")
	route.Get("/all", courses.All)
}
