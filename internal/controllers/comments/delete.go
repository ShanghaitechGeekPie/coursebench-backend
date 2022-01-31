package comments

import (
	"coursebench-backend/pkg/errors"
	"coursebench-backend/pkg/models"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func Delete(c *fiber.Ctx) (err error) {
	id_s := c.Params("id", "GG")
	id_i, err := strconv.Atoi(id_s)
	if err != nil {
		return errors.Wrap(err, errors.InvalidArgument)
	}
	id := uint(id_i)
	return c.Status(fiber.StatusOK).JSON(models.OKResponse{
		Data:  id,
		Error: false,
	})
}
