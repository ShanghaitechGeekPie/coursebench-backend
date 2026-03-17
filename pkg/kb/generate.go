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
	"bytes"
	"context"
	"coursebench-backend/pkg/database"
	"coursebench-backend/pkg/models"
	"fmt"
	syslog "log"
	"path"
	"strings"
	"time"

	"gorm.io/gorm"
)

var scoreLabels = []string{"内容", "教学", "给分", "收获"}

// FormatSemester converts semester int (e.g. 202401) to human-readable string
func FormatSemester(semester int) string {
	year := semester / 100
	sem := semester % 100
	semName := ""
	switch sem {
	case 1:
		semName = "秋季"
	case 2:
		semName = "春季"
	case 3:
		semName = "夏季"
	}
	return fmt.Sprintf("%d-%d %s", year, year+1, semName)
}

// FormatScoreRanking converts ranking int to readable string
func FormatScoreRanking(ranking int) string {
	switch ranking {
	case 1:
		return "前5%"
	case 2:
		return "5%-10%"
	case 3:
		return "10%-20%"
	case 4:
		return "20%-30%"
	case 5:
		return "30%-40%"
	case 6:
		return "40%-50%"
	case 7:
		return "50%-60%"
	case 8:
		return "60%-70%"
	case 9:
		return "70%-80%"
	case 10:
		return "80%-90%"
	case 11:
		return "90%-100%"
	default:
		return "未知"
	}
}

// SanitizeFilename removes characters unsafe for filenames
func SanitizeFilename(s string) string {
	replacer := strings.NewReplacer(
		"/", "_", "\\", "_", ":", "_", "*", "_",
		"?", "_", "\"", "_", "<", "_", ">", "_",
		"|", "_", " ", "_",
	)
	return replacer.Replace(s)
}

// GenerateCourseFilename creates a RAG-friendly filename with metadata embedded
// Format: course_{id}_{code}_{name}_{institute}.md
func GenerateCourseFilename(course *models.Course) string {
	name := SanitizeFilename(course.Name)
	code := SanitizeFilename(course.Code)
	institute := SanitizeFilename(course.Institute)
	return fmt.Sprintf("course_%d_%s_%s_%s.md", course.ID, code, name, institute)
}

