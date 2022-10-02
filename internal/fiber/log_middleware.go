package fiber

import (
	"coursebench-backend/pkg/log"
	"fmt"
	"github.com/gofiber/fiber/v2"
)

func LogMiddleware(c *fiber.Ctx) error {
	s := "New Request: " + c.String()
	if len(c.Body()) < 500 {
		s += "; Body: " + string(c.Body())
	} else {
		s += fmt.Sprintf("; Body: Too long, %v bytes, skip...", len(c.Body()))
	}
	log.Println(s)
	return c.Next()
}
