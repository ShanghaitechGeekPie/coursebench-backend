package reward

import (
	"coursebench-backend/pkg/models"
	"coursebench-backend/pkg/queries"

	"github.com/gofiber/fiber/v2"
)

func Ranklist(c *fiber.Ctx) error {
	var response []models.RanklistResponse
	response, err := queries.Ranklist(nil)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(models.OKResponse{
		Data:  response,
		Error: false,
	})
}
