package users

import (
	"coursebench-backend/internal/config"
	"coursebench-backend/internal/middlewares/session"
	"coursebench-backend/pkg/errors"
	"coursebench-backend/pkg/models"
	"coursebench-backend/pkg/queries"
	"github.com/gofiber/fiber/v2"
)

type UpdatePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
	Captcha     string `json:"captcha"`
}

func UpdatePassword(c *fiber.Ctx) (err error) {
	c.Accepts("application/json")
	var request UpdatePasswordRequest
	if err = c.BodyParser(&request); err != nil {
		return errors.Wrap(err, errors.InvalidArgument)
	}

	if !config.GlobalConf.DisableCaptchaAndMail && !queries.VerifyCaptcha(c, request.Captcha) {
		return errors.New(errors.CaptchaMismatch)
	}

	id, err := session.GetUserID(c)
	if err != nil {
		return err
	}
	err = queries.UpdatePassword(id, request.OldPassword, request.NewPassword)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(models.OKResponse{
		Data:  nil,
		Error: false,
	})
}
