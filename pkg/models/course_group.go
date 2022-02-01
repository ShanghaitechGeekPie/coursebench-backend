package models

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type CourseGroup struct {
	gorm.Model
	Code string
	//Course       Course
	CourseID     uint
	Scores       pq.Int64Array `gorm:"type:bigint[]"`
	CommentCount int
	Teachers     []*Teacher `gorm:"many2many:course_teachers"`
	Comment      []Comment
}
