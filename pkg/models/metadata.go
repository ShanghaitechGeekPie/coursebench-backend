package models

import (
	"coursebench-backend/pkg/modelRegister"
	"gorm.io/gorm"
)

type Metadata struct {
	gorm.Model
	DBVersion int
}

func init() {
	modelRegister.Register(&Metadata{})
}
