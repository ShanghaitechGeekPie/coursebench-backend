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
	CourseID            uint `gorm:"index"`
	Semester            int
	Scores              pq.Int64Array `gorm:"type:bigint[]"`
	Title               string
	Content             string
	StudentScoreRanking int
	IsAnonymous         bool
	CreateTime          int
	UpdateTime          int
	Like                int
	Dislike             int
	IsFold              bool `gorm:"default:false"`
	IsCovered           bool `gorm:"default:false"`
	CoverTitle          string
	CoverContent        string
	CoverReason         string
	Reward			    int  `gorm:"default:0"`
}

type CommentLike struct {
	gorm.Model
	UserID    uint `gorm:"index"`
	CommentID uint `gorm:"index"`
	IsLike    bool `gorm:"index"`
}

func init() {
	modelRegister.Register(&Comment{})
	modelRegister.Register(&CommentLike{})
}
