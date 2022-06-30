package users

import (
	"coursebench-backend/pkg/models"
	"coursebench-backend/pkg/queries"
	"github.com/gofiber/fiber/v2"
)

func GetCaptcha(c *fiber.Ctx) (err error) {
	c.Accepts("application/json")
	s, err := queries.GenerateCaptcha(c)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(models.OKResponse{
		Data:  map[string]interface{}{"img": s},
		Error: false,
	})
}
