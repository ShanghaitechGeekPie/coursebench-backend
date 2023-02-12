package database

import (
	"coursebench-backend/pkg/log"
	"coursebench-backend/pkg/models"
)

// 更新数据库
func Migrate() {
	CurrentDBVersion := 1
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
		fallthrough
	case 1:
	default:
		log.Panicf("The version of database is: %d, which is newer than the backend.", metadata.DBVersion)
	}
	metadata.DBVersion = CurrentDBVersion
	if err = db.Save(&metadata).Error; err != nil {
		log.Panic(err)
	}
}
