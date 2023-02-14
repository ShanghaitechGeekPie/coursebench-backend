package database

import (
	"coursebench-backend/pkg/log"
	"coursebench-backend/pkg/models"
)

// 更新数据库
func Migrate() {
	CurrentDBVersion := 2
	db := GetDB()
	var metadata models.Metadata
	err := db.Take(&metadata).Error
	if err != nil {
		// 初次创建数据库
		metadata = models.Metadata{DBVersion: CurrentDBVersion}
	}
	// 依次更新 DB Version
	// fallthrough 为数不多的用途
	switch metadata.DBVersion {
	case 0:
		log.Println("Updating database version from 0 to 1...")
		fallthrough
	case 1:
		log.Println("Updating database version from 1 to 2...")
		fallthrough
	case 2:
		log.Println("Updating database version from 2 to 3...")
		UpdateFrom2To3()
	default:
		log.Panicf("The version of database is: %d, which is newer than the backend.", metadata.DBVersion)
	}
	metadata.DBVersion = CurrentDBVersion
	if err = db.Save(&metadata).Error; err != nil {
		log.Panic(err)
	}
}

func UpdateFrom2To3() {

}
