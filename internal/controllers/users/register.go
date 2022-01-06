package users

import (
	"coursebench-backend/pkg/models"
	"coursebench-backend/pkg/queries"
	"github.com/gofiber/fiber/v2"
)

type RegisterRequest struct {
	Email    string           `json:"email"`
	Password string           `json:"password"`
	Year     int              `json:"year"`
	Grade    models.GradeType `json:"grade"`
}

func Register(c *fiber.Ctx) (err error) {
	c.Accepts("application/json")
	var userReq RegisterRequest
	if err = c.BodyParser(&userReq); err != nil {
		return
	}
	user := models.User{
		Email:    userReq.Email,
		Password: userReq.Password,
		Year:     userReq.Year,
		Grade:    userReq.Grade,
	}
	if err = queries.Register(&user); err != nil {
		return
	}

	return c.Status(fiber.StatusOK).JSON(models.OKResponse{
		Data:  map[string]interface{}{"UserID": user.ID},
		Error: false,
	})
}
