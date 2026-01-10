package visual_debug_servcie

import (
	"context"
	json2 "encoding/json"
	"fmt"
	"github.com/fansqz/fancode-backend/common/config"
	e "github.com/fansqz/fancode-backend/common/error"
	"github.com/fansqz/fancode-backend/common/logger"
	"github.com/fansqz/fancode-backend/constants"
	"github.com/fansqz/fancode-backend/models/dto"
	"github.com/fansqz/fancode-backend/models/vo"
	"github.com/fansqz/fancode-backend/service/visual_debug_servcie/debug_core"
	"github.com/fansqz/fancode-backend/utils"
	"log"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/gin-gonic/gin"
)

// DebugService
// 用户调试相关
type DebugService interface {
	// CreateDebugSession 创建调试session
	CreateDebugSession(ctx context.Context) (string, error)
	// CreateSseConnect 创建sse连接
	CreateSseConnect(ctx context.Context, id string)

	// Start 加载并启动用户程序
	Start(ctx context.Context, startReq dto.StartDebugRequest) error
	SendToConsole(ctx context.Context, id string, input string) error
	StepIn(ctx context.Context, id string) error
	StepOver(ctx context.Context, id string) error
	StepOut(ctx context.Context, id string) error
	Continue(ctx context.Context, id string) error
	SetBreakpoints(ctx context.Context, id string, breakpoints []int) ([]*debug_core.Breakpoint, error)
	GetStackTrace(ctx context.Context, key string) ([]*debug_core.StackFrame, error)
	GetFrameVariables(ctx context.Context, key string, frameId int) ([]*debug_core.Variable, error)
	GetVariables(ctx context.Context, key string, reference int) ([]*debug_core.Variable, error)
	// Terminate 关闭用户程序并关闭调试session
	Terminate(ctx context.Context, key string) error
}

type debugService struct {
	config *config.AppConfig
}

func NewDebugService(cf *config.AppConfig) DebugService {
	return &debugService{
		config: cf,
	}
}

func (d *debugService) CreateDebugSession(ctx context.Context) (string, error) {
	// 检查是否用户已经有调试程序
	debugID, err := UserDebugManage.GetUserDebugID(ctx)
	if err != nil {
		return "", err
	}
	if debugID != "" {
		debugContext, ok := DebugSessionManage.GetDebugSession(ctx, debugID)
		if ok {
			_ = debugContext.Debugger.Terminate(ctx)
			DebugSessionManage.SendDestroyEvent(ctx, debugID)
		}
	}
	key := utils.GetUUID()
	if err = DebugSessionManage.CreateDebugSession(ctx, d.config.AIConfig, key); err != nil {
		logger.WithCtx(ctx).Errorf("[CreateDebugSession] CreateDebugSession fail, err = %s", err)
		return "", err
	}
	if err = UserDebugManage.StoreUserDebugID(ctx, key); err != nil {
		logger.WithCtx(ctx).Errorf("[CreateDebugSession] StoreUserDebugID fail, err = %s", err)
		return "", err
	}
	return key, nil
}

func (d *debugService) CreateSseConnect(ctx context.Context, key string) {
	ginContext := ctx.(*gin.Context)
	result := vo.NewResult(ginContext)
	session, y := DebugSessionManage.GetDebugSession(ctx, key)
	if !y {
		result.SimpleErrorMessage("key 不存在")
		return
	}
	w := ginContext.Writer
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Transfer-Encoding", "chunked")
	fmt.Fprintf(w, "data: %s\n\n", "connect success")
	// 刷新缓冲，确保立即发送到客户端
	flusher, _ := w.(http.Flusher)
	flusher.Flush()
	// 遍历channel获取event并发送给前端
	for {
		select {
		case event := <-session.DtoEventChan:
			json, err := json2.Marshal(event)
			if err != nil {
				continue
			}
			// 写入事件数据
			fmt.Fprintf(w, "data: %s\n\n", string(json))
			// 刷新缓冲，确保立即发送到客户端
			flusher.Flush()
		case <-session.DestroyEventChan:
			json, _ := json2.Marshal(dto.NewTerminatedEvent())
			// 写入关闭事件
			fmt.Fprintf(w, "data: %s\n\n", json)
			// 刷新缓冲，确保立即发送到客户端
			flusher.Flush()
			// 销毁session
			DebugSessionManage.DestroyDebugSession(ctx, key)
			return
		}
	}
}

func (d *debugService) Terminate(ctx context.Context, debugID string) error {
	ctx = logger.SetDebugID(ctx, debugID)
	// 获取调试上下文
	debugContext, ok := DebugSessionManage.GetDebugSession(ctx, debugID)
	if !ok {
		return e.ErrDebuggerIsClosed
	}
	if err := debugContext.Debugger.Terminate(ctx); err != nil {
		logger.WithCtx(ctx).Errorf("[Terminate] terminate fail, err = %s", err)
		return err
	}
	// 发送关闭session的事件
	DebugSessionManage.SendDestroyEvent(ctx, debugID)
	return nil
}

