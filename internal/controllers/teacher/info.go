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

package teacher

import (
	"coursebench-backend/pkg/database"
	"coursebench-backend/pkg/errors"
	"coursebench-backend/pkg/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"strconv"
)

type InfoResponse struct {
	ID           uint                       `json:"id"`
	Name         string                     `json:"name"`
	Institute    string                     `json:"institute"`
	Job          string                     `json:"job"`
	Introduction string                     `json:"introduction"`
	Photo        string                     `json:"photo"`
	Courses      []models.CourseAllResponse `json:"courses"`
}

func Info(c *fiber.Ctx) (err error) {
	id_s := c.Params("id", "GG")
	id_i, err := strconv.Atoi(id_s)
	if err != nil {
		return errors.Wrap(err, errors.InvalidArgument)
	}
	id := uint(id_i)
	db := database.GetDB()
	teachers := &models.Teacher{}
	err = db.Preload("Courses").Where("id = ?", id).Take(&teachers).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.Wrap(err, errors.DatabaseError)
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.Wrap(err, errors.TeacherNotExists)
	}
	response := InfoResponse{
		ID:           teachers.ID,
		Name:         teachers.Name,
		Institute:    teachers.Institute,
		Job:          teachers.Job,
		Introduction: teachers.Introduction,
		Photo:        teachers.Photo,
		Courses:      make([]models.CourseAllResponse, 0),
	}
	for _, v := range teachers.Courses {
		score := make([]float64, models.ScoreLength)
		if v.CommentCount != 0 {
			for j := 0; j < models.ScoreLength; j++ {
				score[j] = float64(v.Scores[j]) / float64(v.CommentCount)
			}
		}
		response.Courses = append(response.Courses, models.CourseAllResponse{
			ID:         int(v.ID),
			Name:       v.Name,
			Institute:  v.Institute,
			Code:       v.Code,
			Score:      score,
			Credit:     v.Credit,
			CommentNum: v.CommentCount,
		})
	}

	return c.Status(fiber.StatusOK).JSON(models.OKResponse{
		Data:  response,
		Error: false,
	})
}
