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
	"coursebench-backend/internal/middlewares/session"
	"coursebench-backend/pkg/errors"
	"coursebench-backend/pkg/models"
	"coursebench-backend/pkg/queries"
	"github.com/gofiber/fiber/v2"
)

type UpdateProfileRequest struct {
	Year        int              `json:"year"`
	Grade       models.GradeType `json:"grade"`
	NickName    string           `json:"nickname"`
	RealName    string           `json:"realname"`
	IsAnonymous bool             `json:"is_anonymous"`
}

func UpdateProfile(c *fiber.Ctx) (err error) {
	c.Accepts("application/json")
	var request UpdateProfileRequest
	if err = c.BodyParser(&request); err != nil {
		return errors.Wrap(err, errors.InvalidArgument)
	}
	id, err := session.GetUserID(c)
	if err != nil {
		return err
	}
	err = queries.UpdateProfile(nil, id, request.Year, request.Grade, request.NickName, request.RealName, request.IsAnonymous)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(models.OKResponse{
		Data:  nil,
		Error: false,
	})
}
