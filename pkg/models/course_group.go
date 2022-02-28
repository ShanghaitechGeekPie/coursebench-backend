package models

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type CourseGroup struct {
	gorm.Model
	Code         string
	Course       Course
	CourseID     uint          `gorm:"index"`
	Scores       pq.Int64Array `gorm:"type:bigint[]"`
	CommentCount int
	Teachers     []*Teacher `gorm:"many2many:coursegroup_teachers"`
	Comment      []Comment
}
