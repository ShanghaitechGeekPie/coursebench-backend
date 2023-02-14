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
	route.Get("/user/:id", comments.UserComment)
	route.Get("/course_group/:id", comments.CourseGroupComment)
	route.Get("/course/:id", comments.CourseComment)
	route.Get("/recent", comments.RecentComment)
	route.Post("/like", comments.Like)
	route.Post("/fold", comments.Fold)
}
