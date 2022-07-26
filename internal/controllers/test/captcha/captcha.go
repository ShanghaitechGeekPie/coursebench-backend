package captcha

import (
	"coursebench-backend/pkg/errors"
	"coursebench-backend/pkg/models"
	"coursebench-backend/pkg/queries"
	"github.com/gofiber/fiber/v2"
)

type VerifyRequest struct {
	Code string `json:"code"`
}

func Get(c *fiber.Ctx) (err error) {
	c.Accepts("application/json")
	s, err := queries.GenerateCaptcha(c)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(models.OKResponse{
		Data:  s,
		Error: false,
	})
}

func Verify(c *fiber.Ctx) (err error) {
	c.Accepts("application/json")
	var request VerifyRequest
	if err = c.BodyParser(&request); err != nil {
		return errors.Wrap(err, errors.InvalidArgument)
	}
	return c.Status(fiber.StatusOK).JSON(models.OKResponse{
		Data:  queries.VerifyCaptcha(c, request.Code),
		Error: false,
	})
}
