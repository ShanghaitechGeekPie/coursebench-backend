package comments

import (
	"coursebench-backend/internal/middlewares/session"
	"coursebench-backend/internal/utils"
	"coursebench-backend/pkg/database"
	"coursebench-backend/pkg/errors"
	"coursebench-backend/pkg/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type RecentCommentByPageResponse struct {
	PageCount int64             `json:"page_count"`
	HasMore   bool              `json:"has_more"`
	Comments  []CommentResponse `json:"comments"`
}

func RecentCommentByPage(c *fiber.Ctx) (err error) {
	uid, err := session.GetUserID(c)
	if err != nil {
		uid = 0
	}

	id_s := c.Params("id", "1")
	id, err := strconv.Atoi(id_s)
	if err != nil {
		return errors.New(errors.InvalidArgument)
	}

	db := database.GetDB()

	var count int64
	err = db.Model(&models.Comment{}).Count(&count).Error
	if err != nil {
		return errors.Wrap(err, errors.DatabaseError)
	}

	var comments []models.Comment
	err = db.Preload("User").Preload("CourseGroup").Preload("CourseGroup.Course").Preload("CourseGroup.Teachers").
		Order("update_time DESC").Offset((id - 1) * 30).Limit(30).Find(&comments).Error
	if err != nil {
		return errors.Wrap(err, errors.DatabaseError)
	}

	var likeResult []CommentLikeResult
	if uid != 0 {
		db.Raw("SELECT comment_likes.comment_id, comment_likes.is_like from comments, comment_likes where comment_likes.user_id = ? and comment_likes.comment_id = comments.id and comment_likes.deleted_at is NULL and comments.deleted_at is NULL order by create_time desc OFFSET ? LIMIT 30",
			(id-1)*30,
			uid).Scan(&likeResult)
	}
	var response []CommentResponse
	response = GenerateResponse(comments, uid, likeResult, true, utils.GetIP(c))
	return c.Status(fiber.StatusOK).JSON(models.OKResponse{
		Data: RecentCommentByPageResponse{
			PageCount: (count + 29) / 30, // ceil
			HasMore:   count > int64(id)*30,
			Comments:  response,
		},
		Error: false,
	})
}
