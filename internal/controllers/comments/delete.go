package comments

import (
	"coursebench-backend/internal/middlewares/session"
	"coursebench-backend/pkg/database"
	"coursebench-backend/pkg/events"
	"coursebench-backend/pkg/models"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type DeleteRequest struct {
	ID uint `json:"id"`
}

func Delete(c *fiber.Ctx) *events.AttributedEvent {
	c.Accepts("application/json")
	var request DeleteRequest
	if err := c.BodyParser(&request); err != nil {
		return events.Wrap(err, events.InvalidArgument)
	}

	uid, err := session.GetUserID(c)
	if err != nil {
		return err
	}

	db := database.GetDB()
	comment := &models.Comment{}
	if err := db.Preload("CourseGroup").Preload("CourseGroup.Course").Where("id = ?", request.ID).Take(comment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return events.New(events.CommentNotExists)
		} else {
			return events.Wrap(err, events.DatabaseError)
		}
	}

	if uid != comment.UserID {
		return events.New(events.PermissionDenied)
	}

	err = events.Wrap(db.Transaction(func(tx *gorm.DB) error {
		comment.CourseGroup.CommentCount--
		for i := 0; i < models.ScoreLength; i++ {
			comment.CourseGroup.Scores[i] -= comment.Scores[i]
		}
		comment.CourseGroup.Course.CommentCount--
		for i := 0; i < models.ScoreLength; i++ {
			comment.CourseGroup.Course.Scores[i] -= comment.Scores[i]
		}
		err := tx.Select("Scores", "CommentCount").Updates(comment.CourseGroup).Error
		if err != nil {
			return events.Wrap(err, events.DatabaseError).ToError()
		}

		err = tx.Select("Scores", "CommentCount").Updates(comment.CourseGroup.Course).Error
		if err != nil {
			return events.Wrap(err, events.DatabaseError).ToError()
		}

		err = tx.Delete(comment).Error
		if err != nil {
			return events.Wrap(err, events.DatabaseError).ToError()
		}

		return nil
	}), events.DatabaseError)

	return events.Wrap(c.Status(fiber.StatusOK).JSON(models.OKResponse{
		Data:  nil,
		Error: false,
	}), events.InternalServerError)
}
