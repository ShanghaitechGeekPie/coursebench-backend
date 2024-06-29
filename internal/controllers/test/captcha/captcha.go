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

package captcha

import (
	"coursebench-backend/pkg/errors"
	"coursebench-backend/pkg/models"
	"coursebench-backend/pkg/queries"
	"github.com/gofiber/fiber/v2"
)

type VerifyRequest struct {
	Code string `json:"code"`
}

func Get(c *fiber.Ctx) (err error) {
	c.Accepts("application/json")
	s, err := queries.GenerateCaptcha(c)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(models.OKResponse{
		Data:  s,
		Error: false,
	})
}

func Verify(c *fiber.Ctx) (err error) {
	c.Accepts("application/json")
	var request VerifyRequest
	if err = c.BodyParser(&request); err != nil {
		return errors.Wrap(err, errors.InvalidArgument)
	}
	err = queries.VerifyCaptcha(c, request.Code)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(models.OKResponse{
		Data:  nil,
		Error: false,
	})
}
