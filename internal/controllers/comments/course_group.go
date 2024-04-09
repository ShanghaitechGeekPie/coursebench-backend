package comments

import (
	"coursebench-backend/internal/middlewares/session"
	"coursebench-backend/internal/utils"
	"coursebench-backend/pkg/database"
	"coursebench-backend/pkg/errors"
	"coursebench-backend/pkg/models"
	"coursebench-backend/pkg/queries"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func CourseGroupComment(c *fiber.Ctx) (err error) {
	id_s := c.Params("id", "GG")
	id_i, err := strconv.Atoi(id_s)
	if err != nil {
		return errors.New(errors.InvalidArgument)
	}
	id := uint(id_i)

	uid, err := session.GetUserID(c)
	if err != nil {
		uid = 0
	}

	db := database.GetDB()
	var comments []models.Comment
	err = db.Preload("User").Preload("CourseGroup").Preload("CourseGroup.Course").Preload("CourseGroup.Teachers").
		Where("course_group_id = ?", id).Find(&comments).Error
	if err != nil {
		return errors.Wrap(err, errors.DatabaseError)
	}
	var likeResult []CommentLikeResult
	if uid != 0 {
		db.Raw("SELECT comment_likes.comment_id,comment_likes.is_like from comments, comment_likes where comments.course_group_id = ? and comment_likes.user_id = ? and comment_likes.comment_id = comments.id and comment_likes.deleted_at is NULL and comments.deleted_at is NULL",
			id, uid).Scan(&likeResult)
	}
	var response []CommentResponse
	response = GenerateResponse(comments, uid, likeResult, true, utils.GetIP(c))
	currentUser, err := queries.GetUserByID(nil, uid)
	if err != nil {
		return err
	}
	if !currentUser.IsAdmin && !currentUser.IsCommunityAdmin {
		for i := range response {
			// 设置评论的 Reward 字段为 -1，表示不可见
			response[i].Reward = -1
		}
	}
	return c.Status(fiber.StatusOK).JSON(models.OKResponse{
		Data:  response,
		Error: false,
	})
}
