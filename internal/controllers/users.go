// Copyright (C) 2021-2024 ShanghaiTech GeekPie
// This file is part of CourseBench Backend.
//
// CourseBench Backend is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// CourseBench Backend is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with CourseBench Backend.  If not, see <http://www.gnu.org/licenses/>.

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
	route.Get("/my_id", users.MyID)
}
