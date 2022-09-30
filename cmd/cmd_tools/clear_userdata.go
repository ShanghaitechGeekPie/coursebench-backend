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
