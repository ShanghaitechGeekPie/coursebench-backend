package users

import (
	"coursebench-backend/internal/config"
	"coursebench-backend/internal/middlewares/session"
	"coursebench-backend/pkg/database"
	"coursebench-backend/pkg/errors"
	"coursebench-backend/pkg/models"
	"coursebench-backend/pkg/queries"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	syslog "log"
)

func UploadAvatar(c *fiber.Ctx) (err error) {
	id, err := session.GetUserID(c)
	if err != nil {
		return err
	}
	file, err := c.FormFile("avatar")
	if err != nil {
		return errors.Wrap(err, errors.InvalidArgument)
	}
	if config.GlobalConf.AvatarSizeLimit < file.Size {
		return errors.Wrap(err, errors.FileTooLarge)
	}
	r, err := file.Open()
	if err != nil {
		return errors.Wrap(err, errors.InternalServerError)
	}
	user, err := queries.GetUserByID(id)
	if err != nil {
		return err
	}
	oldAvatar := user.Avatar
	nameUUID, err := uuid.NewRandom()
	if err != nil {
		return errors.Wrap(err, errors.InternalServerError)
	}
	name := "avatar/" + nameUUID.String()
	err = database.UploadFile(c.Context(), name, r, file.Size)
	if err != nil {
		return err
	}
	user.Avatar = nameUUID.String()
	db := database.GetDB()
	if err = db.Save(&user).Error; err != nil {
		return errors.Wrap(err, errors.DatabaseError)
	}
	if oldAvatar != "" {
		err = database.DeleteFile(c.Context(), "avatar/"+oldAvatar)
		if err != nil {
			syslog.Println(err)
		}
	}
	return c.Status(fiber.StatusOK).JSON(models.OKResponse{Data: map[string]string{"avatar": nameUUID.String()}, Error: false})
}
