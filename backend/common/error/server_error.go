package error

// 服务器 1000-9999 服务器错误
const (
	CodeUnknownError      = 1000 + iota
	CodeMysqlError        // mysql错误
	CodeServerErr         // 服务器错误
	CodeRedisErr          // redis错误
	CodeServerBusyErr     // 服务器繁忙
	CodeServerMaintenance // 服务器维护中
)

var (
	ErrUnknown           = NewError(CodeUnknownError, "服务器未知错误", ErrTypeServer)
	ErrMysql             = NewError(CodeMysqlError, "服务器数据错误", ErrTypeServer)
	ErrServer            = NewError(CodeServerErr, "服务器错误", ErrTypeServer)
	ErrRedis             = NewError(CodeRedisErr, "服务器数据错误", ErrTypeServer)
	ErrServerBusy        = NewError(CodeServerBusyErr, "服务器繁忙", ErrTypeServer)
	ErrServerMaintenance = NewError(CodeServerMaintenance, "服务器维护中", ErrTypeServer)
)
