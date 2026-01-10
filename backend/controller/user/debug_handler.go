package user

import (
	e "github.com/fansqz/fancode-backend/common/error"
	"github.com/fansqz/fancode-backend/models/dto"
	r "github.com/fansqz/fancode-backend/models/vo"
	"github.com/fansqz/fancode-backend/service/visual_debug_servcie"

	"github.com/gin-gonic/gin"
)

type DebugController interface {
	CreateDebugSession(ctx *gin.Context)
	// Start 启动调试
	Start(ctx *gin.Context)
	// CreateSseConnect 会创建一个sse链接，用于接受服务器响应
	CreateSseConnect(ctx *gin.Context)
	// SendToConsole 提交
	SendToConsole(ctx *gin.Context)
	// StepIn 单步调试，会进入函数内部
	StepIn(ctx *gin.Context)
	// StepOut 单步调试，会跳出当前程序
	StepOut(ctx *gin.Context)
	// StepOver 单步调试，跳过不进入程序内部
	StepOver(ctx *gin.Context)
	// Continue 到达下一个断点
	Continue(ctx *gin.Context)
	// SetBreakpoints 设置断点
	SetBreakpoints(ctx *gin.Context)
	// GetStackTrace 获取当前栈信息
	GetStackTrace(ctx *gin.Context)
	// GetFrameVariables 根据栈帧id获取变量列表
	GetFrameVariables(ctx *gin.Context)
	// GetVariables 根据引用获取变量信息，如果是指针，获取指针指向的内容，如果是结构体，获取结构体内容
	GetVariables(ctx *gin.Context)
	// Terminate 关闭调试
	Terminate(ctx *gin.Context)
}

type debugController struct {
	debugService visual_debug_servcie.DebugService
}

func NewDebugController(ds visual_debug_servcie.DebugService) DebugController {
	return &debugController{
		debugService: ds,
	}
}

func (d *debugController) CreateDebugSession(ctx *gin.Context) {
	result := r.NewResult(ctx)
	id, err := d.debugService.CreateDebugSession(ctx)
	if err != nil {
		result.Error(err)
		return
	}
	result.SuccessData(id)
}

// Start 开始调试
func (d *debugController) Start(ctx *gin.Context) {
	result := r.NewResult(ctx)
	var startReq dto.StartDebugRequest
	if err := ctx.BindJSON(&startReq); err != nil {
		return
	}
	err := d.debugService.Start(ctx, startReq)
	if err != nil {
		result.Error(err)
		return
	}
	result.SuccessMessage("启动成功")
}

// CreateSseConnect
func (d *debugController) CreateSseConnect(ctx *gin.Context) {
	id := ctx.Param("id")
	d.debugService.CreateSseConnect(ctx, id)
}

// SendToConsole 提交
func (d *debugController) SendToConsole(ctx *gin.Context) {
	result := r.NewResult(ctx)
	var req dto.SendToConsoleRequest
	if err := ctx.BindJSON(&req); err != nil {
		result.Error(e.ErrBadRequest)
		return
	}
	if err := d.debugService.SendToConsole(ctx, req.ID, req.Input); err != nil {
		result.Error(err)
		return
	}
	result.SuccessMessage("请求成功")
}

// Step
func (d *debugController) StepIn(ctx *gin.Context) {
	result := r.NewResult(ctx)
	var req dto.BaseDebugRequest
	if err := ctx.BindJSON(&req); err != nil {
		result.Error(e.ErrBadRequest)
		return
	}
	if err := d.debugService.StepIn(ctx, req.ID); err != nil {
		result.Error(err)
		return
	}
	result.SuccessMessage("请求成功")
}

// Step
func (d *debugController) StepOut(ctx *gin.Context) {
	result := r.NewResult(ctx)
	var req dto.BaseDebugRequest
	if err := ctx.BindJSON(&req); err != nil {
		result.Error(e.ErrBadRequest)
		return
	}
	if err := d.debugService.StepOut(ctx, req.ID); err != nil {
		result.Error(err)
		return
	}
	result.SuccessMessage("请求成功")
}

// Step
func (d *debugController) StepOver(ctx *gin.Context) {
	result := r.NewResult(ctx)
	var req dto.BaseDebugRequest
	if err := ctx.BindJSON(&req); err != nil {
		result.Error(e.ErrBadRequest)
		return
	}
	if err := d.debugService.StepOver(ctx, req.ID); err != nil {
		result.Error(err)
		return
	}
	result.SuccessMessage("请求成功")
}

// Continue
func (d *debugController) Continue(ctx *gin.Context) {
	result := r.NewResult(ctx)
	var req dto.BaseDebugRequest
	if err := ctx.BindJSON(&req); err != nil {
		result.Error(e.ErrBadRequest)
		return
	}
	if err := d.debugService.Continue(ctx, req.ID); err != nil {
		result.Error(err)
		return
	}
	result.SuccessMessage("请求成功")
}

// SetBreakpoints
func (d *debugController) SetBreakpoints(ctx *gin.Context) {
	result := r.NewResult(ctx)
	var req dto.SetBreakpointRequest
	if err := ctx.BindJSON(&req); err != nil {
		result.Error(e.ErrBadRequest)
		return
	}
	response, err := d.debugService.SetBreakpoints(ctx, req.ID, req.Breakpoints)
	if err != nil {
		result.Error(e.ErrBadRequest)
		return
	}
	result.SuccessData(response)
}

// Terminate
func (d *debugController) Terminate(ctx *gin.Context) {
	result := r.NewResult(ctx)
	var req dto.BaseDebugRequest
	if err := ctx.BindJSON(&req); err != nil {
		result.Error(e.ErrBadRequest)
		return
	}
	if err := d.debugService.Terminate(ctx, req.ID); err != nil {
		result.Error(err)
		return
	}
	result.SuccessMessage("请求成功")
}

func (d *debugController) GetStackTrace(ctx *gin.Context) {
	result := r.NewResult(ctx)
	var req dto.BaseDebugRequest
	if err := ctx.BindJSON(&req); err != nil {
		result.Error(e.ErrBadRequest)
		return
	}
	stackFrames, err := d.debugService.GetStackTrace(ctx, req.ID)
	if err != nil {
		result.Error(err)
		return
	}
	result.SuccessData(stackFrames)
}

func (d *debugController) GetFrameVariables(ctx *gin.Context) {
	result := r.NewResult(ctx)
	var req dto.GetFrameVariablesRequest
	if err := ctx.BindJSON(&req); err != nil {
		result.Error(e.ErrBadRequest)
		return
	}
	variables, err := d.debugService.GetFrameVariables(ctx, req.ID, req.FrameID)
	if err != nil {
		result.Error(err)
		return
	}
	result.SuccessData(variables)
}

// GetVariables 根据引用获取变量信息，如果是指针，获取指针指向的内容，如果是结构体，获取结构体内容
func (d *debugController) GetVariables(ctx *gin.Context) {
	result := r.NewResult(ctx)
	var req dto.GetVariablesRequest
	if err := ctx.BindJSON(&req); err != nil {
		result.Error(e.ErrBadRequest)
		return
	}
	variables, err := d.debugService.GetVariables(ctx, req.ID, req.Reference)
	if err != nil {
		result.Error(err)
		return
	}
	result.SuccessData(variables)
}
