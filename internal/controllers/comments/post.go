package comments

import (
	"coursebench-backend/pkg/models"
	"github.com/gofiber/fiber/v2"
)

func Post(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(models.OKResponse{
		Data:  nil,
		Error: false,
	})
}
