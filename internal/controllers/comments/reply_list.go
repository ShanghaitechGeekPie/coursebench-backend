package comments

import (
	"coursebench-backend/internal/middlewares/session"
	"coursebench-backend/pkg/database"
	"coursebench-backend/pkg/errors"
	"coursebench-backend/pkg/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func ReplyList(c *fiber.Ctx) (err error) {
	commentIDString := c.Params("id", "0")
	commentIDInt, err := strconv.Atoi(commentIDString)
	if err != nil || commentIDInt <= 0 {
		return errors.New(errors.InvalidArgument)
	}
	commentID := uint(commentIDInt)

	uid, err := session.GetUserID(c)
	if err != nil {
		uid = 0
	}

	sortBy := c.Query("sort", "latest")
	orderBy, err := getReplyOrderBy(sortBy)
	if err != nil {
		return err
	}

	showAll := c.Query("all", "0") == "1"

	db := database.GetDB()

	var totalCount int64
	err = db.Model(&models.Reply{}).Where("comment_id = ?", commentID).Count(&totalCount).Error
	if err != nil {
		return errors.Wrap(err, errors.DatabaseError)
	}

	query := db.Preload("User").Preload("ParentReply").Preload("ParentReply.User").
		Where("comment_id = ?", commentID)
	if !showAll {
		query = query.Where("\"like\" > ?", 5)
	}

	var replies []models.Reply
	err = query.Order(orderBy).Find(&replies).Error
	if err != nil {
		return errors.Wrap(err, errors.DatabaseError)
	}

	responses, err := buildReplyResponses(db, replies, uid)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(models.OKResponse{
		Data: ReplyListResponse{
			TotalCount:    totalCount,
			FilteredCount: int64(len(replies)),
			Replies:       responses,
		},
		Error: false,
	})
}
