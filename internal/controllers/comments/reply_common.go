package comments

import (
	"coursebench-backend/pkg/errors"
	"coursebench-backend/pkg/models"
	"coursebench-backend/pkg/queries"

	"gorm.io/gorm"
)

type ReplyUserResponse struct {
	ID          uint   `json:"id"`
	NickName    string `json:"nickname"`
	Avatar      string `json:"avatar"`
	IsAnonymous bool   `json:"is_anonymous"`
}

type ReplyTargetResponse struct {
	ReplyID uint              `json:"reply_id"`
	User    ReplyUserResponse `json:"user"`
}

type ReplyResponse struct {
	ID            uint                 `json:"id"`
	CommentID     uint                 `json:"comment_id"`
	ParentReplyID *uint                `json:"parent_reply_id"`
	Content       string               `json:"content"`
	CreateTime    int                  `json:"post_time"`
	UpdateTime    int                  `json:"update_time"`
	Like          int                  `json:"like"`
	Dislike       int                  `json:"dislike"`
	LikeStatus    int                  `json:"like_status"`
	IsAnonymous   bool                 `json:"is_anonymous"`
	HasSubReplies bool                 `json:"has_sub_replies"`
	User          ReplyUserResponse    `json:"user"`
	ReplyTo       *ReplyTargetResponse `json:"reply_to,omitempty"`
}

type ReplyTreeNode struct {
	Reply    ReplyResponse   `json:"reply"`
	Children []ReplyTreeNode `json:"children"`
}

type ReplyListResponse struct {
	TotalCount    int64           `json:"total_count"`
	FilteredCount int64           `json:"filtered_count"`
	Replies       []ReplyResponse `json:"replies"`
}

type ReplyChainResponse struct {
	Ancestors   []ReplyResponse `json:"ancestors"`
	Current     ReplyResponse   `json:"current"`
	Descendants []ReplyTreeNode `json:"descendants"`
}

func buildReplyUserResponse(db *gorm.DB, userID uint, isAnonymous bool, queryingUserID uint) (ReplyUserResponse, error) {
	if isAnonymous && userID != queryingUserID {
		return ReplyUserResponse{
			ID:          0,
			NickName:    "匿名用户",
			Avatar:      "",
			IsAnonymous: true,
		}, nil
	}

	profile, err := queries.GetProfile(db, userID, queryingUserID)
	if err != nil {
		return ReplyUserResponse{}, err
	}

	return ReplyUserResponse{
		ID:          profile.ID,
		NickName:    profile.NickName,
		Avatar:      profile.Avatar,
		IsAnonymous: isAnonymous,
	}, nil
}

func buildReplyLikeStatusMap(replies []models.Reply, likes []models.ReplyLike) map[uint]int {
	likeStatusMap := make(map[uint]int, len(replies))
	for _, reply := range replies {
		likeStatusMap[reply.ID] = 0
	}
	for _, like := range likes {
		if like.IsLike {
			likeStatusMap[like.ReplyID] = 1
		} else {
			likeStatusMap[like.ReplyID] = 2
		}
	}
	return likeStatusMap
}

func buildHasSubRepliesMap(db *gorm.DB, replies []models.Reply) (map[uint]bool, error) {
	hasSubRepliesMap := make(map[uint]bool, len(replies))
	if len(replies) == 0 {
		return hasSubRepliesMap, nil
	}

	replyIDs := make([]uint, 0, len(replies))
	for _, reply := range replies {
		replyIDs = append(replyIDs, reply.ID)
		hasSubRepliesMap[reply.ID] = false
	}

	type parentReplyIDRow struct {
		ParentReplyID uint `json:"parent_reply_id"`
	}
	var rows []parentReplyIDRow
	err := db.Model(&models.Reply{}).
		Select("parent_reply_id").
		Where("parent_reply_id IN ?", replyIDs).
		Group("parent_reply_id").
		Scan(&rows).Error
	if err != nil {
		return nil, errors.Wrap(err, errors.DatabaseError)
	}

	for _, row := range rows {
		hasSubRepliesMap[row.ParentReplyID] = true
	}
	return hasSubRepliesMap, nil
}

func buildReplyResponse(db *gorm.DB, reply models.Reply, queryingUserID uint, likeStatus int, hasSubReplies bool) (ReplyResponse, error) {
	user, err := buildReplyUserResponse(db, reply.UserID, reply.IsAnonymous, queryingUserID)
	if err != nil {
		return ReplyResponse{}, err
	}

	res := ReplyResponse{
		ID:            reply.ID,
		CommentID:     reply.CommentID,
		ParentReplyID: reply.ParentReplyID,
		Content:       reply.Content,
		CreateTime:    reply.CreateTime,
		UpdateTime:    reply.UpdateTime,
		Like:          reply.Like,
		Dislike:       reply.Dislike,
		LikeStatus:    likeStatus,
		IsAnonymous:   reply.IsAnonymous,
		HasSubReplies: hasSubReplies,
		User:          user,
	}

	if reply.ParentReply != nil {
		replyToUser, err := buildReplyUserResponse(db, reply.ParentReply.UserID, reply.ParentReply.IsAnonymous, queryingUserID)
		if err != nil {
			return ReplyResponse{}, err
		}
		res.ReplyTo = &ReplyTargetResponse{
			ReplyID: reply.ParentReply.ID,
			User:    replyToUser,
		}
	}

	return res, nil
}

func buildReplyResponses(db *gorm.DB, replies []models.Reply, queryingUserID uint) ([]ReplyResponse, error) {
	responses := make([]ReplyResponse, 0, len(replies))
	if len(replies) == 0 {
		return responses, nil
	}

	replyIDs := make([]uint, 0, len(replies))
	for _, reply := range replies {
		replyIDs = append(replyIDs, reply.ID)
	}

	var likes []models.ReplyLike
	if queryingUserID != 0 {
		err := db.Where("user_id = ? AND reply_id IN ?", queryingUserID, replyIDs).Find(&likes).Error
		if err != nil {
			return nil, errors.Wrap(err, errors.DatabaseError)
		}
	}

	likeStatusMap := buildReplyLikeStatusMap(replies, likes)
	hasSubRepliesMap, err := buildHasSubRepliesMap(db, replies)
	if err != nil {
		return nil, err
	}

	for _, reply := range replies {
		response, err := buildReplyResponse(db, reply, queryingUserID, likeStatusMap[reply.ID], hasSubRepliesMap[reply.ID])
		if err != nil {
			return nil, err
		}
		responses = append(responses, response)
	}

	return responses, nil
}

func getReplyOrderBy(sortBy string) (string, error) {
	switch sortBy {
	case "", "latest":
		return "create_time DESC", nil
	case "hottest":
		return "\"like\" DESC, create_time DESC", nil
	default:
		return "", errors.New(errors.InvalidArgument)
	}
}
