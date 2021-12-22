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

	// InvalidRequest 等错误:请求完成,但是产生了业务上的错误
	InvalidRequest          = createDescription("InvalidRequest", "请求非法", SILENT)
	UserDoNotExist          = createDescription("UserDoNotExist", "未找到用户", SILENT)
	UserNameDuplicated      = createDescription("UserNameDuplicated", "用户名重复", SILENT)
	StyleDoNotExist         = createDescription("StyleDoNotExist", "风格不存在", SILENT)
	SharerHasNotViewedStyle = createDescription("SharerHasNotViewedStyle", "被分享者没有看过此风格", SILENT)
	HaveViewedStyle         = createDescription("HaveViewedStyle", "此风格已被看过", SILENT)
	BoardDoNotExist         = createDescription("BoardDoNotExist", "未找到画板", SILENT)
	AlreadyCheckedIn        = createDescription("AlreadyCheckedIn", "您已经签过到了", SILENT)
	InvalidJWT              = createDescription("InvalidJWT", "登陆凭据无效或过期", SILENT)

	StyleAttrServiceError = createDescription("StyleAttrServiceError", "服务器内部错误", ERROR)
	UnCaughtError         = createDescription("UnCaughtError", "服务器内部错误", FATAL)
	FailedToGetRedisLock  = createDescription("FailedToGetRedisLock", "服务器繁忙", SILENT)
	StyleDuplicated       = createDescription("StyleDuplicated", "您已经分享过风格", SILENT)
)
