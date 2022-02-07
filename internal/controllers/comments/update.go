package comments

import (
	"coursebench-backend/internal/middlewares/session"
	"coursebench-backend/pkg/database"
	"coursebench-backend/pkg/errors"
	"coursebench-backend/pkg/models"
	"coursebench-backend/pkg/queries"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"time"
)

type UpdateRequest struct {
	ID                  int     `json:"id"`
	Title               string  `json:"title"`
	Content             string  `json:"content"`
	Semester            int     `json:"semester"`
	IsAnonymous         bool    `json:"is_anonymous"`
	Scores              []int64 `json:"scores"`
	StudentScoreRanking int     `json:"student_score_ranking"`
}

func Update(c *fiber.Ctx) (err error) {
	updateTime := int(time.Now().Unix())

	c.Accepts("application/json")
	var request UpdateRequest
	if err = c.BodyParser(&request); err != nil {
		return errors.Wrap(err, errors.InvalidArgument)
	}

	uid, err := session.GetUserID(c)
	if err != nil {
		return err
	}

	if !queries.CheckCommentTitle(request.Title) {
		return errors.New(errors.InvalidArgument)
	}
	if !queries.CheckCommentContent(request.Content) {
		return errors.New(errors.InvalidArgument)
	}
	if !queries.CheckSemester(request.Semester) {
		return errors.New(errors.InvalidArgument)
	}
	if !queries.CheckCommentScore(request.Scores) {
		return errors.New(errors.InvalidArgument)
	}
	if !queries.CheckCommentScoreRanking(request.StudentScoreRanking) {
		return errors.New(errors.InvalidArgument)
	}

	db := database.GetDB()
	comment := &models.Comment{}
	err = db.Preload("CourseGroup").Preload("CourseGroup.Course").Where("id = ?", request.ID).Take(comment).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New(errors.CommentNotExists)
		} else {
			return errors.Wrap(err, errors.DatabaseError)
		}
	}

	if uid != comment.UserID {
		return errors.New(errors.PermissionDenied)
	}

	err = db.Transaction(func(tx *gorm.DB) error {

		for i := 0; i < models.ScoreLength; i++ {
			comment.CourseGroup.Scores[i] += request.Scores[i] - comment.Scores[i]
		}
		for i := 0; i < models.ScoreLength; i++ {
			comment.CourseGroup.Course.Scores[i] += request.Scores[i] - comment.Scores[i]
		}

		comment.Title = request.Title
		comment.Content = request.Content
		comment.Semester = request.Semester
		comment.IsAnonymous = request.IsAnonymous
		comment.UpdateTime = updateTime
		comment.StudentScoreRanking = request.StudentScoreRanking
		comment.Scores = request.Scores
		err = tx.Session(&gorm.Session{FullSaveAssociations: true}).Select("Title", "Content", "Semester", "IsAnonymous", "UpdateTime",
			"StudentScoreRanking", "Scores", "CourseGroup", "CourseGroup.Scores", "CourseGroup.ID",
			"CourseGroup.Course", "CourseGroup.Course.ID", "CourseGroup.Course.Scores").Updates(comment).Error
		//err = tx.Session(&gorm.Session{FullSaveAssociations: true}).Updates(comment).Error
		if err != nil {
			return errors.Wrap(err, errors.DatabaseError)
		}

		return nil
	})

	return c.Status(fiber.StatusOK).JSON(models.OKResponse{
		Data:  nil,
		Error: false,
	})
}
