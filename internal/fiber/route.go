package fiber

import (
	"coursebench-backend/internal/config"
	"coursebench-backend/internal/controllers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func Routes(app *fiber.App) {
	route := app.Group("/v1")
	route.Use(cors.New())
	controllers.UserRoutes(route)
	if config.GlobalConf.InDevelopment {
		controllers.TestRoutes(route)
	}
	controllers.CourseRoutes(route)
	controllers.CommentRoutes(route)
	controllers.TeacherRoute(route)

}
