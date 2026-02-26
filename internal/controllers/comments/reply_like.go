package comments

import (
	"coursebench-backend/internal/middlewares/session"
	"coursebench-backend/pkg/database"
	"coursebench-backend/pkg/errors"
	"coursebench-backend/pkg/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ReplyLikeRequest struct {
	ID     uint `json:"id"`
	Status int  `json:"status"`
}

func ReplyLike(c *fiber.Ctx) (err error) {
	c.Accepts("application/json")
	var request ReplyLikeRequest
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
	reply := &models.Reply{}
	err = db.Where("id = ?", request.ID).Take(reply).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New(errors.InvalidArgument)
		}
		return errors.Wrap(err, errors.DatabaseError)
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		rl := &models.ReplyLike{}
		err = tx.Where("user_id = ? AND reply_id = ?", uid, request.ID).Take(rl).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.Wrap(err, errors.DatabaseError)
		}

		if err != nil {
			if request.Status == 0 {
				return nil
			}
			rl = &models.ReplyLike{UserID: uid, ReplyID: request.ID, IsLike: request.Status == 1}
			err = tx.Create(rl).Error
			if err != nil {
				return errors.Wrap(err, errors.DatabaseError)
			}
			if request.Status == 1 {
				err = tx.Model(reply).Where("id = ?", reply.ID).Update("\"like\"", reply.Like+1).Error
			} else {
				err = tx.Model(reply).Where("id = ?", reply.ID).Update("dislike", reply.Dislike+1).Error
			}
			if err != nil {
				return errors.Wrap(err, errors.DatabaseError)
			}
			return nil
		}

		if request.Status == 0 {
			err = tx.Delete(rl).Error
			if err != nil {
				return errors.Wrap(err, errors.DatabaseError)
			}
			if rl.IsLike {
				err = tx.Model(reply).Where("id = ?", reply.ID).Update("\"like\"", reply.Like-1).Error
			} else {
				err = tx.Model(reply).Where("id = ?", reply.ID).Update("dislike", reply.Dislike-1).Error
			}
			if err != nil {
				return errors.Wrap(err, errors.DatabaseError)
			}
			return nil
		}

		if rl.IsLike == (request.Status == 1) {
			return nil
		}

		if request.Status == 1 {
			err = tx.Model(reply).Where("id = ?", reply.ID).Update("\"like\"", reply.Like+1).Error
			if err != nil {
				return errors.Wrap(err, errors.DatabaseError)
			}
			err = tx.Model(reply).Where("id = ?", reply.ID).Update("dislike", reply.Dislike-1).Error
		} else {
			err = tx.Model(reply).Where("id = ?", reply.ID).Update("\"like\"", reply.Like-1).Error
			if err != nil {
				return errors.Wrap(err, errors.DatabaseError)
			}
			err = tx.Model(reply).Where("id = ?", reply.ID).Update("dislike", reply.Dislike+1).Error
		}
		if err != nil {
			return errors.Wrap(err, errors.DatabaseError)
		}

		rl.IsLike = request.Status == 1
		err = tx.Save(rl).Error
		if err != nil {
			return errors.Wrap(err, errors.DatabaseError)
		}
		return nil
	})
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(models.OKResponse{Data: map[string]interface{}{}, Error: false})
}
