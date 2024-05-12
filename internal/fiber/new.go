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
	"coursebench-backend/internal/middlewares/logger"
	"coursebench-backend/internal/middlewares/session"
	"coursebench-backend/internal/utils"
	"coursebench-backend/pkg/models"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"time"
)

func New() (app *fiber.App) {
	session.Init()
	app = fiber.New(getFiberConfig())
	app.Get("/metrics_monitor", monitor.New())
	app.Use(limiter.New(limiter.Config{
		KeyGenerator: func(c *fiber.Ctx) string {
			return utils.GetIP(c)[0]
		},
		Max: 300,
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(models.ErrorResponse{
				Error:     true,
				Errno:     "TooManyRequests",
				Message:   "Too many requests, please try again a minute later",
				Timestamp: time.Now(),
			})
		},
	}))
	app.Use(logger.LogMiddleware)
	return app
}
