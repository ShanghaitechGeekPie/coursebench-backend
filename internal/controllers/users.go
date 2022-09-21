package controllers

import (
	"coursebench-backend/internal/controllers/users"
	"github.com/gofiber/fiber/v2"
)

func UserRoutes(r fiber.Router) {
	route := r.Group("/user")
	route.Post("/register", users.Register)
	route.Post("/register_active", users.RegisterActive)
	route.Post("/login", users.Login)
	route.Post("/logout", users.Logout)
	route.Get("/profile/:id", users.Profile)
	route.Post("/update_profile", users.UpdateProfile)
	route.Post("/update_password", users.UpdatePassword)
	route.Post("/get_captcha", users.GetCaptcha)
	route.Post("/upload_avatar", users.UploadAvatar)
	route.Post("/reset_password", users.ResetPassword)
	route.Post("/reset_password_active", users.ResetPasswordActive)
}
