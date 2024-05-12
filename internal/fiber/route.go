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

package fiber

import (
	"coursebench-backend/internal/config"
	"coursebench-backend/internal/controllers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func Routes(app *fiber.App) {
	route := app.Group("/v1")
	route.Use(cors.New())
	controllers.UserRoutes(route)
	if config.GlobalConf.InDevelopment {
		controllers.TestRoutes(route)
	}
	controllers.CourseRoutes(route)
	controllers.CommentRoutes(route)
	controllers.TeacherRoute(route)
	controllers.RewardRoutes(route)
}
