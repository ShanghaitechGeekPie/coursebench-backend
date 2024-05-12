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

package errors

var (
	FATAL  = NewOptionBuilder().SetLogLevel(Fatal).Build()
	ERROR  = NewOptionBuilder().SetLogLevel(Error).Build()
	SILENT = NewOptionBuilder().SetLogLevel(Error).Build() // TODO now we treat slient as errors
)

var (
	// ServerPanics :造成程序崩溃的错误
	ServerPanics = createDescription("ServerPanics", "服务器内部错误", FATAL)
	LogicError   = createDescription("LogicError", "服务器内部错误", FATAL)
	// InternalServerError 等错误:造成请求无法完成的错误
	InternalServerError = createDescription("InternalServerError", "服务器内部错误", ERROR)
	InfluxDBError       = createDescription("InfluxDBError", "服务器内部错误", ERROR)
	NebulaError         = createDescription("NebulaError", "服务器内部错误", ERROR)
	RedisError          = createDescription("RedisError", "服务器内部错误", ERROR)
	DatabaseError       = createDescription("DatabaseError", "服务器内部错误", ERROR)
	MinIOError          = createDescription("MinIOError", "服务器内部错误", ERROR)
	GobEncodingError    = createDescription("GobEncodingError", "服务器内部错误", ERROR)
	GobDecodingError    = createDescription("GobDecodingError", "服务器内部错误", ERROR)
	GPTWorkerError      = createDescription("GPTWorkerError", "GPT Worker 出错", ERROR)

	// InvalidRequest 等错误:请求完成,但是产生了业务上的错误
	InvalidRequest        = createDescription("InvalidRequest", "请求非法", SILENT, 400)
	UserNotExists         = createDescription("UserNotExists", "未找到用户", SILENT, 400)
	UserAlreadyExists     = createDescription("UserAlreadyExists", "用户已存在", SILENT, 400)
	UserEmailDuplicated   = createDescription("UserEmailDuplicated", "用户邮箱重复", SILENT, 400)
	UserPasswordIncorrect = createDescription("UserPasswordIncorrect", "用户密码错误", SILENT, 400)
	UserNotLogin          = createDescription("UserNotLogin", "用户未登录", SILENT, 400)
	UserNotActive         = createDescription("UserNotActive", "用户邮箱未激活", SILENT, 400)
	MailCodeInvalid       = createDescription("MailCodeInvalid", "邮箱验证码错误", SILENT, 400)
	CaptchaMismatch       = createDescription("CaptchaMismatch", "验证码错误", SILENT, 400)
	NoCaptchaToken        = createDescription("NoCaptchaToken", "未请求过验证码Token，请检查您的 Cookie 设置", SILENT, 400)
	CaptchaExpired        = createDescription("CaptchaExpired", "验证码已过期", SILENT, 400)

	TeacherNotExists = createDescription("TeacherNotExists", "未找到教师", SILENT, 400)

	CourseNotExists      = createDescription("CourseNotExists", "未找到课程", SILENT, 400)
	CourseGroupNotExists = createDescription("CourseGroupNotExists", "未找到课程授课组", SILENT, 400)

	CommentAlreadyExists = createDescription("CommentAlreadyExists", "评论已存在", SILENT, 400)
	CommentNotExists     = createDescription("CommentNotExists", "评论不存在", SILENT, 400)

	FileTooLarge = createDescription("FileTooLarge", "文件过大", SILENT, 400)

	InvalidArgument = createDescription("InvalidArgument", "参数非法", SILENT, 400)

	PermissionDenied = createDescription("PermissionDenied", "您没有权限执行此操作", SILENT, 403)

	UnCaughtError        = createDescription("UnCaughtError", "服务器内部错误", FATAL, 500)
	FailedToGetRedisLock = createDescription("FailedToGetRedisLock", "服务器繁忙", SILENT, 500)

	SMTPError = createDescription("SMTPError", "邮件发送失败", FATAL, 500)
)
