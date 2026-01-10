package error

// 服务器 10000 以上业务错误 服务器错误, 不同模块业务错误码间隔 500
// 前10000 留给公用的错误

/************User错误**************/
const (
	CodeUserNameOrPasswordWrong       = 11000 + iota // 用户名或密码错误
	CodeUserEmailIsExist                             // 邮箱已存在
	CodeLoginCodeWrong                               // 登录验证码错误
	CodeUserPasswordNotEnoughAccuracy                // 用户密码精度不够
	CodePasswordEncodeFailed                         // 密码加密失败
	CodeUserNotExist                                 // 用户不存在
	CodeUserUnknownError                             // 用户服务未知错误
	CodeEmailFormatWrong                             // 邮箱格式有问题
	CodeEmailNotRegister                             // 邮箱不存在
)

var (
	ErrUserPasswordNotEnoughAccuracy = NewError(CodeUserPasswordNotEnoughAccuracy, "The password is not accurate enough", ErrTypeBus)
	ErrPasswordEncodeFailed          = NewError(CodePasswordEncodeFailed, "Failed to encode password", ErrTypeServer)
	ErrUserNotExist                  = NewError(CodeUserNotExist, "The user does not exist", ErrTypeBus)
	ErrUserUnknownError              = NewError(CodeUserUnknownError, "Unknown error", ErrTypeServer)
	ErrUserNameOrPasswordWrong       = NewError(CodeUserNameOrPasswordWrong, "账号或密码错误", ErrTypeBus)
	ErrUserEmailIsExist              = NewError(CodeUserEmailIsExist, "email already exist", ErrTypeBus)
	ErrLoginCodeWrong                = NewError(CodeLoginCodeWrong, "登录验证码错误", ErrTypeBus)
	ErrEmailFormatWrong              = NewError(CodeEmailFormatWrong, "邮箱格式异常", ErrTypeBus)
	ErrEmailNotRegister              = NewError(CodeEmailNotRegister, "邮箱未注册，请检查", ErrTypeBus)
)

/************Question错误**************/
const (
	CodeProblemCodeIsExist     = 11500 + iota //题目编号已存在
	CodeProblemCodeCheckFailed                // 题目编号检测失败
	CodeProblemGetFailed                      // 获取题目失败
	CodeProblemInsertFailed                   // 添加题目失败
	CodeProblemUpdateFailed                   // 题目更新失败
	CodeProblemListFailed                     // 获取题目列表失败
	CodeProblemNotExist                       // 题目不存在
)

var (
	ErrProblemCodeIsExist     = NewError(CodeProblemCodeIsExist, "problem code is exist", ErrTypeBus)
	ErrProblemCodeCheckFailed = NewError(CodeProblemCodeCheckFailed, "The problem code check failed", ErrTypeServer)
	ErrProblemGetFailed       = NewError(CodeProblemGetFailed, "The problem get failed", ErrTypeServer)
	ErrProblemInsertFailed    = NewError(CodeProblemInsertFailed, "The problem insert failed", ErrTypeServer)
	ErrProblemUpdateFailed    = NewError(CodeProblemUpdateFailed, "The problem update failed", ErrTypeServer)
	ErrProblemListFailed      = NewError(CodeProblemListFailed, "Failed to get the problem list", ErrTypeServer)
	ErrProblemNotExist        = NewError(CodeProblemNotExist, "The problem does not exist", ErrTypeBus)
)

/************judge错误**************/
const (
	CodeSubmitFailed = 12500 + iota //题目编号已存在
	CodeExecuteFailed
	CodeCompileFailed
	CodeLanguageNotSupported
	CodeDebuggerIsClosed
	CodeProgramIsRunningOptionFail
	CodeErrDebugIsFinish
)

var (
	ErrSubmitFailed               = NewError(CodeSubmitFailed, "Submit error", ErrTypeBus)
	ErrExecuteFailed              = NewError(CodeExecuteFailed, "Execute error", ErrTypeBus)
	ErrCompileFailed              = NewError(CodeCompileFailed, "Compilation error", ErrTypeBus)
	ErrLanguageNotSupported       = NewError(CodeLanguageNotSupported, "This language is not supported", ErrTypeBus)
	ErrDebuggerIsClosed           = NewError(CodeDebuggerIsClosed, "debug is closed", ErrTypeBus)
	ErrProgramIsRunningOptionFail = NewError(CodeProgramIsRunningOptionFail, "The program is running", ErrTypeBus)
)
