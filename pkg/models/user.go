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

type GradeType int

const (
	UnknownGrade  GradeType = 0
	Undergraduate GradeType = 1
	Postgraduate  GradeType = 2
	PhDStudent    GradeType = 3
)

type User struct {
	gorm.Model
	Email             string `gorm:"index"`
	Password          string
	NickName          string
	RealName          string
	Year              int
	Grade             GradeType
	IsActive          bool
	Avatar            string
	IsAnonymous       bool
	IsAdmin           bool `gorm:"default:false"`
	IsCommunityAdmin  bool `gorm:"default:false"`
	InvitationCode    string
	InvitedByUserID   uint
	Reward            int
	HasPostedComments bool `gorm:"default:false"`
}

func init() {
	modelRegister.Register(&User{})
}

type ProfileResponse struct {
	ID               uint      `json:"id"`
	Email            string    `json:"email"`
	Year             int       `json:"year"`
	Grade            GradeType `json:"grade"`
	NickName         string    `json:"nickname"`
	RealName         string    `json:"realname"`
	Avatar           string    `json:"avatar"`
	IsAnonymous      bool      `json:"is_anonymous"`
	IsAdmin          bool      `json:"is_admin"`
	IsCommunityAdmin bool      `json:"is_community_admin"`
	InvitationCode   string    `json:"invitation_code"`
	Reward           int       `json:"reward"`
}
