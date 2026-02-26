package comments

import (
	"coursebench-backend/internal/middlewares/session"
	"coursebench-backend/pkg/database"
	"coursebench-backend/pkg/errors"
	"coursebench-backend/pkg/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func ReplyChain(c *fiber.Ctx) (err error) {
	replyIDString := c.Params("id", "0")
	replyIDInt, err := strconv.Atoi(replyIDString)
	if err != nil || replyIDInt <= 0 {
		return errors.New(errors.InvalidArgument)
	}
	replyID := uint(replyIDInt)

	uid, err := session.GetUserID(c)
	if err != nil {
		uid = 0
	}

	db := database.GetDB()
	currentReply := &models.Reply{}
	err = db.Preload("User").Preload("ParentReply").Preload("ParentReply.User").Where("id = ?", replyID).Take(currentReply).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New(errors.InvalidArgument)
		}
		return errors.Wrap(err, errors.DatabaseError)
	}

	ancestors, err := buildAncestors(db, *currentReply)
	if err != nil {
		return err
	}

	currentResponses, err := buildReplyResponses(db, []models.Reply{*currentReply}, uid)
	if err != nil {
		return err
	}
	currentResponse := currentResponses[0]

	ancestorResponses, err := buildReplyResponses(db, ancestors, uid)
	if err != nil {
		return err
	}

	descendants, err := buildDescendantTree(db, currentReply.ID, uid)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(models.OKResponse{
		Data: ReplyChainResponse{
			Ancestors:   ancestorResponses,
			Current:     currentResponse,
			Descendants: descendants,
		},
		Error: false,
	})
}

func buildAncestors(db *gorm.DB, reply models.Reply) ([]models.Reply, error) {
	ancestors := make([]models.Reply, 0)
	current := reply

	for current.ParentReplyID != nil {
		parentReply := models.Reply{}
		err := db.Preload("User").Preload("ParentReply").Preload("ParentReply.User").Where("id = ?", *current.ParentReplyID).Take(&parentReply).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New(errors.InvalidArgument)
			}
			return nil, errors.Wrap(err, errors.DatabaseError)
		}
		ancestors = append([]models.Reply{parentReply}, ancestors...)
		current = parentReply
	}

	return ancestors, nil
}

func buildDescendantTree(db *gorm.DB, parentReplyID uint, queryingUserID uint) ([]ReplyTreeNode, error) {
	var children []models.Reply
	err := db.Preload("User").Preload("ParentReply").Preload("ParentReply.User").
		Where("parent_reply_id = ?", parentReplyID).
		Order("create_time ASC").
		Find(&children).Error
	if err != nil {
		return nil, errors.Wrap(err, errors.DatabaseError)
	}
	if len(children) == 0 {
		return []ReplyTreeNode{}, nil
	}

	childResponses, err := buildReplyResponses(db, children, queryingUserID)
	if err != nil {
		return nil, err
	}

	nodes := make([]ReplyTreeNode, 0, len(children))
	for idx, child := range children {
		descendants, err := buildDescendantTree(db, child.ID, queryingUserID)
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, ReplyTreeNode{
			Reply:    childResponses[idx],
			Children: descendants,
		})
	}

	return nodes, nil
}
