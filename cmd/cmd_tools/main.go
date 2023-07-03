package main

import (
	"coursebench-backend/internal/config"
	"coursebench-backend/pkg/database"
	"coursebench-backend/pkg/log"
	"coursebench-backend/pkg/modelRegister"
	syslog "log"
	"os"
)

func main() {
	config.SetupViper()
	log.InitLog()
	syslog.Println("Backend Command Line Tool Starting...")
	database.InitDB()
	database.InitRedis()
	database.InitMinio()
	db := database.GetDB()
	err := db.Migrator().AutoMigrate(modelRegister.GetRegisteredTypes()...)
	if err != nil {
		panic(err)
	}
	args := os.Args
	if len(args) < 2 {
		syslog.Fatal("No command specified")
	}
	switch args[1] {
	case "import_teacher":
		if len(args) < 3 {
			syslog.Fatalln("Missing parameters <file path>")
		}
		filePath := args[2]
		ImportTeacher(filePath)
	case "import_course":
		if len(args) < 3 {
			syslog.Fatalln("Missing parameters <file path>")
		}
		filePath := args[2]
		ImportCourse(filePath)
	case "clear_userdata":
		if len(args) != 3 || args[2] != "Yes_Confirm" {
			syslog.Fatalln("Wrong check code")
		}
		ClearUserdata()
	default:
		syslog.Fatal("Unknown command!")
	}
}
