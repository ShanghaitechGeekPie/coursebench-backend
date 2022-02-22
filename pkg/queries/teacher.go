package queries

import (
	"coursebench-backend/pkg/database"
	"coursebench-backend/pkg/errors"
	"coursebench-backend/pkg/models"
	"gorm.io/gorm"
)

func GetTeacher(id uint) (teacher *models.Teacher, err error) {
	db := database.GetDB()

	teacher = &models.Teacher{}
	result := db.Where("id = ?", id).Take(teacher)
	if err := result.Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.Wrap(err, errors.DatabaseError)
	}
	if result.RowsAffected == 0 {
		return nil, errors.New(errors.TeacherNotExists)
	}

	return teacher, nil
}

func AddTeacher(name string, job string, introduction string, EamsID int) (teacher *models.Teacher, err error) {
	db := database.GetDB()

	teacher = &models.Teacher{
		Name:         name,
		Job:          job,
		Introduction: introduction,
		EamsID:       EamsID,
	}
	result := db.Create(teacher)
	if err := result.Error; err != nil {
		return nil, errors.Wrap(err, errors.DatabaseError)
	}

	return teacher, nil
}
