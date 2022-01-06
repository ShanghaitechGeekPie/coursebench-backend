package models

import (
	"coursebench-backend/pkg/modelRegister"
	"gorm.io/gorm"
)

type GradeType int

const (
	Undergraduate GradeType = 1
	Postgraduate  GradeType = 2
	PhDStudent    GradeType = 3
)

type User struct {
	gorm.Model
	Email    string
	Password string
	NickName string
	RealName string
	Year     int
	Grade    GradeType
}

func init() {
	modelRegister.Register(&User{})
}
