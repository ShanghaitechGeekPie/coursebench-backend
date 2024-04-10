package users

import (
	"coursebench-backend/internal/middlewares/session"
	"coursebench-backend/pkg/errors"
	"coursebench-backend/pkg/models"
	"coursebench-backend/pkg/queries"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func Profile(c *fiber.Ctx) error {
	id_s := c.Params("id", "GG")
	id_i, err := strconv.Atoi(id_s)
	if err != nil {
		return errors.Wrap(err, errors.InvalidArgument)
	}
	id := uint(id_i)
	uid, err := session.GetUserID(c)
	if err != nil && !errors.Is(err, errors.UserNotLogin) {
		return err
	}
	response, err := queries.GetProfile(nil, id, uid)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(models.OKResponse{
		Data:  response,
		Error: false,
	})
}
