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

package comments

import (
	"coursebench-backend/internal/middlewares/session"
	"coursebench-backend/pkg/database"
	"coursebench-backend/pkg/errors"
	"coursebench-backend/pkg/models"
	"coursebench-backend/pkg/queries"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type DeleteRequest struct {
	ID uint `json:"id"`
}

func Delete(c *fiber.Ctx) (err error) {
	c.Accepts("application/json")
	var request DeleteRequest
	if err = c.BodyParser(&request); err != nil {
		return errors.Wrap(err, errors.InvalidArgument)
	}

	uid, err := session.GetUserID(c)
	if err != nil {
		return err
	}

	user, err := queries.GetUserByID(nil, uid)
	if err != nil {
		return err
	}

	db := database.GetDB()
	comment := &models.Comment{}
	err = db.Preload("CourseGroup").Preload("CourseGroup.Course").Where("id = ?", request.ID).Take(comment).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New(errors.CommentNotExists)
		} else {
			return errors.Wrap(err, errors.DatabaseError)
		}
	}

	if !user.IsAdmin && uid != comment.UserID { // admin can delete any comment.
		return errors.New(errors.PermissionDenied)
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		comment.CourseGroup.CommentCount--
		for i := 0; i < models.ScoreLength; i++ {
			comment.CourseGroup.Scores[i] -= comment.Scores[i]
		}
		comment.CourseGroup.Course.CommentCount--
		for i := 0; i < models.ScoreLength; i++ {
			comment.CourseGroup.Course.Scores[i] -= comment.Scores[i]
		}
		err = tx.Select("Scores", "CommentCount").Updates(comment.CourseGroup).Error
		if err != nil {
			return errors.Wrap(err, errors.DatabaseError)
		}

		err = tx.Select("Scores", "CommentCount").Updates(comment.CourseGroup.Course).Error
		if err != nil {
			return errors.Wrap(err, errors.DatabaseError)
		}

		err = tx.Delete(comment).Error
		if err != nil {
			return errors.Wrap(err, errors.DatabaseError)
		}

		return nil
	})

	return c.Status(fiber.StatusOK).JSON(models.OKResponse{
		Data:  nil,
		Error: false,
	})
}
