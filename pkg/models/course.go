package models

import (
	"coursebench-backend/pkg/modelRegister"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Course struct {
	gorm.Model
	Teachers    []*Teacher `gorm:"many2many:course_teachers;"`
	Name        string
	Code        string
	Scores      pq.Int64Array `gorm:"type:bigint[]"`
	Evaluations []Evaluation
}

func init() {
	modelRegister.Register(&Course{})
}
