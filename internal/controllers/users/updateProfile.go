package users

import (
	"coursebench-backend/internal/middlewares/session"
	"coursebench-backend/pkg/errors"
	"coursebench-backend/pkg/models"
	"coursebench-backend/pkg/queries"
	"github.com/gofiber/fiber/v2"
)

type UpdateProfileRequest struct {
	Year        int              `json:"year"`
	Grade       models.GradeType `json:"grade"`
	NickName    string           `json:"nickname"`
	RealName    string           `json:"realname"`
	IsAnonymous bool             `json:"is_anonymous"`
}

func UpdateProfile(c *fiber.Ctx) (err error) {
	c.Accepts("application/json")
	var request UpdateProfileRequest
	if err = c.BodyParser(&request); err != nil {
		return errors.Wrap(err, errors.InvalidArgument)
	}
	id, err := session.GetUserID(c)
	if err != nil {
		return err
	}
	err = queries.UpdateProfile(id, request.Year, request.Grade, request.NickName, request.RealName, request.IsAnonymous)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(models.OKResponse{
		Data:  nil,
		Error: false,
	})
}
