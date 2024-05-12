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
	"coursebench-backend/internal/utils"
	"coursebench-backend/pkg/database"
	"coursebench-backend/pkg/errors"
	"coursebench-backend/pkg/models"
	"coursebench-backend/pkg/queries"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type CommentResponse struct {
	ID          uint                    `json:"id"`
	Title       string                  `json:"title"`
	Content     string                  `json:"content"`
	CreateTime  int                     `json:"post_time"`
	UpdateTime  int                     `json:"update_time"`
	Semester    int                     `json:"semester"`
	IsAnonymous bool                    `json:"is_anonymous"`
	Like        int                     `json:"like"`
	Dislike     int                     `json:"dislike"`
	LikeStatus  int                     `json:"like_status"`
	Score       []int64                 `json:"score"`
	User        *models.ProfileResponse `json:"user"`
	Course      struct {
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
	IsFold       bool   `json:"is_fold"`
	IsCovered    bool   `json:"is_covered"`
	CoverTitle   string `json:"cover_title"`
	CoverContent string `json:"cover_content"`
	CoverReason  string `json:"cover_reason"`
}

type CommentLikeResult struct {
	CommentID uint
	IsLike    bool
}

// GenerateResponse generates the response for comments.
// comments: the comments to generate response.
// uid: the user id of the user who is requesting the comments. If the user is not logged in, set it to 0.
// likeResult: the user's like status of the comments.
// showAnonymous: whether to show the anonymous comments. For /comment/user, it should be false.
func GenerateResponse(comments []models.Comment, uid uint, likeResult []CommentLikeResult, showAnonymous bool, ip []string) (response []CommentResponse) {
	mp := make(map[uint]bool)
	for _, v := range likeResult {
		mp[v.CommentID] = v.IsLike
	}
	for _, v := range comments {
		likeStatus := 0
		if like, ok := mp[v.ID]; ok {
			if like {
				likeStatus = 1
			} else {
				likeStatus = 2
			}
		}
		anonymous := v.IsAnonymous
		c := CommentResponse{
			ID:          v.ID,
			Title:       v.Title,
			Content:     v.Content,
			CreateTime:  v.CreateTime,
			UpdateTime:  v.UpdateTime,
			Semester:    v.Semester,
			IsAnonymous: anonymous,
			Like:        v.Like,
			Dislike:     v.Dislike,
			LikeStatus:  likeStatus,
			Score:       v.Scores,
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
			IsFold:       v.IsFold,
			IsCovered:    v.IsCovered,
			CoverTitle:   v.CoverTitle,
			CoverContent: v.CoverContent,
			CoverReason:  v.CoverReason,
		}
		// 该评论未设置匿名，或者是自己的评论，则显示用户信息
		if !anonymous || v.User.ID == uid {
			t, _ := queries.GetProfile(nil, v.UserID, uid)
			c.User = &t
		} else if !showAnonymous {
			continue
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

func UserComment(c *fiber.Ctx) (err error) {
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
		Where("user_id = ?", id).Find(&comments).Error
	if err != nil {
		return errors.Wrap(err, errors.DatabaseError)
	}
	var likeResult []CommentLikeResult
	if uid != 0 {
		db.Raw("SELECT comment_likes.comment_id,comment_likes.is_like from comments, comment_likes where comments.user_id = ? and comment_likes.user_id = ? and comment_likes.comment_id = comments.id and comment_likes.deleted_at is NULL and comments.deleted_at is NULL",
			id, uid).Scan(&likeResult)
	}
	var response []CommentResponse
	response = GenerateResponse(comments, uid, likeResult, false, utils.GetIP(c))
	return c.Status(fiber.StatusOK).JSON(models.OKResponse{
		Data:  response,
		Error: false,
	})
}
