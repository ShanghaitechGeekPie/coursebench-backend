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
	"coursebench-backend/internal/fiber"
	"coursebench-backend/pkg/database"
	"coursebench-backend/pkg/database/upgrade"
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
	upgrade.UpgradeDB()
	app := fiber.New()
	fiber.Routes(app)
	if err := app.Listen(config.FiberConfig.Listen); err != nil {
		panic(err)
	}
}
