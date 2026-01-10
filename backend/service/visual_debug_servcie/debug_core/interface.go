package debug_core

import (
	"context"
	"time"

	"github.com/fansqz/fancode-backend/constants"
)

type NotificationCallback func(interface{})

// Debugger
// 用户的一次调试过程处理
// debugger目前设置为支持多文件的
// 需要保证并发安全
type Debugger interface {
	// Start
	// 开始调试，及调用start命令，callback用来异步处理用户程序输出
	Start(ctx context.Context, option *Option) error
	// Send 输入
	Send(ctx context.Context, input string) error
	// StepOver 下一步，不会进入函数内部
	StepOver(ctx context.Context) error
	// StepIn 下一步，会进入函数内部
	StepIn(ctx context.Context) error
	// StepOut 单步退出
	StepOut(ctx context.Context) error
	// Continue 忽略继续执行
	Continue(ctx context.Context) error
	// SetBreakpoints 添加断点
	// 设置断点
	SetBreakpoints(ctx context.Context, breakpoints []int) ([]*Breakpoint, error)
	// GetStackTrace 获取栈帧
	GetStackTrace(ctx context.Context) ([]*StackFrame, error)
	// GetFrameVariables 获取某个栈帧中的变量列表
	GetFrameVariables(ctx context.Context, frameId int) ([]*Variable, error)
	// GetVariables 查看引用的值
	GetVariables(ctx context.Context, reference int) ([]*Variable, error)
	// Terminate 终止调试
	// 调用完该命令以后可以重新Launch
	Terminate(ctx context.Context) error
	// StructVisual 以结构体为导向的可视化方法，一般用作树、图的可视化
	StructVisual(ctx context.Context, query *StructVisualQuery) (*StructVisualData, error)
	// ArrayVisual 数组可视化
	ArrayVisual(ctx context.Context, query *ArrayVisualQuery) (*ArrayVisualData, error)
	// Array2DVisual 二维数组可视化
	Array2DVisual(ctx context.Context, query *Array2DVisualQuery) (*Array2DVisualData, error)
}

// Option 启动调试的参数
type Option struct {
	// Language 程序的编程语言
	Language constants.LanguageType

	// 调试器镜像名称
	DebuggerImage string

	// Code 用户代码
	Code string
	// BreakPoints 初始化的断点
	BreakPoints []int

	// CompileTimeout 编译超时
	CompileTimeout time.Duration
	// OptionTimeout 操作超时，调试器经常需要和编辑器通信，在OptionLimitTime的时间之内没有响应就取消等待
	OptionTimeout time.Duration
	// DebugTimeout 调试超时时间，如果调试器经过DebugTimeout没有返回信息，那么可以认为调试器超时，将调试终止
	// 可能是因为用户写的程序有死循环，也可能是因为用户太久没有进行操作。
	DebugTimeout time.Duration
	// MemoryLimit 内存限制
	MemoryLimit int64
	// CPUQuota cpu配额
	CPUQuota int64

	// Callback 事件回调
	Callback NotificationCallback

	// TempDir 临时文件目录
	TempDir string
}
