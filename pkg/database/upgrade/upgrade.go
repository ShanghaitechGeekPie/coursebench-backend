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
	CurrentDBVersion := 5
	db := database.GetDB()
	var metadata models.Metadata
	err := db.Take(&metadata).Error
	if err != nil {
		// 初次创建数据库
		metadata = models.Metadata{DBVersion: CurrentDBVersion}
		DatabaseInit()
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
		fallthrough
	case 3:
		log.Println("Upgrading database version from 3 to 4...")
		UpgradeDBFrom3To4()
		fallthrough
	case 4:
		log.Println("Upgrading database version from 4 to 5...")
		UpgradeDBFrom4To5()
	case 5:
	default:
		log.Panicf("The version of database is: %d, which is newer than the backend.", metadata.DBVersion)
	}
	metadata.DBVersion = CurrentDBVersion
	if err = db.Save(&metadata).Error; err != nil {
		log.Panic(err)
	}
}

func DatabaseInit() {
	db := database.GetDB()
	// 创建占位教师“其他”
	err := db.Transaction(func(tx *gorm.DB) error {
		teacher := &models.Teacher{
			Model:  gorm.Model{ID: models.TEACHER_OTHER_ID},
			Name:   "其他",
			EamsID: -1,
			UniID:  -1,
		}
		if err := tx.Create(teacher).Error; err != nil {
			return errors.Wrap(err, errors.DatabaseError)
		}
		return nil
	})
	if err != nil {
		log.Panicln(err)
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

func UpgradeDBFrom3To4() {
	db := database.GetDB()
	err := db.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&models.User{}).Where("reward is null").Update("reward", 0).Error
		if err != nil {
			return errors.Wrap(err, errors.DatabaseError)
		}
		return nil
	})
	if err != nil {
		log.Panicln(err)
	}
}

func UpgradeDBFrom4To5() {
	db := database.GetDB()
	// 升级到v5：添加成就系统的唯一约束
	err := db.Transaction(func(tx *gorm.DB) error {
		// 为 user_achievements 表添加唯一约束（如果表已存在）
		// GORM 的 AutoMigrate 会自动创建表，但我们需要确保唯一约束存在

		// 如果表已存在，添加唯一索引（处理已有数据的情况）
		// 首先删除可能存在的重复数据
		result := tx.Exec(`
			DELETE FROM user_achievements a USING (
				SELECT MIN(id) as id, user_id, achievement_id
				FROM user_achievements 
				GROUP BY user_id, achievement_id HAVING COUNT(*) > 1
			) b
			WHERE a.user_id = b.user_id 
			AND a.achievement_id = b.achievement_id 
			AND a.id <> b.id
		`)
		if result.Error != nil {
			// 表可能不存在，这是正常的
			log.Println("Note: user_achievements table may not exist yet, will be created by AutoMigrate")
		}

		// 添加唯一索引（如果不存在）
		result = tx.Exec(`
			CREATE UNIQUE INDEX IF NOT EXISTS idx_user_achievement 
			ON user_achievements(user_id, achievement_id) 
			WHERE deleted_at IS NULL
		`)
		if result.Error != nil {
			// 索引可能已存在或表不存在
			log.Println("Note: unique index creation skipped or already exists")
		}

		return nil
	})
	if err != nil {
		log.Panicln(err)
	}
}
