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
	"encoding/csv"
	"gorm.io/gorm"
	syslog "log"
	"os"
)

func ImportTeacher(filePath string) {
	syslog.Printf("Starting to import teachers' data: %s\n", filePath)
	db := database.GetDB()
	csvFile, err := os.Open(filePath)
	if err != nil {
		syslog.Fatalln(err)
	}
	records, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		syslog.Fatalln(err)
	}
	for i, record := range records {
		if len(record) < 6 {
			syslog.Printf("Wrong format at line:%d", i)
			continue
		}
		name := record[0]
		photo := record[1]
		job := record[2]
		email := record[3]
		institute := record[4]
		introduction := record[5]
		teacher := &models.Teacher{}
		err = db.Where("name = ?", name).Take(teacher).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			syslog.Fatalln(err)
		}
		if err != nil {
			syslog.Printf("Teacher %s not found", name)
			continue
		}
		teacher.Email = email
		teacher.Introduction = introduction
		teacher.Job = job
		teacher.Photo = photo
		teacher.Institute = institute
		if err = db.Save(teacher).Error; err != nil {
			syslog.Fatalln(err)
		}
		syslog.Printf("Teacher %s updated", name)
	}
	syslog.Println("Finished importing teacher data")
}
