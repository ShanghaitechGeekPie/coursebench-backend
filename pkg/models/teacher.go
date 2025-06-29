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

package models

import (
	"coursebench-backend/pkg/modelRegister"

	"gorm.io/gorm"
)

const TEACHER_OTHER_ID = 100000001

type Teacher struct {
	gorm.Model
	EamsID       int
	UniID        int
	Name         string
	Institute    string
	Job          string
	Introduction string
	Email        string
	Photo        string
	CourseGroups []*CourseGroup `gorm:"many2many:coursegroup_teachers;"`
	Courses      []*Course      `gorm:"many2many:course_teachers;"`
}

func init() {
	modelRegister.Register(&Teacher{})
}
