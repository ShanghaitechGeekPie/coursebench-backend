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
	Course       []Course
}

func init() {
	modelRegister.Register(&Teacher{})
}
