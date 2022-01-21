package users

import (
	"coursebench-backend/internal/middlewares/session"
	"coursebench-backend/pkg/errors"
	"coursebench-backend/pkg/models"
	"coursebench-backend/pkg/queries"
	"github.com/gofiber/fiber/v2"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Login(c *fiber.Ctx) (err error) {
	c.Accepts("application/json")
	var request LoginRequest
	if err = c.BodyParser(&request); err != nil {
		return errors.Wrap(err, errors.InvalidArgument)
	}

	user, err := queries.Login(request.Email, request.Password)
	if err != nil {
		return
	}
	sess, err := session.GetStore().Get(c)
	if err != nil {
		return errors.Wrap(err, errors.InternalServerError)
	}
	// Save session
	if err := sess.Save(); err != nil {
		return errors.Wrap(err, errors.InternalServerError)
	}
	sess.Set("user_id", user.ID)
	/*err = session.SetSession(user.ID, sess)
	if err != nil {
		return
	}*/

	return c.Status(fiber.StatusOK).JSON(models.OKResponse{
		Data:  map[string]interface{}{"UserID": user.ID},
		Error: false,
	})
}
