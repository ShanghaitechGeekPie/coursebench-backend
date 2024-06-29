// Copyright (C) 2021-2024 ShanghaiTech GeekPie
// This file is part of CourseBench Backend.
//
// CourseBench Backend is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// CourseBench Backend is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with CourseBench Backend.  If not, see <http://www.gnu.org/licenses/>.

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