func (d *debugService) Start(ctx context.Context, startReq dto.StartDebugRequest) error {
	// 如果已经获取不到，代表已经被提前close了
	if _, ok := DebugSessionManage.GetDebugSession(ctx, startReq.ID); !ok {
		return e.ErrDebuggerIsClosed
	}
	ctx = logger.SetDebugID(ctx, startReq.ID)

	debugSession, ok := DebugSessionManage.GetDebugSession(ctx, startReq.ID)
	if !ok {
		return e.ErrUnknown
	}
	debugge := debugSession.Debugger
	visualDescriptionAnalyzer := debugSession.VisualDescriptionAnalyzer

	// 启动代码分析
	if err := visualDescriptionAnalyzer.StartAnalyzeCode(ctx, startReq.Code, startReq.Language); err != nil {
		return fmt.Errorf("start analyze code error, err = %v", err)
	}

	//启动用户程序
	return debugge.Start(ctx, &debug_core.Option{
		Language:      startReq.Language,
		Code:          startReq.Code,
		BreakPoints:   startReq.Breakpoints,
		TempDir:       d.config.FilePathConfig.TempDir,
		DebuggerImage: d.config.DebuggerImage,
		// Callback 事件回调
		Callback:       d.notificationCallback(ctx, startReq.ID),
		CompileTimeout: 30 * time.Second,
		OptionTimeout:  2 * time.Second,
		DebugTimeout:   10 * 60 * time.Second,
		MemoryLimit:    1024 * 1024 * 1024,
		CPUQuota:       5 * 60 * 1000000,
	})
}

// notificationCallback 处理调试对象返回的事件
func (d *debugService) notificationCallback(ctx context.Context, debugID string) func(data interface{}) {
	return func(data interface{}) {
		defer func() {
			if r := recover(); r != nil {
				log.Println(r)
			}
		}()
		event := d.getDebuggerEventToDtoEvent(data)
		// 发送event给用户
		DebugSessionManage.SendEvent(ctx, debugID, event)

		// 程序退出，发送销毁事件
		if _, ok := event.(*dto.TerminatedEvent); ok {
			DebugSessionManage.SendDestroyEvent(ctx, debugID)
		}
	}
}

func (d *debugService) sendEventToSse(ctx context.Context, key string, event interface{}) error {
	session, ok := DebugSessionManage.GetDebugSession(ctx, key)
	if !ok {
		logger.WithCtx(ctx).Warnf("[sendEventToSse], key not exist, key = %v", key)
		return e.KeyNotExistError
	}
	session.DtoEventChan <- event
	return nil
}

func (d *debugService) getDebuggerEventToDtoEvent(data interface{}) interface{} {
	var event interface{}
	if oevent, ok := data.(*debug_core.OutputEvent); ok {
		event = dto.NewOutputEvent(oevent.Output)
	}
	if sevent, ok := data.(*debug_core.StoppedEvent); ok {
		// 如果是stop
		event = dto.NewStoppedEvent(sevent.Reason)
	}
	if _, ok := data.(*debug_core.ContinuedEvent); ok {
		event = dto.NewContinuedEvent()
	}
	if _, ok := data.(*debug_core.TerminalEvent); ok {
		event = dto.NewTerminatedEvent()
	}
	if cevent, ok := data.(*debug_core.CompileEvent); ok {
		event = dto.NewCompileEvent(cevent.Success, cevent.Message)
	}
	return event
}

func (d *debugService) SendToConsole(ctx context.Context, debugID string, input string) error {
	ctx = logger.SetDebugID(ctx, debugID)
	// 获取调试上下文
	debugContext, ok := DebugSessionManage.GetDebugSession(ctx, debugID)
	if !ok {
		return e.ErrDebuggerIsClosed
	}
	if err := debugContext.Debugger.Send(ctx, input); err != nil {
		logger.WithCtx(ctx).Errorf("[SendToConsole] send to console error, err = %v", err)
		return err
	}
	return nil
}

func (d *debugService) StepIn(ctx context.Context, debugID string) error {
	ctx = logger.SetDebugID(ctx, debugID)
	// 获取调试上下文
	debugContext, ok := DebugSessionManage.GetDebugSession(ctx, debugID)
	if !ok {
		return e.ErrDebuggerIsClosed
	}
	if err := debugContext.Debugger.StepIn(ctx); err != nil {
		logger.WithCtx(ctx).Errorf("[StepIn] step in, err = %v", err)
		return err
	}
	return nil
}

func (d *debugService) StepOut(ctx context.Context, debugID string) error {
	ctx = logger.SetDebugID(ctx, debugID)
	// 获取调试上下文
	debugContext, ok := DebugSessionManage.GetDebugSession(ctx, debugID)
	if !ok {
		return e.ErrDebuggerIsClosed
	}
	if err := debugContext.Debugger.StepOut(ctx); err != nil {
		logger.WithCtx(ctx).Errorf("[StepOut] step out error, err = %v", err)
		return err
	}
	return nil
}

