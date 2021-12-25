package main

import (
	"coursebench-backend/pkg/database"
	"coursebench-backend/pkg/modelRegister"
	_ "coursebench-backend/pkg/models"
	"fmt"
)

func main() {
	fmt.Println("Hello, World!")
	db := database.GetDB()
	err := db.Migrator().AutoMigrate(modelRegister.GetRegisteredTypes()...)
	if err != nil {
		panic(err)
	}
}
