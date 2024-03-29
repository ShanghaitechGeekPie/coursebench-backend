package controllers

import (
	"coursebench-backend/internal/controllers/test"
	"github.com/gofiber/fiber/v2"
)

func TestRoutes(r fiber.Router) {
	route := r.Group("/test")
	route.Get("/my_profile", test.MyProfile)
	test.CourseRoutes(route)
	test.TeacherRoutes(route)
	test.CourseGroupRoutes(route)
	test.CaptchaRoutes(route)
}
