package debugger

import (
	"github.com/fansqz/ShowMeCode/go-debugger/constants"
)

// StartOption 启动调试的参数
type StartOption struct {
	ExecFile string
	// Main文件代码
	MainCode string
	// Callback 事件回调
	Callback NotificationCallback
}

// Breakpoint 表示断点
type Breakpoint struct {
	File string // 文件名称
	Line int    // 行号
}

func NewBreakpoint(file string, line int) *Breakpoint {
	return &Breakpoint{file, line}
}

// StackFrame 栈帧
type StackFrame struct {
	ID   string `json:"id"`   // 栈帧id
	Name string `json:"name"` // 函数名称
	Path string `json:"path"` // 文件路径
	Line int    `json:"line"`
}

// Scope 作用域
type Scope struct {
	Name      constants.ScopeName
	Reference string // 作用域的引用
}

// Variable 变量
type Variable struct {
	Name  string  `json:"name"`
	Type  string  `json:"type"`
	Value *string `json:"value"`
	// 变量引用
	Reference string `json:"reference"`
	// childrenNumber
	ChildrenNumber int
}

// 定义的一些Event
var (
	CompileSuccessEvent = NewCompileEvent(true, "用户代码编译成功")
	LaunchSuccessEvent  = NewLaunchEvent(true, "目标代码加载成功")
	LaunchFailEvent     = NewLaunchEvent(false, "目标代码加载失败")
)

// BreakpointEvent 断点事件
// 该event指示有关断点的某些信息已更改。
type BreakpointEvent struct {
	Reason      constants.BreakpointReasonType
	Breakpoints []*Breakpoint
}

func NewBreakpointEvent(reason constants.BreakpointReasonType, breakpoints []*Breakpoint) *BreakpointEvent {
	return &BreakpointEvent{
		Reason:      reason,
		Breakpoints: breakpoints,
	}
}

// OutputEvent
// 用户程序输出
type OutputEvent struct {
	Output string // 输出内容
}

func NewOutputEvent(output string) *OutputEvent {
	return &OutputEvent{
		Output: output,
	}
}

// StoppedEvent
// 该event表明，由于某些原因，被调试进程的执行已经停止。
// 这可能是由先前设置的断点、完成的步进请求、执行调试器语句等引起的。
type StoppedEvent struct {
	Reason constants.StoppedReasonType // 停止执行的原因
	File   string                      // 当前停止在哪个文件
	Line   int                         // 停止在某行
}

func NewStoppedEvent(reason constants.StoppedReasonType, file string, line int) *StoppedEvent {
	return &StoppedEvent{
		Reason: reason,
		File:   file,
		Line:   line,
	}
}

// ContinuedEvent
// 该event表明debug的执行已经继续。
// 请注意:debug adapter不期望发送此事件来响应暗示执行继续的请求，例如启动或继续。
// 它只有在没有先前的request暗示这一点时，才有必要发送一个持续的事件。
type ContinuedEvent struct {
}

func NewContinuedEvent() *ContinuedEvent {
	return &ContinuedEvent{}
}

// ExitedEvent
// 该event表明被调试对象已经退出并返回exit code。但是并不意味着调试会话结束
type ExitedEvent struct {
	ExitCode int
	Message  string
}

func NewExitedEvent(code int, message string) *ExitedEvent {
	return &ExitedEvent{
		ExitCode: code,
		Message:  message,
	}
}

// CompileEvent
// 编译事件
type CompileEvent struct {
	Success bool
	Message string // 编译产生的信息
}

func NewCompileEvent(success bool, message string) *CompileEvent {
	return &CompileEvent{
		Success: success,
		Message: message,
	}
}

// LaunchEvent
// 调试资源准备成功
type LaunchEvent struct {
	Success bool
	Message string // 启动gdb的消息
}

func NewLaunchEvent(success bool, message string) *LaunchEvent {
	return &LaunchEvent{
		Success: success,
		Message: message,
	}
}
