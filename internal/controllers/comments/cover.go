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
	"bytes"
	"coursebench-backend/internal/config"
	"coursebench-backend/internal/middlewares/session"
	"coursebench-backend/pkg/database"
	"coursebench-backend/pkg/errors"
	"coursebench-backend/pkg/models"
	"coursebench-backend/pkg/queries"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"io"
	"net/http"
)

type CoverRequest struct {
	ID     int  `json:"id"`
	Status bool `json:"status"`
}

type GPTWorkerResponse struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Reason  string `json:"reason"`
}

func Cover(c *fiber.Ctx) (err error) {
	c.Accepts("application/json")
	var request FoldRequest
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
	if !user.IsAdmin && !user.IsCommunityAdmin {
		return errors.New(errors.PermissionDenied)
	}

	db := database.GetDB()
	err = db.Transaction(func(tx *gorm.DB) error {
		comment := &models.Comment{}
		err = tx.Where("id = ?", request.ID).Take(comment).Error

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New(errors.CommentNotExists)
			} else {
				return errors.Wrap(err, errors.DatabaseError)
			}
		}
		if request.Status { // cover
			requestBody, err := json.Marshal(map[string]string{"title": comment.Title, "content": comment.Content})
			if err != nil {
				return errors.Wrap(err, errors.InternalServerError)
			}
			response, err := http.Post(config.GlobalConf.GPTWorkerURL, "application/json", bytes.NewBuffer(requestBody))
			if err != nil {
				return errors.New(errors.GPTWorkerError)
			}
			defer response.Body.Close()
			var gptWorkerResponse GPTWorkerResponse
			body, err := io.ReadAll(response.Body)
			if err != nil {
				return errors.New(errors.GPTWorkerError)
			}
			err = json.Unmarshal(body, &gptWorkerResponse)
			if err != nil {
				return errors.New(errors.GPTWorkerError)
			}
			comment.CoverTitle = gptWorkerResponse.Title
			comment.CoverContent = gptWorkerResponse.Content
			comment.CoverReason = gptWorkerResponse.Reason
			comment.IsCovered = true
		} else { // uncover
			comment.CoverTitle = ""
			comment.CoverContent = ""
			comment.CoverReason = ""
			comment.IsCovered = false
		}

		err = tx.Select("IsCovered", "CoverTitle", "CoverContent", "CoverReason").Save(comment).Error
		if err != nil {
			return errors.Wrap(err, errors.DatabaseError)
		}
		return nil
	})
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(models.OKResponse{
		Data:  map[string]interface{}{},
		Error: false})
}
