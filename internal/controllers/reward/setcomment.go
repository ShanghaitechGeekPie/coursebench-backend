package reward

import (
	"coursebench-backend/internal/middlewares/session"
	"coursebench-backend/pkg/database"
	"coursebench-backend/pkg/errors"
	"coursebench-backend/pkg/models"
	"coursebench-backend/pkg/queries"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type SetCommentRequest struct {
	Reward int `json:"reward"`
}

func SetComment(c *fiber.Ctx) error {
	c.Accepts("application/json")
	var request SetCommentRequest
	if err := c.BodyParser(&request); err != nil {
		return errors.Wrap(err, errors.InvalidArgument)
	}

	idRaw := c.Params("id", "GG")
	id, err := strconv.Atoi(idRaw)
	if err != nil {
		return errors.New(errors.InvalidArgument)
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
	result := db.First(&comment, id)
	if result.Error != nil {
		return errors.Wrap(result.Error, errors.DatabaseError)
	}

	comment.User.Reward -= comment.Reward
	comment.Reward = request.Reward
	comment.User.Reward += comment.Reward

	if err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(comment).Error; err != nil {
			return err
		}
		if err := tx.Save(user).Error; err != nil {
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
