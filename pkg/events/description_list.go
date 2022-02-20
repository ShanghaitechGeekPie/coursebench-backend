package events

import "go.uber.org/zap/zapcore"

var (
	FATAL = zapcore.DPanicLevel
	ERROR = zapcore.ErrorLevel
	WARN  = zapcore.WarnLevel
	INFO  = zapcore.InfoLevel
)

var (
	// ServerPanics :造成程序崩溃的错误
	ServerPanics = createDescription("ServerPanics", "服务器内部错误", FATAL, 500)
	LogicError   = createDescription("LogicError", "服务器内部错误", FATAL, 500)
	// InternalServerError 等错误:造成请求无法完成的错误
	InternalServerError = createDescription("InternalServerError", "服务器内部错误", ERROR, 500)
	InfluxDBError       = createDescription("InfluxDBError", "服务器内部错误", ERROR, 500)
	NebulaError         = createDescription("NebulaError", "服务器内部错误", ERROR, 500)
	RedisError          = createDescription("RedisError", "服务器内部错误", ERROR, 500)
	DatabaseError       = createDescription("DatabaseError", "服务器内部错误", ERROR, 500)
	MinIOError          = createDescription("MinIOError", "服务器内部错误", ERROR, 500)
	GobEncodingError    = createDescription("GobEncodingError", "服务器内部错误", ERROR, 500)
	GobDecodingError    = createDescription("GobDecodingError", "服务器内部错误", ERROR, 500)
	CacheMiss           = createDescription("CacheMiss", "服务器繁忙", WARN, 500)

	// InvalidRequest 等错误:请求完成,但是产生了业务上的错误
	InvalidRequest        = createDescription("InvalidRequest", "请求非法", INFO, 400)
	UserDoNotExist        = createDescription("UserDoNotExist", "未找到用户", INFO, 400)
	UserAlreadyExists     = createDescription("UserAlreadyExists", "用户已存在", INFO, 400)
	UserEmailDuplicated   = createDescription("UserEmailDuplicated", "用户邮箱重复", INFO, 400)
	UserPasswordIncorrect = createDescription("UserPasswordIncorrect", "用户密码错误", INFO, 400)
	UserNotLogin          = createDescription("UserNotLogin", "用户未登录", INFO, 400)

	TeacherNotExists = createDescription("TeacherNotExists", "未找到教师", INFO, 400)

	CourseNotExists      = createDescription("CourseNotExists", "未找到课程", INFO, 400)
	CourseGroupNotExists = createDescription("CourseGroupNotExists", "未找到课程授课组", INFO, 400)

	CommentAlreadyExists = createDescription("CommentAlreadyExists", "评论已存在", INFO, 400)
	CommentNotExists     = createDescription("CommentNotExists", "评论不存在", INFO, 400)

	InvalidArgument = createDescription("InvalidArgument", "参数非法", INFO, 400)

	PermissionDenied = createDescription("PermissionDenied", "您没有权限执行此操作", INFO, 403)

	UnCaughtError        = createDescription("UnCaughtError", "服务器内部错误", FATAL, 500)
	FailedToGetRedisLock = createDescription("FailedToGetRedisLock", "服务器繁忙", WARN, 500)

	SMTPError = createDescription("SMTPError", "邮件发送失败", ERROR, 500)
)
