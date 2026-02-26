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

type ReplyDeleteRequest struct {
	ID uint `json:"id"`
}

func ReplyDelete(c *fiber.Ctx) (err error) {
	c.Accepts("application/json")
	var request ReplyDeleteRequest
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

	db := database.GetDB()
	reply := &models.Reply{}
	err = db.Where("id = ?", request.ID).Take(reply).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New(errors.InvalidArgument)
		}
		return errors.Wrap(err, errors.DatabaseError)
	}
	if !user.IsAdmin && uid != reply.UserID {
		return errors.New(errors.PermissionDenied)
	}

	err = db.Delete(reply).Error
	if err != nil {
		return errors.Wrap(err, errors.DatabaseError)
	}

	return c.Status(fiber.StatusOK).JSON(models.OKResponse{Data: nil, Error: false})
}
