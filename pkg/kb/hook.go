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

package kb

import (
	"coursebench-backend/pkg/database"
	syslog "log"

	"gorm.io/gorm"
)

// OnCourseCommentChanged is called asynchronously when a comment is posted, updated, or deleted.
// It regenerates the markdown for the affected course and uploads it to S3.
func OnCourseCommentChanged(db *gorm.DB, courseID uint) {
	if !database.IsS3Enabled() {
		return
	}
	go func() {
		err := ExportSingleCourse(db, courseID)
		if err != nil {
			syslog.Printf("[KB] Failed to update knowledge base for course %d: %v\n", courseID, err)
		} else {
			syslog.Printf("[KB] Updated knowledge base for course %d\n", courseID)
		}
	}()
}
