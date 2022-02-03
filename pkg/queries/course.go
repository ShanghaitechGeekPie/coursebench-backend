package queries

import (
	"coursebench-backend/pkg/database"
	"coursebench-backend/pkg/errors"
	"coursebench-backend/pkg/models"
	"github.com/lib/pq"
	"gorm.io/gorm"
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
	course := &models.Course{}
	err = db.Where("id = ?", courseID).First(course).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.Wrap(err, errors.DatabaseError)
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.Wrap(err, errors.CourseNotExists)
	}
	course.Groups = append(course.Groups, *courseGroup)
	result := db.Save(course)

	//result := db.Create(courseGroup)
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

func AllCourseRequest() (Courses []models.CourseAllResponse, err error) {
	c, err := GetAllCourse()
	if err != nil {
		return nil, err
	}

	Courses = make([]models.CourseAllResponse, len(c))
	for i, v := range c {
		score := 0.0
		if v.CommentCount != 0 {
			score = float64(v.Scores[0]) / float64(v.CommentCount)
		}
		Courses[i] = models.CourseAllResponse{
			ID:        int(v.ID),
			Name:      v.Name,
			Institute: v.Institute,
			Code:      v.Code,
			Score:     score,
		}
	}

	return
}
