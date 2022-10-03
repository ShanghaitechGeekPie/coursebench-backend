package logger

import (
	"coursebench-backend/pkg/log"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"strings"
)

// LogMiddleware record the details of each request.
func LogMiddleware(c *fiber.Ctx) error {
	s := "New Request: " + c.String()
	s += " ; X-Forwarded-For: " + strings.Join(c.IPs(), ",")
	if c.Path() == "/v1/user/login" || c.Path() == "/v1/user/register" || c.Path() == "/v1/user/update_password" || c.Path() == "/v1/user/reset_password_active" { // Don't log the password. parsing json is too much trouble here, so I just don't record the whole body.
		s += " ; Body: Sensitive, skip... "
	} else if len(c.Body()) < 500 {
		s += " ; Body: " + string(c.Body())
	} else {
		s += fmt.Sprintf(" ; Body: Too long, %v bytes, skip...", len(c.Body()))
	}
	log.Println(s)
	return c.Next()
}
