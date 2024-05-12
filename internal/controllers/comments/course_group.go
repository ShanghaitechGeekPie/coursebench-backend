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
	"coursebench-backend/internal/utils"
	"coursebench-backend/pkg/database"
	"coursebench-backend/pkg/errors"
	"coursebench-backend/pkg/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func CourseGroupComment(c *fiber.Ctx) (err error) {
	id_s := c.Params("id", "GG")
	id_i, err := strconv.Atoi(id_s)
	if err != nil {
		return errors.New(errors.InvalidArgument)
	}
	id := uint(id_i)

	uid, err := session.GetUserID(c)
	if err != nil {
		uid = 0
	}

	db := database.GetDB()
	var comments []models.Comment
	err = db.Preload("User").Preload("CourseGroup").Preload("CourseGroup.Course").Preload("CourseGroup.Teachers").
		Where("course_group_id = ?", id).Find(&comments).Error
	if err != nil {
		return errors.Wrap(err, errors.DatabaseError)
	}
	var likeResult []CommentLikeResult
	if uid != 0 {
		db.Raw("SELECT comment_likes.comment_id,comment_likes.is_like from comments, comment_likes where comments.course_group_id = ? and comment_likes.user_id = ? and comment_likes.comment_id = comments.id and comment_likes.deleted_at is NULL and comments.deleted_at is NULL",
			id, uid).Scan(&likeResult)
	}
	var response []CommentResponse
	response = GenerateResponse(comments, uid, likeResult, true, utils.GetIP(c))
	return c.Status(fiber.StatusOK).JSON(models.OKResponse{
		Data:  response,
		Error: false,
	})
}
