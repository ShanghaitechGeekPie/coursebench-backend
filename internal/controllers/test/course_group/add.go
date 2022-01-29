package course_group

import (
	"coursebench-backend/pkg/errors"
	"coursebench-backend/pkg/models"
	"coursebench-backend/pkg/queries"
	"github.com/gofiber/fiber/v2"
)

type AddRequest struct {
	Code     string `json:"code"`
	Course   int    `json:"course"`
	Teachers []int  `json:"teachers"`
}

func Add(c *fiber.Ctx) (err error) {
	c.Accepts("application/json")
	var request AddRequest
	if err = c.BodyParser(&request); err != nil {
		return errors.Wrap(err, errors.InvalidArgument)
	}

	courseGroup, err := queries.AddCourseGroup(request.Code, request.Course, request.Teachers)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(models.OKResponse{
		Data:  courseGroup.ID,
		Error: false,
	})
}
