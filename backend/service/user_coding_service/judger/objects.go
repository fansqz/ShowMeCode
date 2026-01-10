package judger

import "github.com/fansqz/fancode-backend/constants"

// ExecuteOptions 执行文件可选操作
type ExecuteOptions struct {
	Language        constants.LanguageType
	LimitTime       int64 // 资源限制
	MemoryLimit     int64
	CPUQuota        int64
	ExcludedPaths   []string // 屏蔽的敏感路径
	ReplacementPath string   // 取代敏感路径的路径
}

// ExecuteResult 程序执行结果
type ExecuteResult struct {
	Executed     bool   // 判题是否执行成功
	ErrorMessage string // 异常信息
	Output       []byte // 输出结果（正常输出结果，如果有）
	UsedTime     int64  // 执行时间（以纳秒为单位）
	UsedMemory   int64  // 内存使用量（以字节为单位）
	UsedCpuTime  int64  // cpu使用
}

// CompileOptions 编译文件可选参数
type CompileOptions struct {
	Language        constants.LanguageType
	LimitTime       int64
	ExcludedPaths   []string // 屏蔽的敏感路径
	ReplacementPath string   // 取代敏感路径的路径
}

// CompileResult 系统编译结果
type CompileResult struct {
	Compiled         bool   // 判题是否编译成功
	ErrorMessage     string // 异常信息
	CompiledFilePath string // 输出文件路径
}
