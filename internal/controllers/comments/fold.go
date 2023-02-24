package comments

import (
	"coursebench-backend/internal/middlewares/session"
	"coursebench-backend/pkg/database"
	"coursebench-backend/pkg/errors"
	"coursebench-backend/pkg/models"
	"coursebench-backend/pkg/queries"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type FoldRequest struct {
	ID     int  `json:"id"`
	Status bool `json:"status"`
}

func Fold(c *fiber.Ctx) (err error) {
	c.Accepts("application/json")
	var request FoldRequest
	if err = c.BodyParser(&request); err != nil {
		return errors.Wrap(err, errors.InvalidArgument)
	}

	uid, err := session.GetUserID(c)
	if err != nil {
		return err
	}

	user, err := queries.GetUserByID(nil, uid)
	if err != nil {
		return err
	}
	if !user.IsAdmin {
		return errors.New(errors.PermissionDenied)
	}

	db := database.GetDB()
	comment := &models.Comment{}
	err = db.Where("id = ?", request.ID).Take(comment).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New(errors.CommentNotExists)
		} else {
			return errors.Wrap(err, errors.DatabaseError)
		}
	}

	err = db.Model(comment).Where("id = ?", request.ID).Update("is_fold", request.Status).Error
	if err != nil {
		return errors.Wrap(err, errors.DatabaseError)
	}

	return c.Status(fiber.StatusOK).JSON(models.OKResponse{
		Data:  map[string]interface{}{},
		Error: false})
}
