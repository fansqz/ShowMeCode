package constants

const (
	// Accepted 成功 在提交时使用
	Accepted = 200 + iota
	// RunSuccess 运行成功 在执行时使用
	RunSuccess
	// WrongAnswer 答案错误
	WrongAnswer
	// CompileError 编译出错
	CompileError
	// RuntimeError 运行出错
	RuntimeError
)
