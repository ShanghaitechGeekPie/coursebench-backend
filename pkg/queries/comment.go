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

package queries

import (
	"coursebench-backend/pkg/models"
	"time"
)

func CheckCommentTitle(title string) bool {
	length := len(title)
	if length == 0 || length > 200 {
		return false
	}
	return true
}

func CheckCommentContent(content string) bool {
	length := len(content)
	if length == 0 || length > 50000 {
		return false
	}
	return true
}

func CheckSemester(semester int) bool {
	if semester < 0 || semester > 1000000 {
		return false
	}
	year := semester / 100
	sem := semester % 100
	if year < 2014 || year > time.Now().Year() {
		return false
	}
	if sem < 1 || sem > 3 {
		return false
	}
	return true
}

func CheckCommentScore(scores []int64) bool {
	if len(scores) != models.ScoreLength {
		return false
	}
	for _, v := range scores {
		if v < 1 || v > 5 {
			return false
		}
	}
	return true
}

func CheckCommentScoreRanking(ranking int) bool {
	if ranking < 1 || ranking > 11 {
		return false
	}
	return true
}