// GenerateCourseMarkdown generates RAG-friendly markdown for a single course with all its comments
func GenerateCourseMarkdown(course *models.Course, groups []models.CourseGroup, comments []models.Comment) string {
	var buf bytes.Buffer
	now := time.Now().Format(time.RFC3339)

	// Collect all teachers across groups
	teacherMap := make(map[uint]models.Teacher)
	for _, g := range groups {
		for _, t := range g.Teachers {
			teacherMap[t.ID] = *t
		}
	}
	teacherNames := make([]string, 0, len(teacherMap))
	for _, t := range teacherMap {
		teacherNames = append(teacherNames, t.Name)
	}

	// Calculate average scores
	avgScores := make([]float64, models.ScoreLength)
	if course.CommentCount > 0 {
		for i := 0; i < models.ScoreLength; i++ {
			avgScores[i] = float64(course.Scores[i]) / float64(course.CommentCount)
		}
	}

	// YAML frontmatter
	buf.WriteString("---\n")
	buf.WriteString(fmt.Sprintf("course_id: %d\n", course.ID))
	buf.WriteString(fmt.Sprintf("course_name: \"%s\"\n", course.Name))
	buf.WriteString(fmt.Sprintf("course_code: \"%s\"\n", course.Code))
	buf.WriteString(fmt.Sprintf("institute: \"%s\"\n", course.Institute))
	buf.WriteString(fmt.Sprintf("credit: %d\n", course.Credit))
	buf.WriteString(fmt.Sprintf("comment_count: %d\n", course.CommentCount))
	buf.WriteString(fmt.Sprintf("avg_score_content: %.2f\n", avgScores[0]))
	buf.WriteString(fmt.Sprintf("avg_score_teaching: %.2f\n", avgScores[1]))
	buf.WriteString(fmt.Sprintf("avg_score_grading: %.2f\n", avgScores[2]))
	buf.WriteString(fmt.Sprintf("avg_score_harvest: %.2f\n", avgScores[3]))
	if len(teacherNames) > 0 {
		buf.WriteString(fmt.Sprintf("teachers: [%s]\n", strings.Join(teacherNames, ", ")))
	}
	buf.WriteString(fmt.Sprintf("last_updated: \"%s\"\n", now))
	buf.WriteString("---\n\n")

	// Title
	buf.WriteString(fmt.Sprintf("# %s (%s)\n\n", course.Name, course.Code))

	// Basic info
	buf.WriteString("## 课程基本信息\n\n")
	buf.WriteString(fmt.Sprintf("- **课程代码**: %s\n", course.Code))
	buf.WriteString(fmt.Sprintf("- **开课学院**: %s\n", course.Institute))
	buf.WriteString(fmt.Sprintf("- **学分**: %d\n", course.Credit))
	if course.CommentCount > 0 {
		overallAvg := (avgScores[0] + avgScores[1] + avgScores[2] + avgScores[3]) / 4.0
		buf.WriteString(fmt.Sprintf("- **综合评分**: %.2f/5.0\n", overallAvg))
		buf.WriteString(fmt.Sprintf("  - 内容: %.2f | 教学: %.2f | 给分: %.2f | 收获: %.2f\n",
			avgScores[0], avgScores[1], avgScores[2], avgScores[3]))
	}
	buf.WriteString(fmt.Sprintf("- **评价数量**: %d\n\n", course.CommentCount))

	// Teachers section
	if len(teacherMap) > 0 {
		buf.WriteString("## 授课教师\n\n")
		for _, t := range teacherMap {
			line := fmt.Sprintf("- **%s**", t.Name)
			if t.Institute != "" {
				line += fmt.Sprintf(" (%s)", t.Institute)
			}
			if t.Job != "" {
				line += fmt.Sprintf(" - %s", t.Job)
			}
			buf.WriteString(line + "\n")
		}
		buf.WriteString("\n")
	}

	// Course groups info
	if len(groups) > 1 {
		buf.WriteString("## 授课组\n\n")
		for _, g := range groups {
			teachers := make([]string, 0, len(g.Teachers))
			for _, t := range g.Teachers {
				teachers = append(teachers, t.Name)
			}
			gAvg := make([]float64, models.ScoreLength)
			if g.CommentCount > 0 {
				for i := 0; i < models.ScoreLength; i++ {
					gAvg[i] = float64(g.Scores[i]) / float64(g.CommentCount)
				}
			}
			buf.WriteString(fmt.Sprintf("- **%s** (教师: %s) - %d条评价",
				g.Code, strings.Join(teachers, ", "), g.CommentCount))
			if g.CommentCount > 0 {
				buf.WriteString(fmt.Sprintf(" | 评分: 内容%.1f 教学%.1f 给分%.1f 收获%.1f",
					gAvg[0], gAvg[1], gAvg[2], gAvg[3]))
			}
			buf.WriteString("\n")
		}
		buf.WriteString("\n")
	}

	// Comments section
	if len(comments) > 0 {
		buf.WriteString("## 课程评价\n\n")
		for i, c := range comments {
			if c.IsCovered {
				continue
			}

			// Find group teachers for this comment
			var groupTeachers []string
			for _, g := range groups {
				if g.ID == c.CourseGroupID {
					for _, t := range g.Teachers {
						groupTeachers = append(groupTeachers, t.Name)
					}
					break
				}
			}

			buf.WriteString(fmt.Sprintf("### 评价 %d: %s\n\n", i+1, c.Title))

			// Metadata line
			meta := fmt.Sprintf("**学期**: %s", FormatSemester(c.Semester))
			if len(groupTeachers) > 0 {
				meta += fmt.Sprintf(" | **授课教师**: %s", strings.Join(groupTeachers, ", "))
			}
			if len(c.Scores) == models.ScoreLength {
				meta += fmt.Sprintf(" | **评分**: %s%d/%s%d/%s%d/%s%d",
					scoreLabels[0], c.Scores[0],
					scoreLabels[1], c.Scores[1],
					scoreLabels[2], c.Scores[2],
					scoreLabels[3], c.Scores[3])
			}
			if c.StudentScoreRanking > 0 {
				meta += fmt.Sprintf(" | **成绩排名**: %s", FormatScoreRanking(c.StudentScoreRanking))
			}
			buf.WriteString(meta + "\n\n")

			// Content
			buf.WriteString(c.Content + "\n\n")

			// Engagement stats
			if c.Like > 0 || c.Dislike > 0 {
				buf.WriteString(fmt.Sprintf("*👍 %d 👎 %d*\n\n", c.Like, c.Dislike))
			}

			buf.WriteString("---\n\n")
		}
	}

	return buf.String()
}

