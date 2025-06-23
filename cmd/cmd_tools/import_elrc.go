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
	"coursebench-backend/pkg/errors"
	"coursebench-backend/pkg/models"
	"coursebench-backend/pkg/queries"
	"encoding/json"
	"fmt"
	"io"
	syslog "log"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

// CourseInfo 课程基本信息
type CourseInfo struct {
	CourseNumber     string   `json:"courseNumber"`
	Name             string   `json:"name_"`
	Teacher          []string `json:"teacher"`       // 工号数组
	TeacherNames     []string `json:"teacher_names"` // 教师姓名列表
	SemesterShowName string   `json:"semester_show_name"`
	SerialNumber     string   `json:"serialNumber"`
	ContentID        string   `json:"contentId_"`
}

// CourseDetailResponse 课程详情API响应
type CourseDetailResponse struct {
	ErrorCode     string `json:"error_code"`
	ErrorMsg      string `json:"error_msg"`
	ExtendMessage struct {
		KczxCourseActivityBk struct {
			ID                 int    `json:"id"`
			CourseName         string `json:"course_name"`
			CourseCode         string `json:"course_code"`
			ShtechCoursenumber string `json:"shtech_coursenumber"`
			SlotStart          string `json:"slot_start"`
			SlotStop           string `json:"slot_stop"`
			TeacherCode        string `json:"teacher_code"`
			College            string `json:"college"`
			Times              string `json:"times"`
			DayOfWeek          int    `json:"day_of_week"`
			Semester           string `json:"semester"`
			Area               string `json:"area"`
			SemesterCode       string `json:"semester_code"`
			TeacherName        string `json:"teacher_name"`
			CollegeName        string `json:"college_name"`
			AreaName           string `json:"area_name"`
		} `json:"KczxCourseActivityBk_instance"`
		JwPkKcxxBk struct {
			ID                 int    `json:"id"`
			CourseName         string `json:"course_name"`
			CourseEnName       string `json:"course_en_name"`
			CourseCode         string `json:"course_code"`
			CourseCategory     string `json:"course_category"`
			CourseCategoryCode string `json:"course_category_code"`
			Credits            string `json:"credits"`
			ClassHour          string `json:"class_hour"`
		} `json:"JwPkKcxxBk_instance"`
	} `json:"extend_message"`
}

// CourseAPIResponse 课程列表API响应
type CourseAPIResponse struct {
	Data struct {
		Results []CourseInfo `json:"results"`

		Page  int64 `json:"page"`
		Size  int64 `json:"size"`
		Total int64 `json:"total"`
	} `json:"data"`
}

// 学期数字到中文映射
type SemesterInfo struct {
	Year      string // 2024-2025
	TermNum   string // 1/2/3
	TermName  string // 秋季/春季/夏季
	FullLabel string // 2024-2025学年秋季
}

func ParseSemesterArg(arg string) (SemesterInfo, error) {
	// 期望格式: 2024-2025-3
	parts := strings.Split(arg, "-")
	if len(parts) != 3 {
		return SemesterInfo{}, fmt.Errorf("invalid semester format, expect 2024-2025-3")
	}
	termMap := map[string]string{"1": "秋季", "2": "春季", "3": "夏季"}
	termName, ok := termMap[parts[2]]
	if !ok {
		return SemesterInfo{}, fmt.Errorf("invalid term number: %s", parts[2])
	}
	return SemesterInfo{
		Year:      parts[0] + "-" + parts[1],
		TermNum:   parts[2],
		TermName:  termName,
		FullLabel: parts[0] + "-" + parts[1] + "学年" + termName,
	}, nil
}

func ImportELRCWithSemester(semesterArg string) {
	semInfo, err := ParseSemesterArg(semesterArg)
	if err != nil {
		syslog.Fatalf("Semester argument error: %v", err)
	}
	syslog.Printf("Parsed semester: %s, term: %s, label: %s", semInfo.Year, semInfo.TermNum, semInfo.FullLabel)
	syslog.Printf("Please confirm: year=%s, termNum=%s, termName=%s, fullLabel=%s", semInfo.Year, semInfo.TermNum, semInfo.TermName, semInfo.FullLabel)
	fmt.Print("Type 'yes' to continue: ")
	var confirm string
	fmt.Scanln(&confirm)
	if strings.ToLower(confirm) != "yes" {
		syslog.Fatalf("User did not confirm. Abort.")
	}

	db := database.GetDB()

	courses, pageCourseMap, totalPages, err := fetchAllCoursesWithPageInfoAndSemester(semInfo)
	if err != nil {
		panic(err)
	}

	syslog.Printf("Fetched %d courses from API\n", len(courses))

	err = db.Transaction(func(tx *gorm.DB) error {
		for _, courseInfo := range courses {
			if err := processCourseWithSemester(tx, courseInfo, semInfo); err != nil {
				syslog.Printf("Error processing course %s: %v\n", courseInfo.SerialNumber, err)
				return err
			}
		}
		return nil
	})

	if err != nil {
		panic(err)
	}

	syslog.Printf("Finished importing courses from ELRC API\n")

	syslog.Printf("Summary: Total pages processed: %d\n", totalPages)
	for page, courseList := range pageCourseMap {
		syslog.Printf("Page %d: %d courses\n", page, len(courseList))
		for _, c := range courseList {
			syslog.Printf("  - %s: %s\n", c.CourseNumber, c.Name)
		}
	}
}

