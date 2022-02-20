package models

import (
	"coursebench-backend/pkg/modelRegister"
	"gorm.io/gorm"
)

type Teacher struct {
	gorm.Model
	EamsID       int
	Name         string
	Institute    string
	Job          string
	Introduction string
	CourseGroups []*CourseGroup `gorm:"many2many:coursegroup_teachers;"`
	Courses      []*Course      `gorm:"many2many:course_teachers;"`
}

func init() {
	modelRegister.Register(&Teacher{})
}
