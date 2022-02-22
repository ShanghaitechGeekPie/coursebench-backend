package main

import (
	"coursebench-backend/internal/config"
	"coursebench-backend/pkg/database"
	"coursebench-backend/pkg/errors"
	"coursebench-backend/pkg/modelRegister"
	"coursebench-backend/pkg/models"
	_ "coursebench-backend/pkg/models"
	"coursebench-backend/pkg/queries"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"os"
	"strconv"
	"strings"
)

func main() {
	config.SetupViper()
	database.InitDB()
	db := database.GetDB()
	err := db.Migrator().AutoMigrate(modelRegister.GetRegisteredTypes()...)
	if err != nil {
		panic(err)
	}

	csvFile, err := os.Open("data_import/course.csv")
	if err != nil {
		panic(err)
	}
	records, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		panic(err)
	}
	/*for _, record := range records {
		fmt.Println(record)
	}*/
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
		err = db.Where("code = ?", code).Take(course).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			panic(err)
		}
		if err != nil {
			course, err = queries.AddCourse(name, institute, credit, code)
			if err != nil {
				panic(err)
			}
		} else {
			fmt.Printf("Find course %s %d\n", code, course.ID)
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
		for i, name := range teacherNames {
			teacher := &models.Teacher{}
			err = db.Where("name = ?", name).Take(teacher).Error
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				panic(err)
			}
			if err != nil {
				teacher, err = queries.AddTeacher(name, "", "", teacherEamsIDs[i])
				if err != nil {
					panic(err)
				}
			}
			teacherIDs = append(teacherIDs, int(teacher.ID))
		}
		fmt.Println(teacherIDs)
		group, err := queries.AddCourseGroup("", int(course.ID), teacherIDs)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Add course group %s %d %d\n", code, course.ID, group.ID)
	}
}
