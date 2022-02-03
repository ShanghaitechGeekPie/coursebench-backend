package models

import (
	"coursebench-backend/pkg/modelRegister"
	"gorm.io/gorm"
)

type Teacher struct {
	gorm.Model
	Name         string
	Institute    string
	Job          string
	Introduction string
	Courses      []*CourseGroup `gorm:"many2many:course_teachers;"`
}

func init() {
	modelRegister.Register(&Teacher{})
}
