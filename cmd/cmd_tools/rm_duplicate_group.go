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
	"fmt"
	"log"
	"sort"
)

// Helper: 返回教师ID升序拼接的字符串
func teacherIDKey(ids []uint) string {
	copyIDs := append([]uint(nil), ids...)
	sort.Slice(copyIDs, func(i, j int) bool { return copyIDs[i] < copyIDs[j] })
	key := ""
	for i, id := range copyIDs {
		if i > 0 {
			key += ","
		}
		key += fmt.Sprintf("%d", id)
	}
	return key
}

func RmDuplicateCourseGroup() {
	db := database.GetDB()
	var courses []models.Course
	if err := db.Preload("Groups.Teachers").Preload("Groups.Comment").Find(&courses).Error; err != nil {
		log.Fatalf("Failed to load courses: %v", err)
	}
	for _, course := range courses {
		groupMap := make(map[string][]*models.CourseGroup)
		for i := range course.Groups {
			group := &course.Groups[i]
			var teacherIDs []uint
			for _, t := range group.Teachers {
				teacherIDs = append(teacherIDs, t.ID)
			}
			key := teacherIDKey(teacherIDs)
			groupMap[key] = append(groupMap[key], group)
		}
		for key, groups := range groupMap {
			if len(groups) <= 1 {
				continue
			}
			log.Printf("Course %s has %d duplicate groups for teacher set [%s]", course.Code, len(groups), key)
			// 选第一个为保留组
			mainGroup := groups[0]
			var mergedScores []int64
			if len(mainGroup.Scores) > 0 {
				mergedScores = append([]int64(nil), mainGroup.Scores...)
			} else {
				mergedScores = make([]int64, 4)
			}
			mergedCommentCount := mainGroup.CommentCount
			// 合并分数和评论
			for _, dup := range groups[1:] {
				for i := 0; i < len(mergedScores) && i < len(dup.Scores); i++ {
					mergedScores[i] += dup.Scores[i]
				}
				mergedCommentCount += dup.CommentCount
				// 评论迁移
				for _, c := range dup.Comment {
					c.CourseGroupID = mainGroup.ID
					if err := db.Save(&c).Error; err != nil {
						log.Printf("Failed to move comment %d: %v", c.ID, err)
					}
				}
				// TODO: 迁移其它依赖于 course_group_id 的表（如点赞等）
				if err := db.Delete(&models.CourseGroup{}, dup.ID).Error; err != nil {
					log.Printf("Failed to delete duplicate group %d: %v", dup.ID, err)
				}
				log.Printf("Merged and deleted duplicate group %d into %d", dup.ID, mainGroup.ID)
			}
			mainGroup.Scores = mergedScores
			mainGroup.CommentCount = mergedCommentCount
			if err := db.Save(mainGroup).Error; err != nil {
				log.Printf("Failed to update main group %d: %v", mainGroup.ID, err)
			}
		}
	}
	log.Println("Duplicate course group merge complete.")
}
