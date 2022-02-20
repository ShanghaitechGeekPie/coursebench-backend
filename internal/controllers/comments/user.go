package comments

import (
	"coursebench-backend/internal/middlewares/session"
	"coursebench-backend/pkg/database"
	"coursebench-backend/pkg/events"
	"coursebench-backend/pkg/models"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type CommentResponse struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	CreateTime  int    `json:"post_time"`
	UpdateTime  int    `json:"update_time"`
	Semester    int    `json:"semester"`
	IsAnonymous bool   `json:"is_anonymous"`
	User        *struct {
		ID       uint   `json:"id"`
		Nickname string `json:"nickname"`
	} `json:"user"`
	Course struct {
		ID        uint   `json:"id"`
		Name      string `json:"name"`
		Code      string `json:"code"`
		Institute string `json:"institute"`
	} `json:"course"`
	Group struct {
		ID       uint   `json:"id"`
		Code     string `json:"code"`
		Teachers []struct {
			ID   uint   `json:"id"`
			Name string `json:"name"`
		} `json:"teachers"`
	} `json:"group"`
}

func GenerateResponse(comments []models.Comment, uid uint) (response []CommentResponse) {
	for _, v := range comments {
		c := CommentResponse{
			ID:          v.ID,
			Title:       v.Title,
			Content:     v.Content,
			CreateTime:  v.CreateTime,
			UpdateTime:  v.UpdateTime,
			Semester:    v.Semester,
			IsAnonymous: v.IsAnonymous,
			Course: struct {
				ID        uint   `json:"id"`
				Name      string `json:"name"`
				Code      string `json:"code"`
				Institute string `json:"institute"`
			}{
				ID:        v.CourseGroup.Course.ID,
				Name:      v.CourseGroup.Course.Name,
				Code:      v.CourseGroup.Course.Code,
				Institute: v.CourseGroup.Course.Institute,
			},
			Group: struct {
				ID       uint   `json:"id"`
				Code     string `json:"code"`
				Teachers []struct {
					ID   uint   `json:"id"`
					Name string `json:"name"`
				} `json:"teachers"`
			}{
				ID:   v.CourseGroup.ID,
				Code: v.CourseGroup.Code,
			},
		}
		// 该评论未设置匿名，或者是自己的评论，则显示用户信息
		if !c.IsAnonymous || v.User.ID == uid {
			c.User = &struct {
				ID       uint   `json:"id"`
				Nickname string `json:"nickname"`
			}{
				ID:       v.User.ID,
				Nickname: v.User.NickName,
			}
		}
		for _, t := range v.CourseGroup.Teachers {
			c.Group.Teachers = append(c.Group.Teachers, struct {
				ID   uint   `json:"id"`
				Name string `json:"name"`
			}{
				ID:   t.ID,
				Name: t.Name,
			})
		}
		response = append(response, c)
	}
	return
}

func UserComment(c *fiber.Ctx) *events.AttributedEvent {
	id_s := c.Params("id", "GG")
	id_i, err := strconv.Atoi(id_s)
	if err != nil {
		return events.New(events.InvalidArgument)
	}
	id := uint(id_i)

	uid, event := session.GetUserID(c)
	if event != nil {
		return event
	}

	db := database.GetDB()
	var comments []models.Comment
	err = db.Preload("User").Preload("CourseGroup").Preload("CourseGroup.Course").Preload("CourseGroup.Teachers").
		Where("user_id = ?", id).Find(&comments).Error
	if err != nil {
		return events.Wrap(err, events.DatabaseError)
	}
	var response []CommentResponse
	response = GenerateResponse(comments, uid)
	return events.Wrap(c.Status(fiber.StatusOK).JSON(models.OKResponse{
		Data:  response,
		Error: false,
	}), events.InternalServerError)
}