func (d *debugService) StepOver(ctx context.Context, debugID string) error {
	ctx = logger.SetDebugID(ctx, debugID)
	// 获取调试上下文
	debugContext, ok := DebugSessionManage.GetDebugSession(ctx, debugID)
	if !ok {
		return e.ErrDebuggerIsClosed
	}
	if err := debugContext.Debugger.StepOver(ctx); err != nil {
		logger.WithCtx(ctx).Errorf("[StepOver] step over error, err = %v", err)
		return err
	}
	return nil
}

func (d *debugService) Continue(ctx context.Context, debugID string) error {
	ctx = logger.SetDebugID(ctx, debugID)
	// 获取调试上下文
	debugContext, ok := DebugSessionManage.GetDebugSession(ctx, debugID)
	if !ok {
		return e.ErrDebuggerIsClosed
	}
	if err := debugContext.Debugger.Continue(ctx); err != nil {
		logger.WithCtx(ctx).Errorf("[Continue] continue error, err = %v", err)
		return err
	}
	return nil
}

func (d *debugService) SetBreakpoints(ctx context.Context, debugID string, breakpoints []int) ([]*debug_core.Breakpoint, error) {
	ctx = logger.SetDebugID(ctx, debugID)
	// 获取调试上下文
	debugContext, ok := DebugSessionManage.GetDebugSession(ctx, debugID)
	if !ok {
		return nil, e.ErrDebuggerIsClosed
	}
	answer, err := debugContext.Debugger.SetBreakpoints(ctx, breakpoints)
	if err != nil {
		logger.WithCtx(ctx).Errorf("[AddBreakpoints] add break point error, err = %v", err)
		return nil, e.ErrUnknown
	}
	return answer, nil
}

func (d *debugService) GetStackTrace(ctx context.Context, debugID string) ([]*debug_core.StackFrame, error) {
	ctx = logger.SetDebugID(ctx, debugID)
	debugContext, ok := DebugSessionManage.GetDebugSession(ctx, debugID)
	if !ok {
		return nil, e.ErrDebuggerIsClosed
	}
	stackFrames, err := debugContext.Debugger.GetStackTrace(ctx)
	if err != nil {
		logger.WithCtx(ctx).Errorf("[GetStackTrace] get stack trace error, err = %v", err)
		return nil, err
	}
	return stackFrames, nil
}

func (d *debugService) GetFrameVariables(ctx context.Context, debugID string, frameId int) ([]*debug_core.Variable, error) {
	ctx = logger.SetDebugID(ctx, debugID)
	debugContext, ok := DebugSessionManage.GetDebugSession(ctx, debugID)
	if !ok {
		return nil, e.ErrDebuggerIsClosed
	}
	variables, err := debugContext.Debugger.GetFrameVariables(ctx, frameId)
	if err != nil {
		logger.WithCtx(ctx).Errorf("[GetFrameVariables] get frame variables error, err = %v", err)
		return nil, err
	}
	return variables, nil
}

func (d *debugService) GetVariables(ctx context.Context, debugID string, reference int) ([]*debug_core.Variable, error) {
	ctx = logger.SetDebugID(ctx, debugID)
	debugContext, ok := DebugSessionManage.GetDebugSession(ctx, debugID)
	if !ok {
		return nil, e.ErrDebuggerIsClosed
	}
	variables, err := debugContext.Debugger.GetVariables(ctx, reference)
	if err != nil {
		logger.WithCtx(ctx).Errorf("[GetVariables] get variables error, err = %v", err)
		return nil, err
	}
	return variables, nil
}

// saveUserCode
// 保存用户代码到用户的executePath，并返回需要编译的文件列表
func (d *debugService) saveUserCode(ctx context.Context, language constants.LanguageType, codeStr string, executePath string) ([]string, error) {
	var compileFiles []string
	var mainFile string
	var err error

	if mainFile, err = getMainFileNameByLanguage(language); err != nil {
		logger.WithCtx(ctx).Errorf("[saveUserCode] get mail test_file error, err = %v", err)
		return nil, err
	}
	if err = os.WriteFile(path.Join(executePath, mainFile), []byte(codeStr), 0644); err != nil {
		logger.WithCtx(ctx).Errorf("[saveUserCode] write test_file error, err = %v", err)
		return nil, err
	}
	// 将main文件进行编译即可
	compileFiles = []string{path.Join(executePath, mainFile)}

	return compileFiles, nil
}

// 根据编程语言获取该编程语言的Main文件名称
func getMainFileNameByLanguage(language constants.LanguageType) (string, error) {
	switch language {
	case constants.LanguageC:
		return "main.c", nil
	case constants.LanguageJava:
		return "Main.java", nil
	case constants.LanguageGo:
		return "main.go", nil
	case constants.LanguageCPP:
		return "main.cpp", nil
	default:
		return "", e.ErrLanguageNotSupported
	}
}
