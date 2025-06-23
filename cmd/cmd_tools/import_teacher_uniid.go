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
	"encoding/json"
	"io/ioutil"
	"log"
	"sort"
	"strconv"

	"gorm.io/gorm"
)

func ImportTeacherUniID(filePath string) {
	log.Println("Reading teachers file from", filePath)
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	var teachersJSON map[string]map[string]string
	err = json.Unmarshal(file, &teachersJSON)
	if err != nil {
		log.Fatalf("Failed to parse JSON: %v", err)
	}

	teacherData := make(map[string]map[string][]int)
	for courseCode, teachers := range teachersJSON {
		if _, ok := teacherData[courseCode]; !ok {
			teacherData[courseCode] = make(map[string][]int)
		}
		for uniIDStr, teacherName := range teachers {
			uniID, err := strconv.Atoi(uniIDStr)
			if err != nil {
				log.Printf("Warning: could not parse uniID '%s' for teacher '%s' in course '%s'. Skipping.", uniIDStr, teacherName, courseCode)
				continue
			}
			if _, ok := teacherData[courseCode][teacherName]; !ok {
				teacherData[courseCode][teacherName] = make([]int, 0)
			}
			teacherData[courseCode][teacherName] = append(teacherData[courseCode][teacherName], uniID)
		}
	}

	for _, course := range teacherData {
		for _, ids := range course {
			sort.Ints(ids)
		}
	}

	db := database.GetDB()
	err = db.Transaction(func(tx *gorm.DB) error {
		var courseGroups []models.CourseGroup
		err = tx.Preload("Course").Preload("Teachers").Find(&courseGroups).Error
		if err != nil {
			return err
		}

		for _, cg := range courseGroups {
			if cg.Course.Code == "" {
				continue
			}
			courseCode := cg.Course.Code
			courseTeachers, ok := teacherData[courseCode]
			if !ok {
				continue
			}

			for _, teacher := range cg.Teachers {
				teacherUniIDs, ok := courseTeachers[teacher.Name]
				if !ok {
					continue
				}

				if len(teacherUniIDs) > 0 {
					uniID := teacherUniIDs[0]
					if teacher.UniID != uniID {
						log.Printf("Updating teacher '%s' (ID: %d) UniID from %d to %d for course '%s'", teacher.Name, teacher.ID, teacher.UniID, uniID, courseCode)
						err := tx.Model(teacher).Update("uni_id", uniID).Error
						if err != nil {
							log.Printf("ERROR: Failed to update teacher UniID: %v", err)
						}
					}
				}
			}
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Transaction failed: %v", err)
	}

	log.Println("Successfully finished updating teacher UniIDs.")

	// 生成工号-姓名-teacher 映射并打印部分内容
	mapping := GenerateUniIDNameTeacherMap()
	log.Println("工号-姓名-teacher 映射样例:")
	count := 0
	for uniID, nameMap := range mapping {
		for name, teacher := range nameMap {
			log.Printf("工号: %d, 姓名: %s, TeacherID: %d", uniID, name, teacher.ID)
			count++
			if count > 10 {
				break
			}
		}
		if count > 10 {
			break
		}
	}
}

func GenerateUniIDNameTeacherMap() map[int]map[string]*models.Teacher {
	db := database.GetDB()
	var teachers []models.Teacher
	err := db.Find(&teachers).Error
	if err != nil {
		log.Fatalf("Failed to query teachers: %v", err)
	}
	result := make(map[int]map[string]*models.Teacher)
	for i := range teachers {
		uniID := teachers[i].UniID
		name := teachers[i].Name
		if _, ok := result[uniID]; !ok {
			result[uniID] = make(map[string]*models.Teacher)
		}
		result[uniID][name] = &teachers[i]
	}
	return result
}

func ImportAndFixTeacherUniIDAndRelations(jsonPath string) {
	log.Println("开始导入并修正 teacher uniid 及 course_group_teachers 关系，teachers.json 路径:", jsonPath)

	// 1. 读取 teachers.json
	file, err := ioutil.ReadFile(jsonPath)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}
	var teachersJSON map[string]map[string]string
	err = json.Unmarshal(file, &teachersJSON)
	if err != nil {
		log.Fatalf("Failed to parse JSON: %v", err)
	}

	// 2. 建立 course_code => name => uniID 列表
	courseTeacherUniID := make(map[string]map[string][]int)
	for courseCode, teachers := range teachersJSON {
		if _, ok := courseTeacherUniID[courseCode]; !ok {
			courseTeacherUniID[courseCode] = make(map[string][]int)
		}
		for uniIDStr, teacherName := range teachers {
			uniID, err := strconv.Atoi(uniIDStr)
			if err != nil {
				log.Printf("Warning: could not parse uniID '%s' for teacher '%s' in course '%s'. Skipping.", uniIDStr, teacherName, courseCode)
				continue
			}
			courseTeacherUniID[courseCode][teacherName] = append(courseTeacherUniID[courseCode][teacherName], uniID)
		}
	}
	for _, nameMap := range courseTeacherUniID {
		for _, ids := range nameMap {
			sort.Ints(ids)
		}
	}

	db := database.GetDB()
	err = db.Transaction(func(tx *gorm.DB) error {
		type CGT struct {
			CourseGroupID uint
			TeacherID     uint
		}
		var cgts []CGT
		if err := tx.Table("coursegroup_teachers").Find(&cgts).Error; err != nil {
			return err
		}

		for _, cgt := range cgts {
			// 找到 course_code
			var group models.CourseGroup
			if err := tx.Preload("Course").Where("id = ?", cgt.CourseGroupID).First(&group).Error; err != nil {
				log.Printf("Skip cgt (%d, %d): cannot find course_group %d", cgt.CourseGroupID, cgt.TeacherID, cgt.CourseGroupID)
				continue
			}
			courseCode := group.Course.Code

			// 找到 teacher_name
			var teacher models.Teacher
			if err := tx.Where("id = ?", cgt.TeacherID).First(&teacher).Error; err != nil {
				log.Printf("Skip cgt (%d, %d): cannot find teacher %d", cgt.CourseGroupID, cgt.TeacherID, cgt.TeacherID)
				continue
			}
			teacherName := teacher.Name

			// 查找正确工号
			uniIDs := courseTeacherUniID[courseCode][teacherName]
			if len(uniIDs) == 0 {
				log.Printf("Skip cgt (%d, %d): cannot find uniID for %s in course %s", cgt.CourseGroupID, cgt.TeacherID, teacherName, courseCode)
				continue
			}
			uniID := uniIDs[0] // 按工号升序取最小

			// 查找/新建 Teacher，并修正 uniid 字段
			var rightTeacher models.Teacher
			if err := tx.Where("uni_id = ? AND name = ?", uniID, teacherName).First(&rightTeacher).Error; err != nil {
				// 新建
				rightTeacher = models.Teacher{
					UniID: uniID,
					Name:  teacherName,
				}
				if err := tx.Create(&rightTeacher).Error; err != nil {
					log.Printf("Failed to create teacher %s uniID %d: %v", teacherName, uniID, err)
					continue
				}
				log.Printf("Created new teacher: %s uniID %d id %d", teacherName, uniID, rightTeacher.ID)
			} else {
				// 如果 Teacher 的 uniid 不对，修正
				if rightTeacher.UniID != uniID {
					if err := tx.Model(&rightTeacher).Update("uni_id", uniID).Error; err != nil {
						log.Printf("Failed to update teacher %s id %d uniID to %d: %v", teacherName, rightTeacher.ID, uniID, err)
						continue
					}
					log.Printf("Updated teacher %s id %d uniID to %d", teacherName, rightTeacher.ID, uniID)
				}
			}

			// 修正关系（如 teacher_id 需变更，先删后插）
			if cgt.TeacherID != rightTeacher.ID {
				if err := tx.Table("coursegroup_teachers").
					Where("course_group_id = ? AND teacher_id = ?", cgt.CourseGroupID, cgt.TeacherID).
					Delete(nil).Error; err != nil {
					log.Printf("Failed to delete old cgt (%d, %d): %v", cgt.CourseGroupID, cgt.TeacherID, err)
					continue
				}
				insert := map[string]interface{}{
					"course_group_id": cgt.CourseGroupID,
					"teacher_id":      rightTeacher.ID,
				}
				if err := tx.Table("coursegroup_teachers").Create(insert).Error; err != nil {
					log.Printf("Failed to insert new cgt (%d, %d): %v", cgt.CourseGroupID, rightTeacher.ID, err)
					continue
				}
				log.Printf("Fixed cgt: (%d, %d) => (%d, %d)", cgt.CourseGroupID, cgt.TeacherID, cgt.CourseGroupID, rightTeacher.ID)
			}
		}

		// 修复 course_teachers
		type CT struct {
			CourseID  uint
			TeacherID uint
		}
		var cts []CT
		if err := tx.Table("course_teachers").Find(&cts).Error; err != nil {
			return err
		}
		for _, ct := range cts {
			// 找到 course_code
			var course models.Course
			if err := tx.Where("id = ?", ct.CourseID).First(&course).Error; err != nil {
				log.Printf("Skip ct (%d, %d): cannot find course %d", ct.CourseID, ct.TeacherID, ct.CourseID)
				continue
			}
			courseCode := course.Code

			// 找到 teacher_name
			var teacher models.Teacher
			if err := tx.Where("id = ?", ct.TeacherID).First(&teacher).Error; err != nil {
				log.Printf("Skip ct (%d, %d): cannot find teacher %d", ct.CourseID, ct.TeacherID, ct.TeacherID)
				continue
			}
			teacherName := teacher.Name

			// 查找正确工号
			uniIDs := courseTeacherUniID[courseCode][teacherName]
			if len(uniIDs) == 0 {
				log.Printf("Skip ct (%d, %d): cannot find uniID for %s in course %s", ct.CourseID, ct.TeacherID, teacherName, courseCode)
				continue
			}
			uniID := uniIDs[0]

			// 查找/新建 Teacher，并修正 uniid 字段
			var rightTeacher models.Teacher
			if err := tx.Where("uni_id = ? AND name = ?", uniID, teacherName).First(&rightTeacher).Error; err != nil {
				rightTeacher = models.Teacher{
					UniID: uniID,
					Name:  teacherName,
				}
				if err := tx.Create(&rightTeacher).Error; err != nil {
					log.Printf("Failed to create teacher %s uniID %d: %v", teacherName, uniID, err)
					continue
				}
				log.Printf("Created new teacher: %s uniID %d id %d", teacherName, uniID, rightTeacher.ID)
			} else {
				if rightTeacher.UniID != uniID {
					if err := tx.Model(&rightTeacher).Update("uni_id", uniID).Error; err != nil {
						log.Printf("Failed to update teacher %s id %d uniID to %d: %v", teacherName, rightTeacher.ID, uniID, err)
						continue
					}
					log.Printf("Updated teacher %s id %d uniID to %d", teacherName, rightTeacher.ID, uniID)
				}
			}

			// 修正关系（如 teacher_id 需变更，先删后插）
			if ct.TeacherID != rightTeacher.ID {
				if err := tx.Table("course_teachers").
					Where("course_id = ? AND teacher_id = ?", ct.CourseID, ct.TeacherID).
					Delete(nil).Error; err != nil {
					log.Printf("Failed to delete old ct (%d, %d): %v", ct.CourseID, ct.TeacherID, err)
					continue
				}
				insert := map[string]interface{}{
					"course_id":  ct.CourseID,
					"teacher_id": rightTeacher.ID,
				}
				if err := tx.Table("course_teachers").Create(insert).Error; err != nil {
					log.Printf("Failed to insert new ct (%d, %d): %v", ct.CourseID, rightTeacher.ID, err)
					continue
				}
				log.Printf("Fixed ct: (%d, %d) => (%d, %d)", ct.CourseID, ct.TeacherID, ct.CourseID, rightTeacher.ID)
			}
		}

		// 删除孤儿老师
		if err := tx.Exec(`
			DELETE FROM teachers
			WHERE id NOT IN (SELECT teacher_id FROM course_teachers)
			  AND id NOT IN (SELECT teacher_id FROM coursegroup_teachers)
		`).Error; err != nil {
			log.Printf("Failed to delete orphan teachers: %v", err)
		} else {
			log.Println("Orphan teachers deleted.")
		}

		// 合并同名老师（未在json出现但数据库有同名有uniid的老师）
		jsonNameUniID := make(map[string]map[int]struct{})
		for _, teachers := range teachersJSON {
			for uniIDStr, name := range teachers {
				uniID, err := strconv.Atoi(uniIDStr)
				if err != nil {
					continue
				}
				if _, ok := jsonNameUniID[name]; !ok {
					jsonNameUniID[name] = make(map[int]struct{})
				}
				jsonNameUniID[name][uniID] = struct{}{}
			}
		}

		var allTeachers []models.Teacher
		if err := tx.Find(&allTeachers).Error; err != nil {
			return err
		}
		for _, t := range allTeachers {
			if t.UniID > 0 {
				if m, ok := jsonNameUniID[t.Name]; ok {
					if _, ok2 := m[t.UniID]; ok2 {
						continue
					}
				}
			}
			var target models.Teacher
			if err := tx.Where("name = ? AND uni_id > 0", t.Name).First(&target).Error; err != nil {
				continue
			}
			if target.ID == t.ID {
				continue
			}
			if err := tx.Exec("UPDATE course_teachers SET teacher_id = ? WHERE teacher_id = ?", target.ID, t.ID).Error; err != nil {
				log.Printf("Failed to migrate course_teachers for %s: %v", t.Name, err)
			}
			if err := tx.Exec("UPDATE coursegroup_teachers SET teacher_id = ? WHERE teacher_id = ?", target.ID, t.ID).Error; err != nil {
				log.Printf("Failed to migrate coursegroup_teachers for %s: %v", t.Name, err)
			}
			log.Printf("Merged teacher %s (id %d) into teacher id %d", t.Name, t.ID, target.ID)
		}

		// 再次删除孤儿老师
		if err := tx.Exec(`
			DELETE FROM teachers
			WHERE id NOT IN (SELECT teacher_id FROM course_teachers)
			  AND id NOT IN (SELECT teacher_id FROM coursegroup_teachers)
		`).Error; err != nil {
			log.Printf("Failed to delete orphan teachers: %v", err)
		} else {
			log.Println("Orphan teachers deleted after merge.")
		}

		return nil
	})
	if err != nil {
		log.Fatalf("导入并修正 teacher uniid 及关系失败: %v", err)
	}
	log.Println("导入并修正 teacher uniid 及关系完成")
}
