package models

import (
	"coursebench-backend/pkg/modelRegister"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Course struct {
	gorm.Model
	Teachers    []Teacher
	Name        string
	Code        string
	Scores      pq.Int64Array
	Evaluations []Evaluation
}

func init() {
	modelRegister.Register(&Course{})
}
