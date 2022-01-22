package users

import (
	"coursebench-backend/pkg/errors"
	"coursebench-backend/pkg/models"
	"coursebench-backend/pkg/queries"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type ProfileResponse struct {
	Email    string           `json:"email"`
	Year     int              `json:"year"`
	Grade    models.GradeType `json:"grade"`
	NickName string           `json:"nickname"`
	RealName string           `json:"realname"`
}

func Profile(c *fiber.Ctx) error {
	id_s := c.Params("id", "GG")
	id_i, err := strconv.Atoi(id_s)
	if err != nil {
		return errors.Wrap(err, errors.InvalidArgument)
	}
	id := uint(id_i)
	user, err := queries.GetUserByID(id)
	response := ProfileResponse{Email: user.Email, Year: user.Year, Grade: user.Grade, NickName: user.NickName, RealName: user.RealName}
	return c.Status(fiber.StatusOK).JSON(models.OKResponse{
		Data:  response,
		Error: false,
	})
}
