package users

import (
	"coursebench-backend/pkg/errors"
	"coursebench-backend/pkg/models"
	"coursebench-backend/pkg/queries"
	"github.com/gofiber/fiber/v2"
)

type RegisterActiveRequest struct {
	ID   uint   `json:"id"`
	Code string `json:"code"`
}

func RegisterActive(c *fiber.Ctx) (err error) {
	c.Accepts("application/json")
	var activeReq RegisterActiveRequest
	if err = c.BodyParser(&activeReq); err != nil {
		return errors.Wrap(err, errors.InvalidArgument)
	}
	err = queries.RegisterActive(nil, activeReq.ID, activeReq.Code)
	if errors.Is(err, errors.UserNotExists) {
		return errors.New(errors.MailCodeInvalid)
	}
	if err != nil {
		return
	}

	return c.Status(fiber.StatusOK).JSON(models.OKResponse{
		Data:  map[string]interface{}{"OK": true},
		Error: false,
	})
}
