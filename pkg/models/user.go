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
	Email            string `gorm:"index"`
	Password         string
	NickName         string
	RealName         string
	Year             int
	Grade            GradeType
	IsActive         bool
	Avatar           string
	IsAnonymous      bool
	IsAdmin          bool `gorm:"default:false"`
	IsCommunityAdmin bool `gorm:"default:false"`
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
}
