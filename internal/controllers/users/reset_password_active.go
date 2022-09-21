package users

import (
	"coursebench-backend/pkg/errors"
	"coursebench-backend/pkg/models"
	"coursebench-backend/pkg/queries"
	"github.com/gofiber/fiber/v2"
)

type ResetPasswordActiveRequest struct {
	ID       uint   `json:"id"`
	Code     string `json:"code"`
	Password string `json:"password"`
}

func ResetPasswordActive(c *fiber.Ctx) (err error) {
	c.Accepts("application/json")
	var activeReq ResetPasswordActiveRequest
	if err = c.BodyParser(&activeReq); err != nil {
		return errors.Wrap(err, errors.InvalidArgument)
	}
	err = queries.ResetPasswordActive(activeReq.ID, activeReq.Code, activeReq.Password)
	if err != nil {
		return
	}

	return c.Status(fiber.StatusOK).JSON(models.OKResponse{
		Data:  map[string]interface{}{"OK": true},
		Error: false,
	})
}