// ExportAllCourses generates markdown for all courses and uploads to S3
func ExportAllCourses(db *gorm.DB) error {
	if !database.IsS3Enabled() {
		return fmt.Errorf("S3 is not enabled, please configure s3 section in config.json")
	}

	ctx := context.Background()
	prefix := database.S3Conf.Prefix

	// Load all courses
	var courses []models.Course
	if err := db.Find(&courses).Error; err != nil {
		return fmt.Errorf("failed to load courses: %w", err)
	}
	syslog.Printf("Found %d courses to export\n", len(courses))

	exported := 0
	skipped := 0
	for _, course := range courses {
		// Load course groups with teachers
		var groups []models.CourseGroup
		if err := db.Preload("Teachers").Where("course_id = ?", course.ID).Find(&groups).Error; err != nil {
			syslog.Printf("Warning: failed to load groups for course %d (%s): %v\n", course.ID, course.Name, err)
			continue
		}

		// Load comments (exclude soft-deleted)
		var comments []models.Comment
		if err := db.Where("course_id = ? AND deleted_at IS NULL", course.ID).
			Order("create_time DESC").
			Find(&comments).Error; err != nil {
			syslog.Printf("Warning: failed to load comments for course %d (%s): %v\n", course.ID, course.Name, err)
			continue
		}

		// Skip courses with no comments (no review data to index)
		if len(comments) == 0 {
			skipped++
			continue
		}

		// Generate markdown
		markdown := GenerateCourseMarkdown(&course, groups, comments)
		filename := GenerateCourseFilename(&course)
		objectKey := path.Join(prefix, filename)

		// Upload to S3
		reader := bytes.NewReader([]byte(markdown))
		err := database.S3Upload(ctx, objectKey, reader, int64(len(markdown)), "text/markdown; charset=utf-8")
		if err != nil {
			syslog.Printf("Warning: failed to upload %s: %v\n", objectKey, err)
			continue
		}

		exported++
		if exported%50 == 0 {
			syslog.Printf("Progress: exported %d courses...\n", exported)
		}
	}

	syslog.Printf("Export complete: %d courses exported, %d skipped (no comments)\n", exported, skipped)
	return nil
}

// ExportSingleCourse generates and uploads markdown for a single course
func ExportSingleCourse(db *gorm.DB, courseID uint) error {
	if !database.IsS3Enabled() {
		return fmt.Errorf("S3 is not enabled")
	}

	ctx := context.Background()
	prefix := database.S3Conf.Prefix

	var course models.Course
	if err := db.First(&course, courseID).Error; err != nil {
		return fmt.Errorf("course %d not found: %w", courseID, err)
	}

	var groups []models.CourseGroup
	if err := db.Preload("Teachers").Where("course_id = ?", courseID).Find(&groups).Error; err != nil {
		return fmt.Errorf("failed to load groups: %w", err)
	}

	var comments []models.Comment
	if err := db.Where("course_id = ? AND deleted_at IS NULL", courseID).
		Order("create_time DESC").
		Find(&comments).Error; err != nil {
		return fmt.Errorf("failed to load comments: %w", err)
	}

	markdown := GenerateCourseMarkdown(&course, groups, comments)
	filename := GenerateCourseFilename(&course)
	objectKey := path.Join(prefix, filename)

	reader := bytes.NewReader([]byte(markdown))
	return database.S3Upload(ctx, objectKey, reader, int64(len(markdown)), "text/markdown; charset=utf-8")
}
