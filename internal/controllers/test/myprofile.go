package test

import (
	"coursebench-backend/internal/middlewares/session"
	"coursebench-backend/pkg/models"
	"github.com/gofiber/fiber/v2"
)

func MyProfile(c *fiber.Ctx) (err error) {
	id, err := session.GetUserID(c)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(models.OKResponse{
		Data:  map[string]interface{}{"UserID": id},
		Error: false,
	})
}
