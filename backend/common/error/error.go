package error

import (
	"net/http"
)

const (
	CodeRecordNotFound = 10000 + iota // 记录不存在
	CodeParamLen                      // 参数长度不一致
	CodeCustomMsg                     // 自定义错误消息
)

// error 对应的 http 请求
const (
	ErrTypeBus        = http.StatusOK                  // 业务错误
	ErrTypeBadReq     = http.StatusBadRequest          // 请求错误
	ErrTypeAuth       = http.StatusUnauthorized        // 未授权操作
	ErrTypePermission = http.StatusForbidden           // 权限不足
	ErrTypeServer     = http.StatusInternalServerError // 服务器错误
	ErrTypeNotFound   = http.StatusNotFound            // 记录不存在
	ErrTooManyReq     = http.StatusTooManyRequests
)

type Error struct {
	Code     int
	Message  string
	HttpCode int
}

func (e *Error) Error() string {
	return e.Message
}

// NewError 封装错误使用
func NewError(code int, msg string, httpCode int) *Error {
	return &Error{
		Code:     code,
		Message:  msg,
		HttpCode: httpCode,
	}
}

// NewCustomMsg 自定义错误消息，直接传入msg
func NewCustomMsg(msg string) *Error {
	return &Error{
		Code:     CodeCustomMsg,
		Message:  msg,
		HttpCode: ErrTypeBus,
	}
}

// NewRecordNotFoundErr 记录不存在
func NewRecordNotFoundErr(msg string) *Error {
	return &Error{
		Code:     CodeRecordNotFound,
		Message:  msg,
		HttpCode: ErrTypeBus,
	}
}

// NewParamErr 参数错误
func NewParamErr(msg string) *Error {
	return &Error{
		Code:     CodeBadRequest,
		Message:  msg,
		HttpCode: ErrTypeBadReq,
	}
}
