package teacher

import (
	"coursebench-backend/pkg/database"
	"coursebench-backend/pkg/errors"
	"coursebench-backend/pkg/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"strconv"
)

type InfoResponse struct {
	ID           uint                       `json:"id"`
	Name         string                     `json:"name"`
	Institute    string                     `json:"institute"`
	Job          string                     `json:"job"`
	Introduction string                     `json:"introduction"`
	Courses      []models.CourseAllResponse `json:"courses"`
}

func Info(c *fiber.Ctx) (err error) {
	id_s := c.Params("id", "GG")
	id_i, err := strconv.Atoi(id_s)
	if err != nil {
		return errors.Wrap(err, errors.InvalidArgument)
	}
	id := uint(id_i)
	db := database.GetDB()
	teachers := &models.Teacher{}
	err = db.Preload("Courses").Where("id = ?", id).Take(&teachers).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.Wrap(err, errors.DatabaseError)
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.Wrap(err, errors.TeacherNotExists)
	}
	response := InfoResponse{
		ID:           teachers.ID,
		Name:         teachers.Name,
		Institute:    teachers.Institute,
		Job:          teachers.Job,
		Introduction: teachers.Introduction,
		Courses:      make([]models.CourseAllResponse, 0),
	}
	for _, v := range teachers.Courses {
		score := make([]float64, models.ScoreLength)
		if v.CommentCount != 0 {
			for j := 0; j < models.ScoreLength; j++ {
				score[j] = float64(v.Scores[j]) / float64(v.CommentCount)
			}
		}
		response.Courses = append(response.Courses, models.CourseAllResponse{
			ID:        int(v.ID),
			Name:      v.Name,
			Institute: v.Institute,
			Code:      v.Code,
			Score:     score,
		})
	}

	return c.Status(fiber.StatusOK).JSON(models.OKResponse{
		Data:  response,
		Error: false,
	})
}
