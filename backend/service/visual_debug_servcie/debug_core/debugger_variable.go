package debug_core

import (
	"context"
	"fmt"
	"strings"

	"github.com/fansqz/fancode-backend/common/logger"
	"github.com/google/go-dap"
)

func (d *debugger) GetStackTrace(ctx context.Context) ([]*StackFrame, error) {
	logger.WithCtx(ctx).Infof("[debugger] GetStackTrace")
	stackFrames, err := d.getStackTrace()
	if err != nil {
		return nil, err
	}
	var answer []*StackFrame
	for _, s := range stackFrames {
		line := 0
		if s.Source.Path == d.compileFile {
			line = s.Line
		}
		answer = append(answer, &StackFrame{
			ID:   s.Id,
			Name: s.Name,
			Path: strings.Replace(s.Source.Path, d.workPath, "", 1),
			Line: line,
		})
	}
	return answer, nil
}

func (d *debugger) getStackTrace() ([]dap.StackFrame, error) {
	// 读取栈帧
	request := &dap.StackTraceRequest{Request: *d.dapCli.NewRequest("stackTrace")}
	request.Arguments.ThreadId = 1
	request.Arguments.StartFrame = 0
	request.Arguments.Levels = 20
	message, err := d.dapCli.SendWithTimeout(request, d.option.OptionTimeout)
	if err != nil {
		return nil, err
	}
	dapResponse, ok := message.(*dap.StackTraceResponse)
	if !ok {
		return nil, fmt.Errorf("message is not *dap.StackTraceResponse")
	}
	if !dapResponse.Success {
		return nil, fmt.Errorf(dapResponse.Message)
	}
	return dapResponse.Body.StackFrames, nil
}

func (d *debugger) GetFrameVariables(ctx context.Context, frameId int) ([]*Variable, error) {
	logger.WithCtx(ctx).Infof("[debugger] GetFrameVariables")

	// 读取scope
	srequest := &dap.ScopesRequest{Request: *d.dapCli.NewRequest("scopes")}
	srequest.Arguments.FrameId = frameId
	message, err := d.dapCli.SendWithTimeout(srequest, d.option.OptionTimeout)
	if err != nil {
		return nil, err
	}
	sresponse, ok := message.(*dap.ScopesResponse)
	if !ok {
		return nil, fmt.Errorf("message is not *dap.ScopesResponse")
	}
	if !sresponse.Success {
		return nil, fmt.Errorf(sresponse.Message)
	}

	// 读取local scope
	localScope := -1
	for _, scope := range sresponse.Body.Scopes {
		if strings.Contains(strings.ToLower(scope.Name), "local") {
			localScope = scope.VariablesReference
		}
	}

	// 读取变量列表
	if localScope == -1 {
		return []*Variable{}, nil
	}
	return d.GetVariables(ctx, localScope)
}

func (d *debugger) GetVariables(ctx context.Context, reference int) ([]*Variable, error) {
	logger.WithCtx(ctx).Infof("[GoDebugger] GetVariables")
	// 读取local的变量值
	resp, err := d.variables(reference)
	if err != nil {
		return nil, err
	}
	// 对象转换
	answer := []*Variable{}
	for _, v := range resp.Body.Variables {
		// 过滤~r0相关变量
		if strings.HasPrefix(v.Name, "~r") {
			continue
		}
		// 如果是指针，特殊处理
		if strings.HasPrefix(v.Type, "*") && v.VariablesReference != 0 {
			// 1. 获取指针指向的结构体的引用
			resp2, err2 := d.variables(v.VariablesReference)
			if err2 != nil {
				return nil, err2
			}
			var structV dap.Variable
			if len(resp2.Body.Variables) > 0 {
				structV = resp2.Body.Variables[0]
				v.VariablesReference = structV.VariablesReference
			}
		}
		answer = append(answer, &Variable{
			Name:             v.Name,
			Type:             v.Type,
			Value:            v.Value,
			Reference:        v.VariablesReference,
			NamedVariables:   v.NamedVariables,
			IndexedVariables: v.IndexedVariables,
		})
	}
	return answer, nil
}

func (d *debugger) variables(reference int) (*dap.VariablesResponse, error) {
	// 读取local的变量值
	vrequest := &dap.VariablesRequest{
		Request: *d.dapCli.NewRequest("variables"),
		Arguments: dap.VariablesArguments{
			VariablesReference: reference,
			Format: &dap.ValueFormat{
				Hex: true,
			},
		},
	}
	message, err := d.dapCli.SendWithTimeout(vrequest, d.option.OptionTimeout)
	if err != nil {
		return nil, err
	}
	vrespose, ok := message.(*dap.VariablesResponse)
	if !ok {
		return nil, fmt.Errorf("message is not *dap.VariablesResponse")
	}
	if !vrespose.Success {
		return nil, fmt.Errorf(vrespose.Message)
	}
	return vrespose, nil
}

func (d *debugger) evaluate(expr string, fid int, context string) (*dap.EvaluateResponse, error) {
	request := &dap.EvaluateRequest{Request: *d.dapCli.NewRequest("evaluate")}
	request.Arguments.Expression = expr
	request.Arguments.FrameId = fid
	request.Arguments.Context = context
	message, err := d.dapCli.SendWithTimeout(request, d.option.OptionTimeout)
	if err != nil {
		return nil, err
	}
	vrespose, ok := message.(*dap.EvaluateResponse)
	if !ok {
		return nil, fmt.Errorf("message is not *dap.EvaluateResponse")
	}
	if !vrespose.Success {
		return nil, fmt.Errorf(vrespose.Message)
	}
	return vrespose, nil
}
