package comments

import (
	"coursebench-backend/internal/middlewares/session"
	"coursebench-backend/pkg/database"
	"coursebench-backend/pkg/errors"
	"coursebench-backend/pkg/models"
	"coursebench-backend/pkg/queries"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ReplyPostRequest struct {
	ParentReplyID *uint  `json:"parent_reply_id"`
	Content       string `json:"content"`
	IsAnonymous   bool   `json:"is_anonymous"`
}

func ReplyPost(c *fiber.Ctx) (err error) {
	commentIDString := c.Params("id", "0")
	commentIDInt, err := strconv.Atoi(commentIDString)
	if err != nil || commentIDInt <= 0 {
		return errors.New(errors.InvalidArgument)
	}
	commentID := uint(commentIDInt)

	c.Accepts("application/json")
	var request ReplyPostRequest
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
	postTime := int(time.Now().Unix())

	reply := &models.Reply{
		CommentID:     commentID,
		ParentReplyID: request.ParentReplyID,
		UserID:        uid,
		Content:       request.Content,
		IsAnonymous:   request.IsAnonymous,
		Like:          0,
		Dislike:       0,
		CreateTime:    postTime,
		UpdateTime:    postTime,
		IsFold:        false,
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		comment := &models.Comment{}
		err := tx.Where("id = ?", commentID).Take(comment).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New(errors.CommentNotExists)
			}
			return errors.Wrap(err, errors.DatabaseError)
		}

		if request.ParentReplyID != nil {
			parentReply := &models.Reply{}
			err = tx.Where("id = ?", *request.ParentReplyID).Take(parentReply).Error
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return errors.New(errors.InvalidArgument)
				}
				return errors.Wrap(err, errors.DatabaseError)
			}
			if parentReply.CommentID != commentID {
				return errors.New(errors.InvalidArgument)
			}
		}

		err = tx.Create(reply).Error
		if err != nil {
			return errors.Wrap(err, errors.DatabaseError)
		}
		return nil
	})
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(models.OKResponse{
		Data:  map[string]interface{}{"reply_id": reply.ID},
		Error: false,
	})
}
