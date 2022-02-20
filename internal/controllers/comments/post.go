package comments

import (
	"coursebench-backend/internal/middlewares/session"
	"coursebench-backend/pkg/database"
	"coursebench-backend/pkg/events"
	"coursebench-backend/pkg/models"
	"coursebench-backend/pkg/queries"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"time"
)

type PostRequest struct {
	Group               uint    `json:"group"`
	Title               string  `json:"title"`
	Content             string  `json:"content"`
	Semester            int     `json:"semester"`
	IsAnonymous         bool    `json:"is_anonymous"`
	Scores              []int64 `json:"scores"`
	StudentScoreRanking int     `json:"student_score_ranking"`
}

func Post(c *fiber.Ctx) *events.AttributedEvent {
	postTime := int(time.Now().Unix())

	c.Accepts("application/json")
	var request PostRequest
	if err := c.BodyParser(&request); err != nil {
		return events.Wrap(err, events.InvalidArgument)
	}

	uid, event := session.GetUserID(c)
	if event != nil {
		return event
	}

	if !queries.CheckCommentTitle(request.Title) {
		return events.New(events.InvalidArgument)
	}
	if !queries.CheckCommentContent(request.Content) {
		return events.New(events.InvalidArgument)
	}
	if !queries.CheckSemester(request.Semester) {
		return events.New(events.InvalidArgument)
	}
	if !queries.CheckCommentScore(request.Scores) {
		return events.New(events.InvalidArgument)
	}
	if !queries.CheckCommentScoreRanking(request.StudentScoreRanking) {
		return events.New(events.InvalidArgument)
	}

	db := database.GetDB()

	// Check if comment already exists
	err := db.Where("user_id=? AND course_group_id=?", uid, request.Group).Take(&models.Comment{}).Error
	if err == nil {
		return events.New(events.CommentAlreadyExists)
	} else if err != gorm.ErrRecordNotFound {
		return events.Wrap(err, events.DatabaseError)
	}
	err = db.Where("id=?", request.Group).Take(&models.CourseGroup{}).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return events.Wrap(err, events.InvalidArgument)
		} else {
			return events.Wrap(err, events.DatabaseError)
		}
	}
	comment := &models.Comment{
		UserID:              uid,
		CourseGroupID:       request.Group,
		Title:               request.Title,
		Content:             request.Content,
		Semester:            request.Semester,
		IsAnonymous:         request.IsAnonymous,
		StudentScoreRanking: request.StudentScoreRanking,
		CreateTime:          postTime,
		UpdateTime:          postTime,
		Scores:              request.Scores,
	}

	// 插入评论作为一个事务，若插入失败则回滚
	err = db.Transaction(func(tx *gorm.DB) error {
		err = tx.Create(comment).Error
		if err != nil {
			return events.Wrap(err, events.DatabaseError).ToError()

		}

		group := &models.CourseGroup{}
		err = tx.Preload("Course").Where("id=?", request.Group).Take(group).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return events.New(events.CourseGroupNotExists).ToError()
		} else if err != nil {
			return events.Wrap(err, events.DatabaseError).ToError()
		}
		group.CommentCount++
		for i := 0; i < models.ScoreLength; i++ {
			group.Scores[i] += request.Scores[i]
		}
		err = tx.Select("Scores", "CommentCount").Save(group).Error
		if err != nil {
			return events.Wrap(err, events.DatabaseError).ToError()
		}
		group.Course.CommentCount++
		for i := 0; i < models.ScoreLength; i++ {
			group.Course.Scores[i] += request.Scores[i]
		}
		err = tx.Select("Scores", "CommentCount").Save(group.Course).Error
		if err != nil {
			return events.Wrap(err, events.DatabaseError).ToError()
		}
		return nil
	})
	if err != nil {
		return events.Wrap(err, events.DatabaseError)
	}

	return events.Wrap(c.Status(fiber.StatusOK).JSON(models.OKResponse{
		Data:  map[string]interface{}{"comment_id": comment.ID},
		Error: false,
	}), events.InternalServerError)
}
