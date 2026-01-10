package protocol

import "github.com/fansqz/ShowMeCode/go-debugger/constants"

type RequestInterface interface {
	GetSequence() uint
	SetSequence(uint)
}

type BaseRequest struct {
	Type     constants.RequestType `json:"type"`
	Sequence uint                  `json:"sequence"`
}

func (b *BaseRequest) GetSequence() uint {
	return b.Sequence
}

func (b *BaseRequest) SetSequence(s uint) {
	b.Sequence = s
}

// StartDebugRequest 启动调试请求
type StartDebugRequest struct {
	BaseRequest
	// 请求序列号
	// 用户代码
	Code string `json:"code"`
	// Language 调试语言
	Language constants.LanguageType `json:"language"`
	// 初始断点
	Breakpoints []int `json:"breakpoints"`
}

// AddBreakpointRequest 增加断点请求
type AddBreakpointRequest struct {
	BaseRequest
	Breakpoints []int `json:"breakpoints"`
}

// RemoveBreakpointRequest 移除断点请求
type RemoveBreakpointRequest struct {
	BaseRequest
	Breakpoints []int `json:"breakpoints"`
}

// SendToConsoleRequest 输入到控制台
type SendToConsoleRequest struct {
	BaseRequest
	Content string `json:"content"`
}

// StepRequest next
type StepRequest struct {
	BaseRequest
	StepType constants.StepType `json:"stepType"`
}

// ContinueRequest continue
type ContinueRequest struct {
	BaseRequest
}

// TerminateRequest 关闭调试
type TerminateRequest struct {
	BaseRequest
}

// GetFrameVariables 根据栈帧获取变量列表
type GetFrameVariables struct {
	BaseRequest
	FrameID string `json:"frameID"`
}

// GetStackTrace 获取栈列表
type GetStackTrace struct {
	BaseRequest
}

// GetVariablesRequest 根据引用获取变量信息，如果是指针，获取指针指向的内容，如果是结构体，获取结构体内容
type GetVariablesRequest struct {
	BaseRequest
	Reference constants.RequestType `json:"reference"`
}
