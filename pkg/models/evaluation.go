package models

import (
	"coursebench-backend/pkg/modelRegister"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Evaluation struct {
	gorm.Model
	User                User
	Course              Course
	Semester            int
	Scores              pq.Int64Array
	Title               string
	Comment             string
	StudentScoreNumber  int
	StudentScoreRanking int
	IsAnonymous         bool
}

func init() {
	modelRegister.Register(&Evaluation{})
}
