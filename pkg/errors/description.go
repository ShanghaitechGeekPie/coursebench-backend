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
	InvalidRequest      = createDescription("InvalidRequest", "请求非法", SILENT)
	UserDoNotExist      = createDescription("UserDoNotExist", "未找到用户", SILENT)
	UserAlreadyExists   = createDescription("UserAlreadyExists", "用户已存在", SILENT)
	UserEmailDuplicated = createDescription("UserEmailDuplicated", "用户邮箱重复", SILENT)
	InvalidArgument     = createDescription("InvalidArgument", "参数非法", SILENT)

	StyleAttrServiceError = createDescription("StyleAttrServiceError", "服务器内部错误", ERROR)
	UnCaughtError         = createDescription("UnCaughtError", "服务器内部错误", FATAL)
	FailedToGetRedisLock  = createDescription("FailedToGetRedisLock", "服务器繁忙", SILENT)
)
