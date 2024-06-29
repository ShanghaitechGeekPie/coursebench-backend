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

package test

import (
	"coursebench-backend/internal/middlewares/session"
	"coursebench-backend/pkg/models"
	"coursebench-backend/pkg/queries"
	"fmt"
	"github.com/gofiber/fiber/v2"
)

func MyProfile(c *fiber.Ctx) (err error) {
	fmt.Println(c.Cookies("session_id", "GG"))
	id, err := session.GetUserID(c)
	if err != nil {
		return err
	}
	user, err := queries.GetUserByID(nil, id)
	return c.Status(fiber.StatusOK).JSON(models.OKResponse{
		Data:  map[string]interface{}{"UserID": id, "User": user},
		Error: false,
	})
}
