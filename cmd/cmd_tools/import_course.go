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
	"coursebench-backend/pkg/errors"
	"coursebench-backend/pkg/models"
	"coursebench-backend/pkg/queries"
	"encoding/csv"
	"encoding/json"
	"github.com/lib/pq"
	"gorm.io/gorm"
	syslog "log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

func ImportCourse(filePath string) {
	syslog.Printf("Starting to import courses' data from directory: %s\n", filePath)
	db := database.GetDB()
	path, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	pathInfo, err := path.Stat()
	if err != nil {
		panic(err)
	}
	if !pathInfo.IsDir() {
		panic("Given path should be a directory!")
	}
	err = path.Close()
	if err != nil {
		panic(err)
	}
	err = db.Transaction(func(tx *gorm.DB) error {
		return filepath.Walk(filePath, func(file string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}

			syslog.Printf("Starting to import courses' data: %s\n", file)
			csvFile, err := os.Open(file)
			if err != nil {
				panic(err)
			}
			reader := csv.NewReader(csvFile)
			reader.LazyQuotes = true
			records, err := reader.ReadAll()
			if err != nil {
				panic(err)
			}
			for i, record := range records {
				if i == 0 {
					continue
				}
				code := record[3]
				name := record[2]
				institute := record[10]
				creditF, err := strconv.ParseFloat(record[4], 64)
				if err != nil {
					panic(err)
				}
				credit := int(creditF)
				course := &models.Course{}
				err = tx.Where("code = ?", code).Take(course).Error
				if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
					panic(err)
				}
				if err != nil {
					syslog.Printf("Add course %s %s %s %d\n", code, name, institute, credit)
					course, err = queries.AddCourse(tx, name, institute, credit, code)
					if err != nil {
						panic(err)
					}
					other_teachers := []int{models.TEACHER_OTHER_ID}
					_, err = queries.AddCourseGroup(tx, "", int(course.ID), other_teachers)
					if err != nil {
						return errors.Wrap(err, errors.DatabaseError)
					}
				} else {
					syslog.Printf("Find course %s %d\n", code, course.ID)
				}
				t := strings.ReplaceAll(record[12], `'`, `"`)
				var teacherNames []string
				var teacherEamsIDs []int
				var teacherIDs []int
				err = json.Unmarshal([]byte(record[13]), &teacherEamsIDs)
				if err != nil {
					panic(err)
				}
				err = json.Unmarshal([]byte(t), &teacherNames)
				if err != nil {
					panic(err)
				}
				var data []struct {
					CourseGroupID int
					ArrayAgg      pq.StringArray `gorm:"type:string[]"`
				}
				tx.Raw(`select  course_group_id, array_agg(name) from course_groups
    inner join coursegroup_teachers on course_groups.id = coursegroup_teachers.course_group_id
    inner join teachers on coursegroup_teachers.teacher_id = teachers.id
    where course_groups.course_id=? group by course_group_id;`, course.ID).Scan(&data)
				names := make([]string, len(teacherNames))
				copy(names, teacherNames)
				sort.Strings(names)
				flag := false
				for _, c := range data {
					sort.Strings(c.ArrayAgg)
					flag2 := true
					if len(c.ArrayAgg) == len(names) {
						for i, name := range c.ArrayAgg {
							if name != names[i] {
								flag2 = false
								break
							}
						}
					} else {
						flag2 = false
					}
					if flag2 {
						flag = true
						break
					}
				}
				if flag {
					continue
				}
				for i, name := range teacherNames {
					teacher := &models.Teacher{}
					err = tx.Where("name = ?", name).Take(teacher).Error
					if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
						panic(err)
					}
					if err != nil {
						teacher, err = queries.AddTeacher(tx, name, "", "", teacherEamsIDs[i])
						if err != nil {
							panic(err)
						}
					}
					teacherIDs = append(teacherIDs, int(teacher.ID))
				}
				syslog.Println(teacherIDs)
				group, err := queries.AddCourseGroup(tx, "", int(course.ID), teacherIDs)
				if err != nil {
					panic(err)
				}
				syslog.Printf("Add course group %s %d %d\n", code, course.ID, group.ID)
			}

			syslog.Printf("Finished importing courses' data: %s\n", file)
			return nil
		})
	})
	syslog.Printf("Finished importing courses' data from directory: %s\n", filePath)
}
