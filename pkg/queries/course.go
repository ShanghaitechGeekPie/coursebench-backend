package queries

import (
	"coursebench-backend/internal/controllers/courses"
	"coursebench-backend/pkg/database"
	"coursebench-backend/pkg/errors"
	"coursebench-backend/pkg/models"
	"github.com/lib/pq"
)

func AddCourse(name string, institute string, credit int, code string) (course *models.Course, err error) {
	db := database.GetDB()
	/*
		scores64 := make([]int64, len(scores))
		for i, v := range scores {
			scores64[i] = int64(v)
		}
		teachersT := make([]*models.Teacher, len(teachers))
		for i, v := range teachers {
			teachersT[i], err = GetTeacher(uint(v))
			if err != nil {
				return nil, err
			}
		}*/
	course = &models.Course{Name: name, Institute: institute, Credit: credit, Code: code, Scores: pq.Int64Array([]int64{0, 0, 0, 0}), CommentCount: 0}
	result := db.Create(course)
	if result.Error != nil {
		return nil, errors.Wrap(result.Error, errors.DatabaseError)
	}
	return
}

func AddCourseGroup(code string, courseID int, teachers []int) (courseGroup *models.CourseGroup, err error) {
	db := database.GetDB()
	teachersT := make([]*models.Teacher, len(teachers))
	for i, v := range teachers {
		teachersT[i], err = GetTeacher(uint(v))
		if err != nil {
			return nil, err
		}
	}
	courseGroup = &models.CourseGroup{CourseID: uint(courseID), Code: code, Teachers: teachersT, Scores: pq.Int64Array{0, 0, 0, 0}, CommentCount: 0}
	result := db.Create(courseGroup)
	if result.Error != nil {
		return nil, errors.Wrap(result.Error, errors.DatabaseError)
	}
	return
}
func GetAllCourse() (courses []*models.Course, err error) {
	db := database.GetDB()
	result := db.Find(&courses)
	if result.Error != nil {
		return nil, errors.Wrap(result.Error, errors.DatabaseError)
	}
	return
}

func AllCourseRequest() (Courses []courses.CourseAllResponse, err error) {
	c, err := GetAllCourse()
	if err != nil {
		return nil, err
	}
	Courses = make([]courses.CourseAllResponse, len(c))
	/*
		for i, v := range c {
			Courses = append(Courses, courses.CourseAllResponse{
				ID:        v.ID,
				Name:      v.Name,
				Institute: v.Institute,
				Credit:    v.Credit,
				Code:      v.Code,
				Group:     v.Group,
				Scores:    v.Scores,
				Teachers:  v.Teachers,
			})
		}*/
	return
}
