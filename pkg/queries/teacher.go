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

package queries

import (
	"coursebench-backend/pkg/database"
	"coursebench-backend/pkg/errors"
	"coursebench-backend/pkg/models"
	"gorm.io/gorm"
)

func GetTeacher(db *gorm.DB, id uint) (teacher *models.Teacher, err error) {
	if db == nil {
		db = database.GetDB()
	}

	teacher = &models.Teacher{}
	result := db.Where("id = ?", id).Take(teacher)
	if err := result.Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.Wrap(err, errors.DatabaseError)
	}
	if result.RowsAffected == 0 {
		return nil, errors.New(errors.TeacherNotExists)
	}

	return teacher, nil
}

func AddTeacher(db *gorm.DB, name string, job string, introduction string, EamsID int) (teacher *models.Teacher, err error) {
	if db == nil {
		db = database.GetDB()
	}

	teacher = &models.Teacher{
		Name:         name,
		Job:          job,
		Introduction: introduction,
		EamsID:       EamsID,
	}
	result := db.Create(teacher)
	if err := result.Error; err != nil {
		return nil, errors.Wrap(err, errors.DatabaseError)
	}

	return teacher, nil
}
