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
	"coursebench-backend/pkg/database"
	"coursebench-backend/pkg/models"
	"github.com/lib/pq"
	"gorm.io/gorm"
	syslog "log"
)

// ClearUserdata Drops all user data
// Be careful when using this function
func ClearUserdata() {
	syslog.Println("Clearing user data...")
	db := database.GetDB()
	if err := db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.User{}).Error; err != nil {
		syslog.Fatalln(err.Error())
	}
	if err := db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.Comment{}).Error; err != nil {
		syslog.Fatalln(err.Error())
	}
	if err := db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.CommentLike{}).Error; err != nil {
		syslog.Fatalln(err.Error())
	}
	if err := db.Session(&gorm.Session{AllowGlobalUpdate: true}).Model(&models.Course{}).Updates(map[string]interface{}{"comment_count": 0, "scores": pq.Int64Array{0, 0, 0, 0}}).Error; err != nil {
		syslog.Fatalln(err.Error())
	}
	if err := db.Session(&gorm.Session{AllowGlobalUpdate: true}).Model(&models.CourseGroup{}).Updates(map[string]interface{}{"comment_count": 0, "scores": pq.Int64Array{0, 0, 0, 0}}).Error; err != nil {
		syslog.Fatalln(err.Error())
	}

	syslog.Println("Finished clearing user data")
}
