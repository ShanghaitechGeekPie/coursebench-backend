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
)

type AllResponse struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	Institute    string `json:"institute"`
	Photo        string `json:"photo"`
	Job          string `json:"job"`
	Introduction string `json:"introduction"`
}

func All(c *fiber.Ctx) error {
	db := database.GetDB()
	var teachers []models.Teacher
	result := db.Find(&teachers)
	if result.Error != nil {
		return errors.Wrap(result.Error, errors.DatabaseError)
	}

	var response []AllResponse
	for _, v := range teachers {
		response = append(response, AllResponse{
			ID:           v.ID,
			Name:         v.Name,
			Institute:    v.Institute,
			Job:          v.Job,
			Introduction: v.Introduction,
			Photo:        v.Photo,
		})
	}

	return c.Status(fiber.StatusOK).JSON(models.OKResponse{
		Data:  response,
		Error: false,
	})
}
