package users

import (
	"coursebench-backend/pkg/models"
	"github.com/gofiber/fiber/v2"
)

type RegisterUser struct {
	Email    string           `json:"email"`
	Password string           `json:"password"`
	Year     int              `json:"year"`
	Grade    models.GradeType `json:"grade"`
}

func Register(c *fiber.Ctx) (err error) {
	var userReq RegisterUser
	println("!")
	if err = c.BodyParser(&userReq); err != nil {
		return
	}
	println("?")
	user := models.User{
		Email:    userReq.Email,
		Password: userReq.Password,
		Year:     userReq.Year,
		Grade:    userReq.Grade,
	}
	if err = user.Register(); err != nil {
		return
	}

	return c.Status(fiber.StatusOK).JSON(models.OKResponse{
		Data:  nil,
		Error: false,
	})
}