// fetchAllCoursesWithPageInfoAndSemester 获取所有课程数据，支持学期参数
func fetchAllCoursesWithPageInfoAndSemester(semInfo SemesterInfo) ([]CourseInfo, map[int][]CourseInfo, int, error) {
	baseURL := "https://elrc.shanghaitech.edu.cn/learn/shanghai/tech/get/course"

	var allCourses []CourseInfo
	pageCourseMap := make(map[int][]CourseInfo)
	page := 1
	pageSize := 20  // 初始默认
	totalPages := 1 // 初始默认
	var total int64 = 0

	for {
		syslog.Printf("Fetching page %d...\n", page)

		params := url.Values{}
		params.Add("page", strconv.Itoa(page))
		params.Add("size", strconv.Itoa(pageSize))
		params.Add("courseType", "2")
		params.Add("semester", semInfo.FullLabel)

		fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

		resp, err := http.Get(fullURL)
		if err != nil {
			syslog.Printf("Error fetching page %d: %v\n", page, err)
			break
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			syslog.Printf("Bad response status for page %d: %d\n", page, resp.StatusCode)
			break
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			syslog.Printf("Error reading response body for page %d: %v\n", page, err)
			break
		}

		var apiResponse CourseAPIResponse
		if err := json.Unmarshal(body, &apiResponse); err != nil {
			syslog.Printf("Error parsing JSON for page %d: %v\n", page, err)
			break
		}

		if len(apiResponse.Data.Results) == 0 {
			syslog.Printf("No results found on page %d\n", page)
			break
		}

		allCourses = append(allCourses, apiResponse.Data.Results...)
		pageCourseMap[page] = append([]CourseInfo{}, apiResponse.Data.Results...)

		if page == 1 {
			pageSize = int(apiResponse.Data.Size)
			total = apiResponse.Data.Total
			if pageSize > 0 {
				totalPages = int((total + int64(pageSize) - 1) / int64(pageSize))
			} else {
				totalPages = 1
			}
			syslog.Printf("API reports total=%d, pageSize=%d, totalPages=%d\n", total, pageSize, totalPages)
		}

		if page >= totalPages {
			break
		}
		page++
		time.Sleep(100 * time.Millisecond)
	}

	return allCourses, pageCourseMap, totalPages, nil
}

