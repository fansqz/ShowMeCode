package error

const (
	CodeBadRequest = 3000 + iota
	CodeSessionExpire
	CodeSessionInvalid
	CodeSessionNotEqual
	CodePermissionInvalid
	CodeLoginTypeErr
	CodeTooManyRequestsErr
)

// 如果参数错误，不做准确的提示可以用这个
var (
	ErrBadRequest        = NewError(CodeBadRequest, "参数错误", ErrTypeBadReq)
	ErrSessionExpire     = NewError(CodeSessionExpire, "登入已过期", ErrTypeAuth)
	ErrSessionInvalid    = NewError(CodeSessionInvalid, "账号未登入，请先登入", ErrTypeAuth)
	ErrSessionNotEqual   = NewError(CodeSessionNotEqual, "账号已在其他地方登入", ErrTypeAuth)
	ErrPermissionInvalid = NewError(CodePermissionInvalid, "权限不足", ErrTypePermission)
	ErrLoginType         = NewError(CodeLoginTypeErr, "登入类型错误", ErrTypeBadReq)
	ErrTooManyRequests   = NewError(CodeTooManyRequestsErr, "服务器忙碌中，请稍后再试", ErrTooManyReq)
)
