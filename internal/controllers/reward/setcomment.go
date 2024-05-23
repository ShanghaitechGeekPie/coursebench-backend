package reward

import (
	"coursebench-backend/internal/middlewares/session"
	"coursebench-backend/pkg/database"
	"coursebench-backend/pkg/errors"
	"coursebench-backend/pkg/models"
	"coursebench-backend/pkg/queries"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

	err = db.Transaction(func(tx *gorm.DB) error {
		// lock comment & user's reward
		comment := &models.Comment{}
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(comment, request.ID).Error
		if err != nil {
			return err
		}

		user := &models.User{}
		err = tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(user, comment.UserID).Error
		if err != nil {
			return err
		}

		// update reward
		user.Reward -= comment.Reward
		user.Reward += request.Reward
		comment.Reward = request.Reward

		// save
		if err := tx.Select("Reward").Save(user).Error; err != nil {
			return err
		}
		if err := tx.Select("Reward").Save(comment).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return errors.Wrap(err, errors.DatabaseError)
	}

	return c.Status(fiber.StatusOK).JSON(models.OKResponse{
		Error: false,
	})
}
