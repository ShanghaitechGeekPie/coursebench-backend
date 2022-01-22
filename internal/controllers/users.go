package controllers

import (
	"coursebench-backend/internal/controllers/users"
	"github.com/gofiber/fiber/v2"
)

func UserRoutes(r fiber.Router) {
	route := r.Group("/user")
	route.Post("/register", users.Register)
	route.Post("/login", users.Login)
	route.Post("/logout", users.Logout)
}
