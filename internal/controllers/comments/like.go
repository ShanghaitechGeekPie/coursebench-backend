package comments

import (
	"coursebench-backend/internal/middlewares/session"
	"coursebench-backend/pkg/database"
	"coursebench-backend/pkg/errors"
	"coursebench-backend/pkg/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type LikeRequest struct {
	ID     int `json:"id"`
	Status int `json:"status"`
}

func Like(c *fiber.Ctx) (err error) {
	c.Accepts("application/json")
	var request LikeRequest
	if err = c.BodyParser(&request); err != nil {
		return errors.Wrap(err, errors.InvalidArgument)
	}
	if request.Status < 0 || request.Status > 2 {
		return errors.New(errors.InvalidArgument)
	}

	uid, err := session.GetUserID(c)
	if err != nil {
		return err
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

	if uid == comment.UserID {
		return errors.New(errors.PermissionDenied)
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		cl := &models.CommentLike{}
		err = tx.Where("user_id = ? AND comment_id = ?", uid, request.ID).Take(cl).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.Wrap(err, errors.DatabaseError)
		}
		if err != nil {
			// No record found, create a new one
			if request.Status != 0 {
				cl = &models.CommentLike{
					UserID:    uid,
					CommentID: uint(request.ID),
					IsLike:    request.Status == 1,
				}
				err = tx.Create(cl).Error
				if err != nil {
					return errors.Wrap(err, errors.DatabaseError)
				}
				if request.Status == 1 {
					err = tx.Model(comment).Where("id = ?", comment.ID).Update("like", comment.Like+1).Error
				} else {
					err = tx.Model(comment).Where("id = ?", comment.ID).Update("dislike", comment.Dislike+1).Error
				}
				if err != nil {
					return errors.Wrap(err, errors.DatabaseError)
				}
			}
		} else {
			if request.Status == 0 {
				err = tx.Delete(cl).Error
				if err != nil {
					return errors.Wrap(err, errors.DatabaseError)
				}
				if cl.IsLike {
					err = tx.Model(comment).Where("id = ?", comment.ID).Update("like", comment.Like-1).Error
				} else {
					err = tx.Model(comment).Where("id = ?", comment.ID).Update("dislike", comment.Dislike-1).Error
				}
				if err != nil {
					return errors.Wrap(err, errors.DatabaseError)
				}
			} else {
				if cl.IsLike != (request.Status == 1) {
					var k int
					if request.Status == 1 {
						k = 1
					} else {
						k = -1
					}
					err = tx.Model(comment).Where("id = ?", comment.ID).Update("like", comment.Like+k).Error
					if err != nil {
						return errors.Wrap(err, errors.DatabaseError)
					}
					err = tx.Model(comment).Where("id = ?", comment.ID).Update("dislike", comment.Dislike-k).Error
					if err != nil {
						return errors.Wrap(err, errors.DatabaseError)
					}
					cl.IsLike = request.Status == 1
					err = tx.Save(cl).Error
					if err != nil {
						return errors.Wrap(err, errors.DatabaseError)
					}
				}
			}
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
