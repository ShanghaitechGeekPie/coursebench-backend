package queries

import (
	"coursebench-backend/pkg/database"
	"coursebench-backend/pkg/errors"
	"coursebench-backend/pkg/models"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

func AddCourse(db *gorm.DB, name string, institute string, credit int, code string) (course *models.Course, err error) {
	if db == nil {
		db = database.GetDB()
	}
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

func AddCourseGroup(db *gorm.DB, code string, courseID int, teachers []int) (courseGroup *models.CourseGroup, err error) {
	if db == nil {
		db = database.GetDB()
	}
	teachersT := make([]*models.Teacher, len(teachers))
	for i, v := range teachers {
		teachersT[i], err = GetTeacher(db, uint(v))
		if err != nil {
			return nil, err
		}
	}
	courseGroup = &models.CourseGroup{CourseID: uint(courseID), Code: code, Teachers: teachersT, Scores: pq.Int64Array{0, 0, 0, 0}, CommentCount: 0}
	err = db.Transaction(func(tx *gorm.DB) error {
		course := &models.Course{}
		result := tx.Save(courseGroup)
		if result.Error != nil {
			return errors.Wrap(result.Error, errors.DatabaseError)
		}
		err = tx.Preload("Teachers").Preload("Groups").Where("id = ?", courseID).First(course).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.Wrap(err, errors.DatabaseError)
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.Wrap(err, errors.CourseNotExists)
		}
		course.Groups = append(course.Groups, *courseGroup)

		// Update course-teacher relation
		// Should check if teacher is already in course
		// insert into course.teachers

		for _, teacher := range teachersT {
			flag := true
			for _, t := range course.Teachers {
				if t.ID == teacher.ID {
					flag = false
					break
				}
			}
			if flag {
				course.Teachers = append(course.Teachers, teacher)
			}
		}

		result = tx.Select("Teachers", "Groups").Save(course)

		if result.Error != nil {
			return errors.Wrap(result.Error, errors.DatabaseError)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return
}
func GetAllCourse(db *gorm.DB) (courses []*models.Course, err error) {
	if db == nil {
		db = database.GetDB()
	}
	result := db.Find(&courses)
	if result.Error != nil {
		return nil, errors.Wrap(result.Error, errors.DatabaseError)
	}
	return
}

func AllCourseRequest(db *gorm.DB) (Courses []models.CourseAllResponse, err error) {
	if db == nil {
		db = database.GetDB()
	}
	c, err := GetAllCourse(db)
	if err != nil {
		return nil, err
	}

	Courses = make([]models.CourseAllResponse, len(c))
	for i, v := range c {
		score := make([]float64, models.ScoreLength)
		if v.CommentCount != 0 {
			for j := 0; j < models.ScoreLength; j++ {
				score[j] = float64(v.Scores[j]) / float64(v.CommentCount)
			}
		}
		Courses[i] = models.CourseAllResponse{
			ID:         int(v.ID),
			Name:       v.Name,
			Institute:  v.Institute,
			Code:       v.Code,
			Score:      score,
			Credit:     v.Credit,
			CommentNum: v.CommentCount,
		}
	}

	return
}
