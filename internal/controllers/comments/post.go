// Copyright (C) 2021-2024 ShanghaiTech GeekPie
// This file is part of CourseBench Backend.
//
// CourseBench Backend is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// CourseBench Backend is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with CourseBench Backend.  If not, see <http://www.gnu.org/licenses/>.

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

type PostRequest struct {
	Group               uint    `json:"group"`
	Title               string  `json:"title"`
	Content             string  `json:"content"`
	Semester            int     `json:"semester"`
	IsAnonymous         bool    `json:"is_anonymous"`
	Scores              []int64 `json:"scores"`
	StudentScoreRanking int     `json:"student_score_ranking"`
}

func Post(c *fiber.Ctx) (err error) {
	postTime := int(time.Now().Unix())

	c.Accepts("application/json")
	var request PostRequest
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

	// Check if comment already exists
	err = db.Where("user_id=? AND course_group_id=?", uid, request.Group).Take(&models.Comment{}).Error
	if err == nil {
		return errors.New(errors.CommentAlreadyExists)
	} else if err != gorm.ErrRecordNotFound {
		return errors.Wrap(err, errors.DatabaseError)
	}
	group := models.CourseGroup{}
	err = db.Where("id=?", request.Group).Take(&group).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.Wrap(err, errors.InvalidArgument)
		} else {
			return errors.Wrap(err, errors.DatabaseError)
		}
	}
	comment := &models.Comment{
		UserID:              uid,
		CourseGroupID:       request.Group,
		CourseID:            group.CourseID,
		Title:               request.Title,
		Content:             request.Content,
		Semester:            request.Semester,
		IsAnonymous:         request.IsAnonymous,
		StudentScoreRanking: request.StudentScoreRanking,
		CreateTime:          postTime,
		UpdateTime:          postTime,
		Scores:              request.Scores,
		Like:                0,
		Dislike:             0,
		IsFold:              false,
	}

	// 插入评论作为一个事务，若插入失败则回滚
	err = db.Transaction(func(tx *gorm.DB) error {
		err = tx.Create(comment).Error
		if err != nil {
			return errors.Wrap(err, errors.DatabaseError)

		}

		group := &models.CourseGroup{}
		err = tx.Preload("Course").Where("id=?", request.Group).Take(group).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New(errors.CourseGroupNotExists)
		} else if err != nil {
			return errors.Wrap(err, errors.DatabaseError)
		}
		group.CommentCount++
		for i := 0; i < models.ScoreLength; i++ {
			group.Scores[i] += request.Scores[i]
		}
		err = tx.Select("Scores", "CommentCount").Save(group).Error
		if err != nil {
			return errors.Wrap(err, errors.DatabaseError)
		}
		group.Course.CommentCount++
		for i := 0; i < models.ScoreLength; i++ {
			group.Course.Scores[i] += request.Scores[i]
		}
		err = tx.Select("Scores", "CommentCount").Save(group.Course).Error
		if err != nil {
			return errors.Wrap(err, errors.DatabaseError)
		}
		return nil
	})
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(models.OKResponse{
		Data:  map[string]interface{}{"comment_id": comment.ID},
		Error: false,
	})
}
