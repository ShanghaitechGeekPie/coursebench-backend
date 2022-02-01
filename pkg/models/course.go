package models

import (
	"coursebench-backend/pkg/modelRegister"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Course struct {
	gorm.Model
	Name         string
	Institute    string
	Credit       int
	Code         string
	Scores       pq.Int64Array `gorm:"type:bigint[]"`
	CommentCount int
	Groups       []CourseGroup `gorm:"foreignKey:CourseID"`
}

func init() {
	modelRegister.Register(&Course{})
}
