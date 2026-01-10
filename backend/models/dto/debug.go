package dto

import "github.com/fansqz/fancode-backend/constants"

// =======================以下是request=============================

type CreateDebugSession struct {
	Language constants.LanguageType `json:"language"`
}

// StartDebugRequest 启动调试请求
type StartDebugRequest struct {
	ID string `json:"id"`
	// 用户代码
	Code string `json:"code"`
	// Language 调试语言
	Language constants.LanguageType `json:"language"`
	// 初始断点
	Breakpoints []int `json:"breakpoints"`
}

type BaseDebugRequest struct {
	ID string `json:"id"` // 调试id
}

type SendToConsoleRequest struct {
	ID    string `json:"id"`
	Input string `json:"input"`
}

type SetBreakpointRequest struct {
	ID          string `json:"id"`
	Breakpoints []int  `json:"breakpoints"`
}

type GetFrameVariablesRequest struct {
	ID      string `json:"id"`
	FrameID int    `json:"frameID"`
}

type GetVariablesRequest struct {
	ID        string `json:"id"`
	Reference int    `json:"reference"`
}

// CompileEvent
// 编译事件
type CompileEvent struct {
	Event   constants.DebugEventType `json:"event"`
	Success bool                     `json:"success"`
	Message string                   `json:"message"` // 编译产生的信息
}

func NewCompileEvent(success bool, message string) *CompileEvent {
	return &CompileEvent{
		Event:   constants.CompileEvent,
		Success: success,
		Message: message,
	}
}

// OutputEvent
// 该事件表明目标已经产生了一些输出。
type OutputEvent struct {
	Event  constants.DebugEventType `json:"event"`
	Output string                   `json:"output"` // 输出内容
}

func NewOutputEvent(output string) *OutputEvent {
	return &OutputEvent{
		Event:  constants.OutputEvent,
		Output: output,
	}
}

// StoppedEvent
// 该event表明，由于某些原因，被调试进程的执行已经停止。
// 这可能是由先前设置的断点、完成的步进请求、执行调试器语句等引起的。
type StoppedEvent struct {
	Event  constants.DebugEventType    `json:"event"`
	Reason constants.StoppedReasonType `json:"reason"` // 停止执行的原因
}

func NewStoppedEvent(reason constants.StoppedReasonType) *StoppedEvent {
	return &StoppedEvent{
		Event:  constants.StoppedEvent,
		Reason: reason,
	}
}

// ContinuedEvent
// 该event表明debug的执行已经继续。
// 请注意:debug adapter不期望发送此事件来响应暗示执行继续的请求，例如启动或继续。
// 它只有在没有先前的request暗示这一点时，才有必要发送一个持续的事件。
type ContinuedEvent struct {
	Event constants.DebugEventType `json:"event"`
}

func NewContinuedEvent() *ContinuedEvent {
	return &ContinuedEvent{
		Event: constants.ContinuedEvent,
	}
}

// TerminatedEvent
// 程序退出事件
type TerminatedEvent struct {
	Event constants.DebugEventType `json:"event"`
}

func NewTerminatedEvent() *TerminatedEvent {
	return &TerminatedEvent{
		Event: constants.TerminatedEvent,
	}
}
