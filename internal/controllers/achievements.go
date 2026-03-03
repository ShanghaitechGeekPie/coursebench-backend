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
	"coursebench-backend/internal/controllers/achievement"

	"github.com/gofiber/fiber/v2"
)

func AchievementRoutes(r fiber.Router) {
	route := r.Group("/achievement")

	// 用户相关接口
	route.Get("/list", achievement.List)                 // 获取当前用户的成就列表
	route.Get("/stats", achievement.Stats)               // 获取当前用户的成就统计
	route.Get("/user/:id", achievement.UserAchievements) // 获取指定用户的公开成就

	// 管理员接口
	route.Post("/admin/initialize", achievement.AdminInitialize)   // 初始化默认成就
	route.Post("/admin/recalculate", achievement.AdminRecalculate) // 重新计算用户成就
}
