package users

import (
	"coursebench-backend/internal/config"
	"coursebench-backend/pkg/errors"
	"coursebench-backend/pkg/models"
	"coursebench-backend/pkg/queries"
	"github.com/gofiber/fiber/v2"
)

type ResetPasswordRequest struct {
	Email   string `json:"email"`
	Captcha string `json:"captcha"`
}

func ResetPassword(c *fiber.Ctx) (err error) {
	c.Accepts("application/json")
	var userReq ResetPasswordRequest
	if err = c.BodyParser(&userReq); err != nil {
		return errors.Wrap(err, errors.InvalidArgument)
	}
	if !config.GlobalConf.DisableCaptcha {
		if err = queries.VerifyCaptcha(c, userReq.Captcha); err != nil {
			return
		}
	}

	if err = queries.ResetPassword(userReq.Email); err != nil {
		return
	}

	return c.Status(fiber.StatusOK).JSON(models.OKResponse{
		Data:  map[string]interface{}{"OK": true},
		Error: false,
	})
}
