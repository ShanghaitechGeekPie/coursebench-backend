package main

import (
	"coursebench-backend/internal/config"
	"coursebench-backend/internal/fiber"
	"coursebench-backend/pkg/database"
	"coursebench-backend/pkg/log"
	"coursebench-backend/pkg/mail"
	"coursebench-backend/pkg/modelRegister"
	_ "coursebench-backend/pkg/models"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	config.SetupViper()
	log.InitLog()
	database.InitDB()
	database.InitRedis()
	database.InitMinio()
	mail.InitSMTP()
	db := database.GetDB()
	err := db.Migrator().AutoMigrate(modelRegister.GetRegisteredTypes()...)
	if err != nil {
		log.Panicln(err)
	}
	database.Migrate()
	app := fiber.New()
	fiber.Routes(app)
	if err := app.Listen(config.FiberConfig.Listen); err != nil {
		panic(err)
	}
}
