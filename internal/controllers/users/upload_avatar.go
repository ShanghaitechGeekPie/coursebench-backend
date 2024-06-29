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
	"coursebench-backend/pkg/database"
	"coursebench-backend/pkg/errors"
	"coursebench-backend/pkg/models"
	"coursebench-backend/pkg/queries"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	syslog "log"
)

func UploadAvatar(c *fiber.Ctx) (err error) {
	id, err := session.GetUserID(c)
	if err != nil {
		return err
	}
	file, err := c.FormFile("avatar")
	if err != nil {
		return errors.Wrap(err, errors.InvalidArgument)
	}
	if config.GlobalConf.AvatarSizeLimit < file.Size {
		return errors.Wrap(err, errors.FileTooLarge)
	}
	r, err := file.Open()
	if err != nil {
		return errors.Wrap(err, errors.InternalServerError)
	}
	user, err := queries.GetUserByID(nil, id)
	if err != nil {
		return err
	}
	oldAvatar := user.Avatar
	nameUUID, err := uuid.NewRandom()
	if err != nil {
		return errors.Wrap(err, errors.InternalServerError)
	}
	name := "avatar/" + nameUUID.String()
	err = database.UploadFile(c.Context(), name, r, file.Size)
	if err != nil {
		return err
	}
	user.Avatar = nameUUID.String()
	db := database.GetDB()
	if err = db.Save(&user).Error; err != nil {
		return errors.Wrap(err, errors.DatabaseError)
	}
	if oldAvatar != "" {
		err = database.DeleteFile(c.Context(), "avatar/"+oldAvatar)
		if err != nil {
			syslog.Println(err)
		}
	}
	profile, err := queries.GetProfile(nil, id, id)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(models.OKResponse{Data: map[string]string{"avatar": profile.Avatar}, Error: false})
}
