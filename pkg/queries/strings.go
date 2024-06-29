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

// 返回一个字符串的渲染长度
// 非常粗糙的方案，一个ascii字符的长度为1，一个中文字符的长度为2
func GetActualLength(s string) int {
	sum := 0
	for _, c := range s {
		if c > 0x80 {
			sum += 2
		} else {
			sum += 1
		}
	}
	return sum
}
