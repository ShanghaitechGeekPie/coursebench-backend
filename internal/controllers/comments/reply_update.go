package comments

import (
	"coursebench-backend/internal/middlewares/session"
	"coursebench-backend/pkg/database"
	"coursebench-backend/pkg/errors"
	"coursebench-backend/pkg/models"
	"coursebench-backend/pkg/queries"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ReplyUpdateRequest struct {
	ID      uint   `json:"id"`
	Content string `json:"content"`
}

func ReplyUpdate(c *fiber.Ctx) (err error) {
	c.Accepts("application/json")
	var request ReplyUpdateRequest
	if err = c.BodyParser(&request); err != nil {
		return errors.Wrap(err, errors.InvalidArgument)
	}
	if !queries.CheckReplyContent(request.Content) {
		return errors.New(errors.InvalidArgument)
	}

	uid, err := session.GetUserID(c)
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
	if uid != reply.UserID {
		return errors.New(errors.PermissionDenied)
	}

	reply.Content = request.Content
	reply.UpdateTime = int(time.Now().Unix())
	err = db.Select("Content", "UpdateTime").Save(reply).Error
	if err != nil {
		return errors.Wrap(err, errors.DatabaseError)
	}

	return c.Status(fiber.StatusOK).JSON(models.OKResponse{Data: map[string]interface{}{}, Error: false})
}
