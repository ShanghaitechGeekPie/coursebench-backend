package course

import (
	"coursebench-backend/pkg/errors"
	"coursebench-backend/pkg/models"
	"coursebench-backend/pkg/queries"
	"github.com/gofiber/fiber/v2"
)

type AddRequest struct {
	Name      string `json:"name"`
	Institute string `json:"institute"`
	Credit    int    `json:"credit"`
	Code      string `json:"code"`
}

func Add(c *fiber.Ctx) (err error) {
	c.Accepts("application/json")
	var request AddRequest
	if err = c.BodyParser(&request); err != nil {
		return errors.Wrap(err, errors.InvalidArgument)
	}

	course, err := queries.AddCourse(nil, request.Name, request.Institute, request.Credit, request.Code)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(models.OKResponse{
		Data:  course.ID,
		Error: false,
	})
}
