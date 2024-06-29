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

package models

import "time"

type OKResponse struct {
	Error bool        `json:"error" example:"false"`
	Data  interface{} `json:"data,omitempty"`
}

type ErrorResponse struct {
	Error       bool      `json:"error" example:"true"`
	Errno       string    `json:"code"  example:"123"`
	Message     string    `json:"msg"`
	Timestamp   time.Time `json:"timestamp"`
	FullMessage string    `json:"full_msg"`
}
