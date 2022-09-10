package main

import (
	"coursebench-backend/internal/config"
	"coursebench-backend/pkg/database"
	"coursebench-backend/pkg/errors"
	"coursebench-backend/pkg/log"
	"coursebench-backend/pkg/modelRegister"
	"coursebench-backend/pkg/models"
	"encoding/csv"
	"gorm.io/gorm"
	syslog "log"
	"os"
)

func main() {
	config.SetupViper()
	log.InitLog()
	database.InitDB()
	syslog.Println("Starting to import teacher data...")
	db := database.GetDB()
	err := db.Migrator().AutoMigrate(modelRegister.GetRegisteredTypes()...)
	if err != nil {
		syslog.Fatalln(err)
	}

	csvFile, err := os.Open("data_import/teacher.csv")
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
