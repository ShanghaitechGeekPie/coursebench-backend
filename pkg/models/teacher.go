package models

import (
	"coursebench-backend/pkg/modelRegister"
	"gorm.io/gorm"
)

type Teacher struct {
	gorm.Model
	Name         string
	Job          string
	Introduction string
	Course       []*Course `gorm:"many2many:course_teachers;"`
}

func init() {
	modelRegister.Register(&Teacher{})
}
