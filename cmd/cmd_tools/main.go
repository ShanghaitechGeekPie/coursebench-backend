package main

import (
	"coursebench-backend/internal/config"
	"coursebench-backend/pkg/database"
	"coursebench-backend/pkg/log"
	"coursebench-backend/pkg/modelRegister"
	syslog "log"
	"os"
	"strconv"
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
	case "set_admin":
		if len(args) < 3 {
			syslog.Fatalln("Missing parameters <user id>")
		}
		userId, err := strconv.Atoi(args[2])
		if err != nil {
			syslog.Fatalln(err)
		}
		SetAdmin(userId, true)
	case "unset_admin":
		if len(args) < 3 {
			syslog.Fatalln("Missing parameters <user id>")
		}
		userId, err := strconv.Atoi(args[2])
		if err != nil {
			syslog.Fatalln(err)
		}
		SetAdmin(userId, false)
	case "set_community_admin":
		if len(args) < 3 {
			syslog.Fatalln("Missing parameters <user id>")
		}
		userId, err := strconv.Atoi(args[2])
		if err != nil {
			syslog.Fatalln(err)
		}
		SetCommunityAdmin(userId, true)
	case "unset_community_admin":
		if len(args) < 3 {
			syslog.Fatalln("Missing parameters <user id>")
		}
		userId, err := strconv.Atoi(args[2])
		if err != nil {
			syslog.Fatalln(err)
		}
		SetCommunityAdmin(userId, false)
	default:
		syslog.Fatal("Unknown command!")
	}
}
