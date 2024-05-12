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
	route.Get("/recent/:id", comments.RecentCommentByPage)
	route.Post("/like", comments.Like)
	route.Post("/fold", comments.Fold)
	route.Post("/cover", comments.Cover)
}
