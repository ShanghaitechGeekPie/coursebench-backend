package main

import (
	"coursebench-backend/internal/fiber"
	"coursebench-backend/pkg/database"
	"coursebench-backend/pkg/modelRegister"
	_ "coursebench-backend/pkg/models"
	_ "github.com/joho/godotenv/autoload"
	"os"
)

func main() {
	database.InitDB()
	database.InitRedis()
	db := database.GetDB()
	err := db.Migrator().AutoMigrate(modelRegister.GetRegisteredTypes()...)
	if err != nil {
		panic(err)
	}
	app := fiber.New()
	fiber.Routes(app)
	if err := app.Listen(os.Getenv("SERVER_URL")); err != nil {
		panic(err)
	}
}
