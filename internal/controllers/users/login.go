package users

import (
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
		return
	}

	user, err := queries.Login(request.Email, request.Password)
	if err != nil {
		return
	}
	/*
		if errors.Is(err, errors.UserDoNotExist) {
			return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
				Error:   true,
				Errno:   errors.UserDoNotExist.Name(),
				Message: errors.UserDoNotExist.Error()})
		} else if errors.Is(err, errors.UserPasswordIncorrect) {
			return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
				Error:   true,
				Errno:   errors.UserPasswordIncorrect.Name(),
				Message: errors.UserPasswordIncorrect.Error()})
		} else if err != nil {
			return
		}
	*/

	return c.Status(fiber.StatusOK).JSON(models.OKResponse{
		Data:  map[string]interface{}{"UserID": user.ID},
		Error: false,
	})
}
