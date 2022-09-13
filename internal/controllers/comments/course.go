package comments

import (
	"coursebench-backend/internal/middlewares/session"
	"coursebench-backend/pkg/database"
	"coursebench-backend/pkg/errors"
	"coursebench-backend/pkg/models"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func CourseComment(c *fiber.Ctx) (err error) {
	id_s := c.Params("id", "GG")
	id_i, err := strconv.Atoi(id_s)
	if err != nil {
		return errors.New(errors.InvalidArgument)
	}
	id := uint(id_i)

	uid, err := session.GetUserID(c)
	if err != nil {
		uid = 0
	}

	db := database.GetDB()
	var comments []models.Comment
	err = db.Preload("User").Preload("CourseGroup").Preload("CourseGroup.Course").Preload("CourseGroup.Teachers").
		Where("course_id = ?", id).Find(&comments).Error
	if err != nil {
		return errors.Wrap(err, errors.DatabaseError)
	}
	var response []CommentResponse
	response = GenerateResponse(comments, uid)
	return c.Status(fiber.StatusOK).JSON(models.OKResponse{
		Data:  response,
		Error: false,
	})
}
