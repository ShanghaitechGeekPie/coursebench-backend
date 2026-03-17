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
	"coursebench-backend/pkg/kb"
	syslog "log"
	"strconv"
)

// ExportKB exports all courses with comments to S3 as RAG-friendly markdown
func ExportKB() {
	db := database.GetDB()
	err := kb.ExportAllCourses(db)
	if err != nil {
		syslog.Fatalf("Failed to export knowledge base: %v\n", err)
	}
	syslog.Println("Knowledge base export completed successfully")
}

// ExportKBCourse exports a single course to S3 as RAG-friendly markdown
func ExportKBCourse(courseIDStr string) {
	courseID, err := strconv.Atoi(courseIDStr)
	if err != nil {
		syslog.Fatalf("Invalid course ID: %s\n", courseIDStr)
	}
	db := database.GetDB()
	err = kb.ExportSingleCourse(db, uint(courseID))
	if err != nil {
		syslog.Fatalf("Failed to export course %d: %v\n", courseID, err)
	}
	syslog.Printf("Course %d exported successfully\n", courseID)
}
