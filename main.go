package main

import (
	"coursebench-backend/pkg/database"
	"fmt"
)

func main() {
	fmt.Println("Hello, World!")
	_ = database.GetDB()
}
