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
	"coursebench-backend/internal/middlewares/session"
	"coursebench-backend/pkg/errors"
	"coursebench-backend/pkg/models"
	"coursebench-backend/pkg/queries"
	"github.com/gofiber/fiber/v2"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Captcha  string `json:"captcha"`
}

func Login(c *fiber.Ctx) (err error) {
	c.Accepts("application/json")
	var request LoginRequest
	if err = c.BodyParser(&request); err != nil {
		return errors.Wrap(err, errors.InvalidArgument)
	}
	if !config.GlobalConf.DisableCaptcha {
		if err = queries.VerifyCaptcha(c, request.Captcha); err != nil {
			return err
		}
	}

	user, err := queries.Login(nil, request.Email, request.Password)
	if err != nil {
		return
	}
	sess, err := session.GetStore().Get(c)
	if err != nil {
		return errors.Wrap(err, errors.InternalServerError)
	}
	sess.Set("user_id", user.ID)
	// Save session
	if err := sess.Save(); err != nil {
		return errors.Wrap(err, errors.InternalServerError)
	}
	/*err = session.SetSession(user.ID, sess)
	if err != nil {
		return
	}*/

	response, err := queries.GetProfile(nil, user.ID, user.ID)
	if err != nil {
		return
	}

	return c.Status(fiber.StatusOK).JSON(models.OKResponse{
		Data:  response,
		Error: false,
	})
}
