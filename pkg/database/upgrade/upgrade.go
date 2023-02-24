package upgrade

import (
	"coursebench-backend/pkg/database"
	"coursebench-backend/pkg/errors"
	"coursebench-backend/pkg/log"
	"coursebench-backend/pkg/models"
	"coursebench-backend/pkg/queries"
	"gorm.io/gorm"
)

// 更新数据库
func UpgradeDB() {
	CurrentDBVersion := 3
	db := database.GetDB()
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
		log.Println("Upgrading database version from 0 to 1...")
		fallthrough
	case 1:
		log.Println("Upgrading database version from 1 to 2...")
		fallthrough
	case 2:
		log.Println("Upgrading database version from 2 to 3...")
		UpgradeDBFrom2To3()
	case 3:
	default:
		log.Panicf("The version of database is: %d, which is newer than the backend.", metadata.DBVersion)
	}
	metadata.DBVersion = CurrentDBVersion
	if err = db.Save(&metadata).Error; err != nil {
		log.Panic(err)
	}
}

func UpgradeDBFrom2To3() {
	db := database.GetDB()
	// 创建占位教师“其他”，并为所有课程创建一个“其他”授课组
	err := db.Transaction(func(tx *gorm.DB) error {
		teacher := &models.Teacher{
			Model:  gorm.Model{ID: models.TEACHER_OTHER_ID},
			Name:   "其他",
			EamsID: -1,
		}
		if err := tx.Create(teacher).Error; err != nil {
			return errors.Wrap(err, errors.DatabaseError)
		}
		//teachers := []int{models.TEACHER_OTHER_ID}
		teachers := []int{models.TEACHER_OTHER_ID}
		var courses []models.Course
		if err := tx.Find(&courses).Error; err != nil {
			return errors.Wrap(err, errors.DatabaseError)
		}
		for _, c := range courses {
			_, err := queries.AddCourseGroup(tx, "", int(c.ID), teachers)
			if err != nil {
				return errors.Wrap(err, errors.DatabaseError)
			}
		}
		return nil
	})
	if err != nil {
		log.Panicln(err)
	}
}
