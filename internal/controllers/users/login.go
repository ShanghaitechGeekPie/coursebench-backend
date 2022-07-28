package users

import (
	"coursebench-backend/internal/config"
	"coursebench-backend/internal/middlewares/session"
	"coursebench-backend/pkg/errors"
	"coursebench-backend/pkg/models"
	"coursebench-backend/pkg/queries"
	"github.com/gofiber/fiber/v2"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Captcha  string `json:"captcha"`
}

type LoginResponse struct {
	UserID   uint             `json:"user_id"`
	Email    string           `json:"email"`
	Year     int              `json:"year"`
	Grade    models.GradeType `json:"grade"`
	NickName string           `json:"nickname"`
	RealName string           `json:"realname"`
}

func Login(c *fiber.Ctx) (err error) {
	c.Accepts("application/json")
	var request LoginRequest
	if err = c.BodyParser(&request); err != nil {
		return errors.Wrap(err, errors.InvalidArgument)
	}
	if !config.GlobalConf.DisableCaptchaAndMail && !queries.VerifyCaptcha(c, request.Captcha) {
		return errors.New(errors.CaptchaMismatch)
	}

	user, err := queries.Login(request.Email, request.Password)
	if err != nil {
		return
	}
	sess, err := session.GetStore().Get(c)
	if err != nil {
		return errors.Wrap(err, errors.InternalServerError)
	}
	sess.Set("user_id", user.ID)
	// Save session
	if err := sess.Save(); err != nil {
		return errors.Wrap(err, errors.InternalServerError)
	}
	/*err = session.SetSession(user.ID, sess)
	if err != nil {
		return
	}*/

	response := LoginResponse{
		UserID:   user.ID,
		Email:    user.Email,
		Year:     user.Year,
		Grade:    user.Grade,
		NickName: user.NickName,
		RealName: user.RealName,
	}

	return c.Status(fiber.StatusOK).JSON(models.OKResponse{
		Data:  response,
		Error: false,
	})
}
