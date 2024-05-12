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
	"coursebench-backend/internal/config"
	"coursebench-backend/pkg/errors"
	"coursebench-backend/pkg/models"
	"coursebench-backend/pkg/queries"

	"github.com/gofiber/fiber/v2"
)

type RegisterRequest struct {
	Email          string           `json:"email"`
	Password       string           `json:"password"`
	Year           int              `json:"year"`
	Grade          models.GradeType `json:"grade"`
	Captcha        string           `json:"captcha"`
	Nickname       string           `json:"nickname"`
	InvitationCode string           `json:"invitation_code"`
}

func Register(c *fiber.Ctx) (err error) {
	c.Accepts("application/json")
	var userReq RegisterRequest
	if err = c.BodyParser(&userReq); err != nil {
		return errors.Wrap(err, errors.InvalidArgument)
	}
	if !config.GlobalConf.DisableCaptcha {
		if err = queries.VerifyCaptcha(c, userReq.Captcha); err != nil {
			return
		}
	}

	user := models.User{
		Email:       userReq.Email,
		Password:    userReq.Password,
		Year:        userReq.Year,
		Grade:       userReq.Grade,
		NickName:    userReq.Nickname,
		Avatar:      "",
		IsAnonymous: false,
	}
	if err = queries.Register(nil, &user, userReq.InvitationCode); err != nil {
		return
	}
	if config.GlobalConf.DisableMail {
		err = queries.RegisterActive(nil, user.ID, "")
		if err != nil {
			return err
		}
	}

	return c.Status(fiber.StatusOK).JSON(models.OKResponse{
		Data:  map[string]interface{}{"UserID": user.ID},
		Error: false,
	})
}
