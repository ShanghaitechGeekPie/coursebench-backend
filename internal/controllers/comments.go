package controllers

import (
	"coursebench-backend/internal/controllers/comments"
	"github.com/gofiber/fiber/v2"
)

func CommentRoutes(r fiber.Router) {
	route := r.Group("/comment")
	route.Post("/post", comments.Post)
	route.Post("/update", comments.Update)
	route.Post("/delete", comments.Delete)
}
