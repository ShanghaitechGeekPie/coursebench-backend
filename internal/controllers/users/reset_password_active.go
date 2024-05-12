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

package users

import (
	"coursebench-backend/pkg/errors"
	"coursebench-backend/pkg/models"
	"coursebench-backend/pkg/queries"
	"github.com/gofiber/fiber/v2"
)

type ResetPasswordActiveRequest struct {
	ID       uint   `json:"id"`
	Code     string `json:"code"`
	Password string `json:"password"`
}

func ResetPasswordActive(c *fiber.Ctx) (err error) {
	c.Accepts("application/json")
	var activeReq ResetPasswordActiveRequest
	if err = c.BodyParser(&activeReq); err != nil {
		return errors.Wrap(err, errors.InvalidArgument)
	}
	err = queries.ResetPasswordActive(nil, activeReq.ID, activeReq.Code, activeReq.Password)
	if err != nil {
		return
	}

	return c.Status(fiber.StatusOK).JSON(models.OKResponse{
		Data:  map[string]interface{}{"OK": true},
		Error: false,
	})
}
