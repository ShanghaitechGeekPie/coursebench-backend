package teacher

import (
	"coursebench-backend/pkg/errors"
	"coursebench-backend/pkg/models"
	"coursebench-backend/pkg/queries"
	"github.com/gofiber/fiber/v2"
)

type AddRequest struct {
	Name         string `json:"name"`
	Job          string `json:"job"`
	Introduction string `json:"introduction"`
}

func Add(c *fiber.Ctx) (err error) {
	c.Accepts("application/json")
	var request AddRequest
	if err = c.BodyParser(&request); err != nil {
		return errors.Wrap(err, errors.InvalidArgument)
	}

	teacher, err := queries.AddTeacher(request.Name, request.Job, request.Introduction)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(models.OKResponse{
		Data:  teacher.ID,
		Error: false,
	})
}
