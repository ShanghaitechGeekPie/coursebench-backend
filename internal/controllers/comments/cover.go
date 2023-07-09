package comments

import (
	"bytes"
	"coursebench-backend/internal/config"
	"coursebench-backend/internal/middlewares/session"
	"coursebench-backend/pkg/database"
	"coursebench-backend/pkg/errors"
	"coursebench-backend/pkg/models"
	"coursebench-backend/pkg/queries"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"io"
	"net/http"
)

type CoverRequest struct {
	ID     int  `json:"id"`
	Status bool `json:"status"`
}

type GPTWorkerResponse struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Reason  string `json:"reason"`
}

func Cover(c *fiber.Ctx) (err error) {
	c.Accepts("application/json")
	var request FoldRequest
	if err = c.BodyParser(&request); err != nil {
		return errors.Wrap(err, errors.InvalidArgument)
	}

	uid, err := session.GetUserID(c)
	if err != nil {
		return err
	}

	user, err := queries.GetUserByID(nil, uid)
	if err != nil {
		return err
	}
	if !user.IsAdmin && !user.IsCommunityAdmin {
		return errors.New(errors.PermissionDenied)
	}

	db := database.GetDB()
	err = db.Transaction(func(tx *gorm.DB) error {
		comment := &models.Comment{}
		err = tx.Where("id = ?", request.ID).Take(comment).Error

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New(errors.CommentNotExists)
			} else {
				return errors.Wrap(err, errors.DatabaseError)
			}
		}
		if request.Status { // cover
			response, err := http.Post(config.GlobalConf.GPTWorkerURL, "application/json", bytes.NewBuffer([]byte(`{"title": "`+comment.Title+`", "content": "`+comment.Content+`"}`)))
			if err != nil {
				return errors.New(errors.GPTWorkerError)
			}
			defer response.Body.Close()
			var gptWorkerResponse GPTWorkerResponse
			body, err := io.ReadAll(response.Body)
			if err != nil {
				return errors.New(errors.GPTWorkerError)
			}
			err = json.Unmarshal(body, &gptWorkerResponse)
			if err != nil {
				return errors.New(errors.GPTWorkerError)
			}
			comment.CoverTitle = gptWorkerResponse.Title
			comment.CoverContent = gptWorkerResponse.Content
			comment.CoverReason = gptWorkerResponse.Reason
			comment.IsCovered = true
		} else { // uncover
			comment.CoverTitle = ""
			comment.CoverContent = ""
			comment.CoverReason = ""
			comment.IsCovered = false
		}

		err = tx.Select("IsCovered", "CoverTitle", "CoverContent").Save(comment).Error
		if err != nil {
			return errors.Wrap(err, errors.DatabaseError)
		}
		return nil
	})
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(models.OKResponse{
		Data:  map[string]interface{}{},
		Error: false})
}
