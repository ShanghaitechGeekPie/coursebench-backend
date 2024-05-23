package reward

import (
	"coursebench-backend/internal/middlewares/session"
	"coursebench-backend/pkg/database"
	"coursebench-backend/pkg/errors"
	"coursebench-backend/pkg/models"
	"coursebench-backend/pkg/queries"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type SetCommentRequest struct {
	ID     int `json:"id"`
	Reward int `json:"reward"`
}

func SetComment(c *fiber.Ctx) error {
	c.Accepts("application/json")
	var request SetCommentRequest
	if err := c.BodyParser(&request); err != nil {
		return errors.Wrap(err, errors.InvalidArgument)
	}

	uid, err := session.GetUserID(c)
	if err != nil {
		uid = 0
	}

	db := database.GetDB()
	user, err := queries.GetUserByID(db, uid)
	if err != nil {
		return errors.Wrap(err, errors.DatabaseError)
	}
	if !(user.IsCommunityAdmin || user.IsAdmin) {
		return errors.New(errors.PermissionDenied)
	}

	comment := &models.Comment{}
	result := db.Preload("User").First(&comment, request.ID)
	if result.Error != nil {
		return errors.Wrap(result.Error, errors.DatabaseError)
	}

	comment.User.Reward -= comment.Reward
	comment.Reward = request.Reward
	comment.User.Reward += comment.Reward

	if err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Select("reward").Save(comment).Error; err != nil {
			return err
		}
		if err := tx.Select("reward").Save(&comment.User).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return errors.Wrap(err, errors.DatabaseError)
	}

	return c.Status(fiber.StatusOK).JSON(models.OKResponse{
		Error: false,
	})
}
