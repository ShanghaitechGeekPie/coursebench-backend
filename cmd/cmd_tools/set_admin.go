package main

import (
	"coursebench-backend/pkg/database"
	"coursebench-backend/pkg/models"
	"gorm.io/gorm"
	syslog "log"
)

func SetAdmin(userID int, isAdmin bool) {
	syslog.Printf("Setting user %d to admin: %v", userID, isAdmin)
	db := database.GetDB()
	user := &models.User{}
	err := db.Transaction(func(tx *gorm.DB) error {
		err := tx.First(user, userID).Error
		if err != nil {
			return err
		}
		if isAdmin && user.IsCommunityAdmin {
			syslog.Fatalf("User %d is already community admin", userID)
		}
		err = tx.Model(user).Update("is_admin", isAdmin).Error
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		syslog.Fatalln(err)
	}
	syslog.Println("Finished (un)setting user to admin")
}

func SetCommunityAdmin(userID int, isCommunityAdmin bool) {
	syslog.Printf("Setting user %d to community admin: %v", userID, isCommunityAdmin)
	db := database.GetDB()
	user := &models.User{}
	err := db.Transaction(func(tx *gorm.DB) error {
		err := tx.First(user, userID).Error
		if err != nil {
			return err
		}
		if isCommunityAdmin && user.IsAdmin {
			syslog.Fatalf("User %d is already admin", userID)
		}
		err = tx.Model(user).Update("is_community_admin", isCommunityAdmin).Error
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		syslog.Fatalln(err)
	}
	syslog.Println("Finished (un)setting user to community admin")
}
