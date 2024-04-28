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
	var comments []models.Comment
	result := db.Preload("User").Preload("CourseGroup").Preload("CourseGroup.Course").Preload("CourseGroup.Teachers").
		Order("update_time DESC").Offset((id - 1) * 30).Limit(31).Find(&comments)
	if err := result.Error; err != nil {
		return errors.Wrap(err, errors.DatabaseError)
	}
	if result.RowsAffected == 31 {
		comments = comments[:30]
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
		Data: fiber.Map{
			"has_more": result.RowsAffected == 31,
			"comments": response,
		},
		Error: false,
	})
}
