package users

import (
	"coursebench-backend/internal/middlewares/session"
	"coursebench-backend/pkg/errors"
	"coursebench-backend/pkg/models"
	"github.com/gofiber/fiber/v2"
)

func Logout(c *fiber.Ctx) (err error) {
	sess, err := session.GetStore().Get(c)
	if err != nil {
		return errors.Wrap(err, errors.InternalServerError)
	}
	// Save session
	if err := sess.Destroy(); err != nil {
		return errors.Wrap(err, errors.InternalServerError)
	}

	return c.Status(fiber.StatusOK).JSON(models.OKResponse{
		Data:  map[string]interface{}{},
		Error: false,
	})
}
