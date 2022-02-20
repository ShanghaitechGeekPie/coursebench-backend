package courses

import (
	"coursebench-backend/pkg/events"
	"coursebench-backend/pkg/models"
	"coursebench-backend/pkg/queries"
	"github.com/gofiber/fiber/v2"
)

func All(c *fiber.Ctx) *events.AttributedEvent {
	var response []models.CourseAllResponse
	response, event := queries.AllCourseRequest()
	if event != nil && event.IsError() {
		return event
	}
	return events.Wrap(c.Status(fiber.StatusOK).JSON(models.OKResponse{
		Data:  response,
		Error: false,
	}), events.InternalServerError)
}
