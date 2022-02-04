package controllers

import (
	"coursebench-backend/internal/controllers/teacher"
	"github.com/gofiber/fiber/v2"
)

func TeacherRoute(r fiber.Router) {
	route := r.Group("/teacher")
	route.Get("/all", teacher.All)
	route.Get("/:id", teacher.Info)
}
