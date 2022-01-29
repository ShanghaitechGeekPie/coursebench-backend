package models

import (
	"coursebench-backend/pkg/modelRegister"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	UserID              uint
	User                User
	CourseGroup         CourseGroup
	CourseGroupID       uint
	Semester            int
	Scores              pq.Int64Array `gorm:"type:bigint[]"`
	Title               string
	Comment             string
	StudentScoreNumber  int
	StudentScoreRanking int
	IsAnonymous         bool
}

func init() {
	modelRegister.Register(&Comment{})
}
