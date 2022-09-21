package teacher

import (
	"coursebench-backend/pkg/database"
	"coursebench-backend/pkg/errors"
	"coursebench-backend/pkg/models"
	"github.com/gofiber/fiber/v2"
)

type AllResponse struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	Institute    string `json:"institute"`
	Photo        string `json:"photo"`
	Job          string `json:"job"`
	Introduction string `json:"introduction"`
}

func All(c *fiber.Ctx) error {
	db := database.GetDB()
	var teachers []models.Teacher
	result := db.Find(&teachers)
	if result.Error != nil {
		return errors.Wrap(result.Error, errors.DatabaseError)
	}

	var response []AllResponse
	for _, v := range teachers {
		response = append(response, AllResponse{
			ID:           v.ID,
			Name:         v.Name,
			Institute:    v.Institute,
			Job:          v.Job,
			Introduction: v.Introduction,
			Photo:        v.Photo,
		})
	}

	return c.Status(fiber.StatusOK).JSON(models.OKResponse{
		Data:  response,
		Error: false,
	})
}
