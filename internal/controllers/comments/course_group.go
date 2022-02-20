package comments

import (
	"coursebench-backend/internal/middlewares/session"
	"coursebench-backend/pkg/database"
	"coursebench-backend/pkg/events"
	"coursebench-backend/pkg/models"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func CourseGroupComment(c *fiber.Ctx) *events.AttributedEvent {
	ctx := c.UserContext()
	id_s := c.Params("id", "GG")
	id_i, err := strconv.Atoi(id_s)
	if err != nil {
		return events.New(events.InvalidArgument).Log(ctx)
	}
	id := uint(id_i)

	uid, event := session.GetUserID(c)
	if event != nil {
		uid = 0
	}

	db := database.GetDB()
	var comments []models.Comment
	err = db.Preload("User").Preload("CourseGroup").Preload("CourseGroup.Course").Preload("CourseGroup.Teachers").
		Where("course_group_id = ?", id).Find(&comments).Error
	if err != nil {
		return events.Wrap(err, events.DatabaseError).Log(ctx)
	}
	var response []CommentResponse
	response = GenerateResponse(comments, uid)
	return events.Wrap(c.Status(fiber.StatusOK).JSON(models.OKResponse{
		Data:  response,
		Error: false,
	}), events.InternalServerError)
}