// fetchCourseDetailWithSemester 获取课程详细信息，支持学期参数
func fetchCourseDetailWithSemester(serialNumber string, semInfo SemesterInfo) (*CourseDetailResponse, error) {
	baseURL := "https://elrc.shanghaitech.edu.cn/shanghaitechdatasync/datasync/bksCourse/"
	params := url.Values{}
	params.Add("semester", semInfo.Year)
	params.Add("term", semInfo.TermNum)
	params.Add("course_no", serialNumber)
	params.Add("course_id", "undefined")

	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	resp, err := http.Get(fullURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad response status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var detailResponse CourseDetailResponse
	if err := json.Unmarshal(body, &detailResponse); err != nil {
		return nil, err
	}

	if detailResponse.ErrorCode != "shanghaitech.0000.0000" {
		return nil, fmt.Errorf("API error: %s - %s", detailResponse.ErrorCode, detailResponse.ErrorMsg)
	}

	return &detailResponse, nil
}

// processCourseWithSemester 处理单个课程，支持学期参数
func processCourseWithSemester(tx *gorm.DB, courseInfo CourseInfo, semInfo SemesterInfo) error {
	detail, err := fetchCourseDetailWithSemester(courseInfo.SerialNumber, semInfo)
	if err != nil {
		syslog.Printf("Warning: could not fetch detail for course %s: %v\n", courseInfo.SerialNumber, err)
		detail = nil
	}
	// 解析学分
	var credit int = 0
	if detail != nil && detail.ExtendMessage.JwPkKcxxBk.Credits != "" {
		if creditF, err := strconv.ParseFloat(detail.ExtendMessage.JwPkKcxxBk.Credits, 64); err == nil {
			credit = int(creditF)
		}
	}
	// 获取开课单位
	var institute string = ""
	if detail != nil {
		institute = detail.ExtendMessage.KczxCourseActivityBk.CollegeName
	}
	if institute == "" {
		institute = "未知单位"
	}
	// 查找或创建课程
	course := &models.Course{}
	err = tx.Where("code = ?", courseInfo.CourseNumber).Take(course).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if err != nil {
		syslog.Printf("Add course %s %s %s %d\n", courseInfo.CourseNumber, courseInfo.Name, institute, credit)
		course, err = queries.AddCourse(tx, courseInfo.Name, institute, credit, courseInfo.CourseNumber)
		if err != nil {
			return err
		}
		other_teachers := []int{models.TEACHER_OTHER_ID}
		_, err = queries.AddCourseGroup(tx, "", int(course.ID), other_teachers)
		if err != nil {
			return errors.Wrap(err, errors.DatabaseError)
		}
	} else {
		syslog.Printf("Find course %s %d\n", courseInfo.CourseNumber, course.ID)
	}
	// 检查课程组是否已存在 (基于semester+serialNumber的唯一性)
	groupCode := fmt.Sprintf("%s-%s", semInfo.FullLabel, courseInfo.SerialNumber)
	existingGroup := &models.CourseGroup{}
	err = tx.Where("code = ? AND course_id = ?", groupCode, course.ID).Take(existingGroup).Error
	if err == nil {
		syslog.Printf("Course group already exists: %s\n", groupCode)
		return nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	// 处理教师信息
	var teacherIDs []int
	teachers := make(map[string]string)
	for i, uniIDStr := range courseInfo.Teacher {
		if i < len(courseInfo.TeacherNames) {
			teachers[uniIDStr] = courseInfo.TeacherNames[i]
		}
	}
	for uniIDStr, teacherName := range teachers {
		uniID, err := strconv.Atoi(uniIDStr)
		if err != nil {
			syslog.Printf("Warning: invalid UniID '%s' for teacher '%s', skipping\n", uniIDStr, teacherName)
			continue
		}
		teacher, err := findOrCreateTeacherByUniID(tx, uniID, teacherName, institute)
		if err != nil {
			syslog.Printf("Error handling teacher %s (UniID: %d): %v\n", teacherName, uniID, err)
			continue
		}
		teacherIDs = append(teacherIDs, int(teacher.ID))
	}
	if len(teacherIDs) == 0 {
		teacherIDs = []int{models.TEACHER_OTHER_ID}
	}
	var data []struct {
		CourseGroupID int
		ArrayAgg      []int
	}
	tx.Raw(`select course_groups.id as course_group_id, array_agg(coursegroup_teachers.teacher_id) as array_agg from course_groups
		inner join coursegroup_teachers on course_groups.id = coursegroup_teachers.course_group_id
		where course_groups.course_id=? group by course_groups.id;`, course.ID).Scan(&data)
	sort.Ints(teacherIDs)
	for _, c := range data {
		tids := append([]int(nil), c.ArrayAgg...)
		sort.Ints(tids)
		if len(tids) == len(teacherIDs) {
			match := true
			for i := range tids {
				if tids[i] != teacherIDs[i] {
					match = false
					break
				}
			}
			if match {
				syslog.Printf("Skip adding course group: identical teacher set already exists (group id: %d)\n", c.CourseGroupID)
				return nil
			}
		}
	}
	group, err := queries.AddCourseGroup(tx, "", int(course.ID), teacherIDs) // TODO: 这里的 code 省略了，因为会重复（学期/班级不同但是老师配置相同）
	if err != nil {
		return err
	}
	syslog.Printf("Add course group %s %d %d (semester: %s, serialNumber: %s)\n",
		courseInfo.CourseNumber, course.ID, group.ID, semInfo.FullLabel, courseInfo.SerialNumber)
	time.Sleep(50 * time.Millisecond)
	return nil
}

// findOrCreateTeacherByUniID 根据UniID查找或创建教师
func findOrCreateTeacherByUniID(tx *gorm.DB, uniID int, name string, institute string) (*models.Teacher, error) {
	teacher := &models.Teacher{}

	// 首先根据UniID查找
	err := tx.Where("uni_id = ?", uniID).Take(teacher).Error
	if err == nil {
		// 找到了，更新姓名和学院信息（如果为空的话）
		updateNeeded := false
		if teacher.Name == "" && name != "" {
			teacher.Name = name
			updateNeeded = true
		}
		if teacher.Institute == "" && institute != "" {
			teacher.Institute = institute
			updateNeeded = true
		}
		if updateNeeded {
			if err := tx.Save(teacher).Error; err != nil {
				return nil, err
			}
			syslog.Printf("Updated teacher info for UniID %d: %s\n", uniID, name)
		}
		return teacher, nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// UniID不存在，根据姓名查找
	err = tx.Where("name = ?", name).Take(teacher).Error
	if err == nil {
		// 找到同名教师，更新UniID
		teacher.UniID = uniID
		if teacher.Institute == "" && institute != "" {
			teacher.Institute = institute
		}
		if err := tx.Save(teacher).Error; err != nil {
			return nil, err
		}
		syslog.Printf("Updated UniID for existing teacher %s: %d\n", name, uniID)
		return teacher, nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// 教师不存在，创建新教师
	teacher = &models.Teacher{
		UniID:        uniID,
		Name:         name,
		Institute:    institute,
		Job:          "",
		Introduction: "",
		EamsID:       0, // Deprecated
	}

	result := tx.Create(teacher)
	if result.Error != nil {
		return nil, errors.Wrap(result.Error, errors.DatabaseError)
	}

	syslog.Printf("Created new teacher %s (UniID: %d)\n", name, uniID)
	return teacher, nil
}
