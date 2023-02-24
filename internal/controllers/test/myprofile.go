package test

import (
	"coursebench-backend/internal/middlewares/session"
	"coursebench-backend/pkg/models"
	"coursebench-backend/pkg/queries"
	"fmt"
	"github.com/gofiber/fiber/v2"
)

func MyProfile(c *fiber.Ctx) (err error) {
	fmt.Println(c.Cookies("session_id", "GG"))
	id, err := session.GetUserID(c)
	if err != nil {
		return err
	}
	user, err := queries.GetUserByID(nil, id)
	return c.Status(fiber.StatusOK).JSON(models.OKResponse{
		Data:  map[string]interface{}{"UserID": id, "User": user},
		Error: false,
	})
}
