package courses

import (
	"coursebench-backend/pkg/models"
	"github.com/gofiber/fiber/v2"
)

type CourseAllResponse struct {
	Name      string  `json:"name"`
	Code      string  `json:"code"`
	ID        int     `json:"id"`
	Score     float64 `json:"score"`
	Institute string  `json:"institute"`
}

func CourseAll(c *fiber.Ctx) error {
	var response []CourseAllResponse
	return c.Status(fiber.StatusOK).JSON(models.OKResponse{
		Data:  response,
		Error: false,
	})
}
