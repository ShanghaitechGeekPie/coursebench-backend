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

	"github.com/lib/pq"
	"gorm.io/gorm"
)

const ScoreLength = 4

type Comment struct {
	gorm.Model
	UserID              uint `gorm:"index"`
	User                User
	CourseGroup         CourseGroup
	CourseGroupID       uint `gorm:"index"`
	CourseID            uint `gorm:"index"`
	Semester            int
	Scores              pq.Int64Array `gorm:"type:bigint[]"`
	Title               string
	Content             string
	StudentScoreRanking int
	IsAnonymous         bool
	CreateTime          int
	UpdateTime          int
	Like                int
	Dislike             int
	IsFold              bool `gorm:"default:false"`
	IsCovered           bool `gorm:"default:false"`
	CoverTitle          string
	CoverContent        string
	CoverReason         string
	Reward              int
}

type CommentLike struct {
	gorm.Model
	UserID    uint `gorm:"index"`
	CommentID uint `gorm:"index"`
	IsLike    bool `gorm:"index"`
}

func init() {
	modelRegister.Register(&Comment{})
	modelRegister.Register(&CommentLike{})
}
