package users

import (
	"coursebench-backend/internal/middlewares/session"
	"coursebench-backend/pkg/errors"
	"coursebench-backend/pkg/models"
	"github.com/gofiber/fiber/v2"
)

// MyID returns the ID of the current user.
// It is used to check if the user is logged in.
func MyID(c *fiber.Ctx) (err error) {
	id, err := session.GetUserID(c)
	if err != nil {
		if errors.Is(err, errors.UserNotLogin) {
			id = 0
		} else {
			return err
		}
	}
	return c.Status(fiber.StatusOK).JSON(models.OKResponse{
		Data:  map[string]interface{}{"id": id},
		Error: false,
	})
}
