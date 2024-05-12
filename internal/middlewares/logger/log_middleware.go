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

package logger

import (
	"coursebench-backend/pkg/log"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"strings"
)

// LogMiddleware record the details of each request.
func LogMiddleware(c *fiber.Ctx) error {
	s := "New Request: " + c.String()
	s += " ; X-Forwarded-For: " + strings.Join(c.IPs(), ",")
	if c.Path() == "/v1/user/login" || c.Path() == "/v1/user/register" || c.Path() == "/v1/user/update_password" || c.Path() == "/v1/user/reset_password_active" { // Don't log the password. parsing json is too much trouble here, so I just don't record the whole body.
		s += " ; Body: Sensitive, skip... "
	} else if len(c.Body()) < 500 {
		s += " ; Body: " + string(c.Body())
	} else {
		s += fmt.Sprintf(" ; Body: Too long, %v bytes, skip...", len(c.Body()))
	}
	log.Println(s)
	return c.Next()
}
