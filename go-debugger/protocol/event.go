package protocol

import "github.com/fansqz/go-debugger/constants"

// CompileEvent
// 编译事件
type CompileEvent struct {
	Event   constants.DebugEventType `json:"event"`
	Success bool                     `json:"success"`
	Message string                   `json:"message"` // 编译产生的信息
}

// LaunchEvent
// 加载用户代码事件
type LaunchEvent struct {
	Event   constants.DebugEventType `json:"event"`
	Success bool                     `json:"success"`
	Message string                   `json:"message"` // 启动gdb的消息
}

// BreakpointEvent 断点事件
// 该event指示有关断点的某些信息已更改。
type BreakpointEvent struct {
	Event       constants.DebugEventType       `json:"event"`
	Reason      constants.BreakpointReasonType `json:"reason"`
	Breakpoints []int                          `json:"breakpoints"`
}

// OutputEvent
// 该事件表明目标已经产生了一些输出。
type OutputEvent struct {
	Event  constants.DebugEventType `json:"event"`
	Output string                   `json:"output"` // 输出内容
}

// StoppedEvent
// 该event表明，由于某些原因，被调试进程的执行已经停止。
// 这可能是由先前设置的断点、完成的步进请求、执行调试器语句等引起的。
type StoppedEvent struct {
	Event  constants.DebugEventType    `json:"event"`
	Reason constants.StoppedReasonType `json:"reason"` // 停止执行的原因
	Line   int                         `json:"line"`   // 停止在某行
}

// ContinuedEvent
// 该event表明debug的执行已经继续。
// 请注意:debug adapter不期望发送此事件来响应暗示执行继续的请求，例如启动或继续。
// 它只有在没有先前的request暗示这一点时，才有必要发送一个持续的事件。
type ContinuedEvent struct {
	Event constants.DebugEventType `json:"event"`
}

// ExitedEvent
// 该event表明被调试对象已经退出并返回exit code。
type ExitedEvent struct {
	Event    constants.DebugEventType `json:"event"`
	ExitCode int                      `json:"exitCode"`
	Message  string                   `json:"message"`
}

// TerminatedEvent
// 程序退出事件
type TerminatedEvent struct {
	Event constants.DebugEventType `json:"event"`
}
