package models

import (
	"coursebench-backend/pkg/modelRegister"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

const ScoreLength = 4

type Comment struct {
	gorm.Model
	UserID              uint `gorm:"index"`
	User                User
	CourseGroup         CourseGroup
	CourseGroupID       uint `gorm:"index"`
	Semester            int
	Scores              pq.Int64Array `gorm:"type:bigint[]"`
	Title               string
	Content             string
	StudentScoreRanking int
	IsAnonymous         bool
	CreateTime          int
	UpdateTime          int
}

func init() {
	modelRegister.Register(&Comment{})
}
