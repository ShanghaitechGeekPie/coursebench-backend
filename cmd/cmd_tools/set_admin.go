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
