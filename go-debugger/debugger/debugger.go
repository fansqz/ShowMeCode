package debugger

import (
	"github.com/google/go-dap"
)

type NotificationCallback func(message dap.EventMessage)

// Debugger
// 用户的一次调试过程处理
// debugger目前设置为支持多文件的
// 需要保证并发安全
type Debugger interface {
	// Start
	// 开始调试，及调用start命令，callback用来异步处理用户程序输出
	Start(option *StartOption) error
	// Run 启动程序执行
	Run() error
	// StepOver 下一步，不会进入函数内部
	StepOver() error
	// StepIn 下一步，会进入函数内部
	StepIn() error
	// StepOut 单步退出
	StepOut() error
	// Continue 忽略继续执行
	Continue() error
	// SetBreakpoints 设置断点
	SetBreakpoints(dap.Source, []dap.SourceBreakpoint) error
	// GetStackTrace 获取栈帧
	GetStackTrace() ([]dap.StackFrame, error)
	// GetScopes 获取scopes
	GetScopes(frameId int) ([]dap.Scope, error)
	// GetVariables 查看引用的值
	GetVariables(reference int) ([]dap.Variable, error)
	// Terminate 终止调试
	// 调用完该命令以后可以重新Launch
	Terminate() error
}

// NewEvent builds an Event struct with the specified fields.
func NewEvent(seq int, event string) *dap.Event {
	return &dap.Event{
		ProtocolMessage: dap.ProtocolMessage{
			Seq:  seq,
			Type: "event",
		},
		Event: event,
	}
}
