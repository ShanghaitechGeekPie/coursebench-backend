package main

import (
	"coursebench-backend/database"
	"fmt"
)

func main() {
	fmt.Println("Hello, World!")
	_ = database.GetDB()
}
