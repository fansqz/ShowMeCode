package gdb_debugger

import (
	"context"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/fansqz/ShowMeCode/go-debugger/constants"
	. "github.com/fansqz/ShowMeCode/go-debugger/debugger"
	gdb2 "github.com/fansqz/ShowMeCode/go-debugger/debugger/gdb_debugger/gdb"
	. "github.com/fansqz/ShowMeCode/go-debugger/debugger/utils"
	"github.com/fansqz/ShowMeCode/go-debugger/utils/gosync"
	"github.com/google/go-dap"
	"github.com/sirupsen/logrus"
)

const (
	OptionTimeout = time.Second * 10
)

type GDBDebugger struct {
	startOption *StartOption

	// gdb调试的目标语言
	language constants.LanguageType
	// gdb示例
	GDB *gdb2.Gdb

	// functionInfos 语法解析解析出来的函数内容，局部变量获取需要通过静态代码分析内容获取变量列表
	functionInfos []FunctionInfo

	// 引用工具
	ReferenceUtil *ReferenceUtil

	// 事件产生时，触发该回调
	callback NotificationCallback

	// 调试的状态管理
	StatusManager *StatusManager

	// gdb输出工具，用于处理gdb输出
	GdbOutputUtil *GDBOutputUtil

	// 断点记录
	mutex             sync.RWMutex
	breakpointNumbers []string

	initBreakpointCountLock sync.Mutex
	NotInitBreakpointCount  int
	breakpointInitChannel   chan struct{}

	// 由于为了防止stepIn操作会进入系统依赖内部的特殊处理
	preAction               string // 记录gdb上一个命令
	skipContinuedEventCount int64  //记录需要跳过continue事件的数量，读写时必须加锁
}

func NewGDBDebugger(languageType constants.LanguageType) *GDBDebugger {
	d := &GDBDebugger{
		StatusManager:         NewStatusManager(),
		breakpointInitChannel: make(chan struct{}, 2),
		GdbOutputUtil:         NewGDBOutputUtil(),
		ReferenceUtil:         NewReferenceUtil(),
		language:              languageType,
	}
	return d
}

func (g *GDBDebugger) Start(option *StartOption) error {
	// 设置并创建gdb
	g.callback = option.Callback
	g.startOption = option

	// 异步进行语法解析
	if g.startOption.MainCode != "" {
		go func() {
			var err error
			g.functionInfos, err = ParseSourceFile(g.startOption.MainCode)
			if err != nil {
				logrus.Errorf("ParseSourceFile fail, err = %v", err)
			}
		}()
	}

	gd, err := gdb2.New(g.gdbNotificationCallback)
	if err != nil {
		log.Printf("Start fail, err = %s\n", err)
		return err
	}
	g.GDB = gd
	// 创建命令
	m, _ := g.GDB.Send("file-exec-and-symbols", option.ExecFile)
	if result, ok := m["class"]; ok && result == "done" {
		return nil
	} else {
		return fmt.Errorf("目标代码加载失败")
	}
}

// Run 同步方法，开始运行
func (g *GDBDebugger) Run() error {
	var gdbCallback gdb2.AsyncCallback = func(m map[string]interface{}) {
		gosync.Go(context.Background(), g.processUserInput)
		// 启动协程读取用户输出
		gosync.Go(context.Background(), g.processUserOutput)
	}
	// 设置语言
	_, _ = g.GDB.SendWithTimeout(OptionTimeout, "gdb-set", "language", string(g.language))
	g.preAction = "exec-run"
	if err := g.GDB.SendAsync(gdbCallback, "exec-run"); err != nil {
		log.Printf("Run fail, err = %s\n", err)
		// 启动失败
		return err
	}
	return nil
}

// processUserInput 循环读取用户输入
func (g *GDBDebugger) processUserInput(ctx context.Context) {
	var input string
	for {
		if g.StatusManager == nil || g.GDB == nil {
			break
		}
		if g.StatusManager.Is(Finish) {
			break
		}

		_, err := fmt.Scanln(&input)
		if err == nil {
			if input[len(input)-1] != '\n' {
				input = input + "\n"
			}
			if _, err = g.GDB.Write([]byte(input)); err != nil {
				logrus.Errorf("processUserInput fail, err = %s\n", err)
			}
		}
		time.Sleep(100 * time.Microsecond)
	}
}

// processUserOutput 循环处理用户输出
func (g *GDBDebugger) processUserOutput(ctx context.Context) {
	b := make([]byte, 1024)
	for {
		n, err := g.GDB.Read(b)
		if err != nil {
			return
		}
		output := string(b[0:n])
		fmt.Print(output)
	}
}

func (g *GDBDebugger) StepOver() error {
	if !g.StatusManager.Is(Stopped) {
		return errors.New("程序运行中，无法执行单步调试")
	}
	return g.stepOver()
}

func (g *GDBDebugger) stepOver() error {
	g.preAction = "exec-next"
	err := g.GDB.SendAsync(func(obj map[string]interface{}) {}, "exec-next")
	return err
}

func (g *GDBDebugger) StepIn() error {
	if !g.StatusManager.Is(Stopped) {
		return errors.New("程序运行中，无法执行单步调试")
	}
	return g.stopIn()
}

func (g *GDBDebugger) stopIn() error {
	g.preAction = "exec-step"
	err := g.GDB.SendAsync(func(obj map[string]interface{}) {}, "exec-step")
	return err
}

func (g *GDBDebugger) StepOut() error {
	if !g.StatusManager.Is(Stopped) {
		return errors.New("程序运行中，无法执行单步调试")
	}
	return g.stopOut()
}

func (g *GDBDebugger) stopOut() error {
	g.preAction = "exec-finish"
	err := g.GDB.SendAsync(func(obj map[string]interface{}) {}, "exec-finish")
	return err
}

func (g *GDBDebugger) Continue() error {
	if !g.StatusManager.Is(Stopped) {
		return errors.New("程序运行中，无法执行continue")
	}
	return g.continue2()
}

func (g *GDBDebugger) continue2() error {
	g.preAction = "exec-continue"
	err := g.GDB.SendAsync(func(obj map[string]interface{}) {}, "exec-continue")
	return err
}

func (g *GDBDebugger) SetBreakpoints(source dap.Source, breakpoints []dap.SourceBreakpoint) error {
	g.mutex.Lock()
	// 删除原来的所有断点
	err := g.removeBreakpoints(g.breakpointNumbers)
	if err != nil {
		return err
	}
	for _, bp := range breakpoints {
		result, err := g.GDB.SendWithTimeout(OptionTimeout, "break-insert", source.Path+":"+strconv.Itoa(bp.Line))
		if err != nil {
			continue
		} else {
			success, number := g.GdbOutputUtil.ParseAddBreakpointOutput(result)
			if success {
				g.breakpointNumbers = append(g.breakpointNumbers, number)
			}
		}
	}
	g.mutex.Unlock()
	return nil
}

func (g *GDBDebugger) removeBreakpoints(numbers []string) error {
	for _, number := range numbers {
		var callback gdb2.AsyncCallback = func(m map[string]interface{}) {}
		if err := g.GDB.SendAsync(callback, "break-delete", number); err != nil {
			return err
		}
	}
	return nil
}

func (g *GDBDebugger) Terminate() error {
	if g.StatusManager.Is(Finish) {
		return nil
	}
	// 发送终端给程序
	err := g.GDB.Interrupt()
	if err != nil {
		log.Printf("Terminate fail, err = %s\n", err)
		return err
	}
	_ = g.GDB.Exit()
	// 保证map的线程安全
	g.GdbOutputUtil.lock.Lock()
	defer g.GdbOutputUtil.lock.Unlock()
	g.skipContinuedEventCount = 0
	g.StatusManager.Set(Finish)
	return nil
}

func (g *GDBDebugger) GetStackTrace() ([]dap.StackFrame, error) {
	if !g.StatusManager.Is(Stopped) {
		return nil, errors.New("程序未暂停无法获取栈帧信息")
	}
	m, err := g.sendWithTimeOut(OptionTimeout, "stack-list-frames")
	if err != nil {
		log.Printf("GetStackTrace fail, err = %s\n", err)
		return nil, err
	}
	return g.GdbOutputUtil.ParseStackTraceOutput(m), nil
}

func (g *GDBDebugger) GetScopes(frameId int) ([]dap.Scope, error) {
	// 读取栈帧
	return []dap.Scope{
		{Name: "Global", VariablesReference: globalScopeReference},
		{Name: "Local", VariablesReference: g.ReferenceUtil.GetScopesReference(frameId)},
	}, nil
}

// GetVariables 获取变量列表
func (g *GDBDebugger) GetVariables(reference int) ([]dap.Variable, error) {
	// 解析引用
	refStruct, err := g.ReferenceUtil.ParseVariableReference(reference)
	if err != nil {
		log.Printf("getVariables failed: %v\n", err)
		return nil, err
	}

	// 切换栈帧
	if err = g.SelectFrame(refStruct); err != nil {
		return nil, err
	}

	// 创建变量
	targetVar := "structName"
	_, err = g.CreateVar(refStruct, targetVar)
	if err != nil {
		return nil, err
	}
	defer func(g *GDBDebugger, name string) {
		err := g.DeleteVar(name)
		if err != nil {
			log.Printf("DeleteVar failed: %v\n", err)
		}
	}(g, targetVar)

	// 读取变量的children元素列表
	variables, err := g.varListChildren(targetVar)
	if err != nil {
		return nil, err
	}
	return variables, nil
}

// GetLocalScopeVariables 获取局部变量列表
func (g *GDBDebugger) GetLocalScopeVariables(reference int) ([]dap.Variable, error) {
	var variables []dap.Variable
	var err error
	if g.functionInfos != nil {
		variables, err = g.getLocalVariables2(reference)
	} else {
		variables, err = g.getLocalVariables(reference)
	}
	if err != nil {
		return nil, err
	}
	return variables, nil
}

// GetGlobalVariables 获取全局变量列表
func (g *GDBDebugger) GetGlobalVariables() ([]dap.Variable, error) {
	m, err := g.sendWithTimeOut(OptionTimeout, "symbol-info-variables", "--max-results", "40")
	if err != nil {
		return nil, err
	}
	variables := g.GdbOutputUtil.ParseGlobalVariableOutput(g.GDB, m)
	return variables, nil
}

// getLocalVariables 获取本地变量列表1，通过gdb命令获取，会存在有些未初始化变量却被使用情况
func (g *GDBDebugger) getLocalVariables(reference int) ([]dap.Variable, error) {
	frameId := g.ReferenceUtil.GetFrameIDByLocalReference(reference)
	// 获取当前线程id
	currentThreadId, _ := g.getCurrentThreadId()
	// 获取栈帧中所有局部变量
	m, err := g.sendWithTimeOut(OptionTimeout, "stack-list-variables",
		"--thread", currentThreadId, "--frame", strconv.Itoa(frameId), "--skip-unavailable", "2")
	if err != nil {
		log.Printf("getLocalScopeVariables failed: %v\n", err)
		return nil, err
	}
	variables := g.GdbOutputUtil.ParseFrameVariablesOutput(g.GDB, m)
	return variables, nil
}

// getLocalVariables2 通过静态代码分析获取
func (g *GDBDebugger) getLocalVariables2(reference int) ([]dap.Variable, error) {
	frameId := g.ReferenceUtil.GetFrameIDByLocalReference(reference)
	stackTrace, err := g.GetStackTrace()
	if err != nil {
		return g.getLocalVariables(reference)
	}
	// 找到目标栈帧id
	var targetFrame dap.StackFrame
	for _, f := range stackTrace {
		if f.Id == frameId {
			targetFrame = f
			break
		}
	}
	// 获取当前方法
	var target *FunctionInfo
	for _, f := range g.functionInfos {
		if f.Name == targetFrame.Name {
			target = &f
			break
		}
	}
	if target == nil {
		return g.getLocalVariables(reference)
	}

	// 获取目标变量列表
	targetVariableNames := []string{}
	for _, v := range target.Variables {
		if v.Location.Line < targetFrame.Line {
			targetVariableNames = append(targetVariableNames, v.Name)
		}
	}

	// 切换栈帧
	_, err = g.GDB.SendWithTimeout(OptionTimeout, "stack-select-frame", strconv.Itoa(frameId))
	if err != nil {
		return nil, err
	}
	// 读取变量列表
	var answer []dap.Variable
	for _, variableName := range targetVariableNames {
		m2, err := g.GDB.SendWithTimeout(OptionTimeout, "var-create", variableName, "*", variableName)
		if err != nil {
			logrus.Errorf("getChidrenNumber fail err = %s", err)
			continue
		}
		variable, ok := g.GdbOutputUtil.ParseVarCreate(m2)
		if !ok {
			continue
		}
		answer = append(answer, *variable)
		_, _ = g.GDB.SendWithTimeout(OptionTimeout, "var-delete", variableName)
	}
	return answer, nil
}

// SelectFrame 如果是普通类型需要切换栈帧，同一个变量名，可能在不同栈帧中会有重复，需要定位栈帧和变量名称才能读取到变量值
func (g *GDBDebugger) SelectFrame(ref *ReferenceStruct) error {
	if ref.Type == StructType {
		if _, err := g.sendWithTimeOut(OptionTimeout, "stack-select-frame", ref.FrameId); err != nil {
			return err
		}
	}
	return nil
}

// CreateVar 创建变量，在读取一个值的时候，需要创建变量以后才能读取
func (g *GDBDebugger) CreateVar(ref *ReferenceStruct, structName string) (*dap.Variable, error) {
	_, _ = g.GDB.SendWithTimeout(OptionTimeout, "enable-pretty-printing")
	exp := g.GetExport(ref)
	m, err := g.sendWithTimeOut(OptionTimeout, "var-create", structName, "*", exp)
	if err != nil {
		logrus.Errorf("create var fail %s", err)
		return nil, err
	}
	variable, ok := g.GdbOutputUtil.ParseVarCreate(m)
	if !ok {
		return nil, fmt.Errorf("create var fail")
	}
	return variable, nil
}

// GetExport 通过ReferenceStruct，获取变量表达式
func (g *GDBDebugger) GetExport(ref *ReferenceStruct) string {
	var exp string
	switch ref.Type {
	case StructType:
		exp = ref.VariableName
	case PointType:
		exp = fmt.Sprintf("(%s)%s", ref.VariableType, ref.Address)
	}
	if ref.FieldPath != "" {
		exp = fmt.Sprintf("(%s)%s", exp, ref.FieldPath)
	}
	return exp
}

// DeleteVar 删除变量，创建完变量需要进行删除，避免再次创建时名称重复
func (g *GDBDebugger) DeleteVar(name string) error {
	_, err := g.sendWithTimeOut(OptionTimeout, "var-delete", name)
	return err
}

func (g *GDBDebugger) parseObject2Keys(inputStr string) []string {
	// 定义正则表达式模式，匹配 = 前面的键
	re := regexp.MustCompile(`(\w+)\s*=`)
	// 查找所有匹配项
	matches := re.FindAllStringSubmatch(inputStr, -1)
	answer := []string{}
	for _, match := range matches {
		key := match[1]
		if key != "\000" {
			answer = append(answer, key)
		}
	}
	return answer
}

// varListChildren 读取变量的children元素列表 ）
func (g *GDBDebugger) varListChildren(structName string) ([]dap.Variable, error) {
	// 获取所有children列表并解析
	var m map[string]interface{}
	var err error
	m, err = g.sendWithTimeOut(OptionTimeout, "var-list-children", "1", structName)
	if err != nil {
		log.Printf("getVariables fail, err = %s\n", err)
		return nil, err
	}
	variables := g.GdbOutputUtil.ParseVariablesOutput(m)
	return variables, nil
}

// getCurrentThreadId 获取当前线程id
func (g *GDBDebugger) getCurrentThreadId() (string, error) {
	// 获取当前线程id
	m, err := g.sendWithTimeOut(OptionTimeout, "thread-info")
	if err != nil {
		log.Printf("getCurrentThreadId fail, err = %s\n", err)
		return "", err
	}
	threadMap, success := g.GdbOutputUtil.GetPayloadFromMap(m)
	if !success {
		return "", errors.New("获取线程id失败")
	}
	currentThreadId := g.GdbOutputUtil.GetStringFromMap(threadMap, "current-thread-id")
	return currentThreadId, nil
}

// gdbNotificationCallback 处理gdb异步响应的回调
func (g *GDBDebugger) gdbNotificationCallback(m map[string]interface{}) {
	typ := g.GdbOutputUtil.GetStringFromMap(m, "type")
	switch typ {
	case "exec":
		class := g.GdbOutputUtil.GetStringFromMap(m, "class")
		switch class {
		case "stopped":
			// 处理程序停止的事件
			g.processStoppedData(g.GdbOutputUtil.GetInterfaceFromMap(m, "payload"))
			g.StatusManager.Set(Stopped)
		case "running":
			g.processRunningData()
			g.StatusManager.Set(Running)
		}
	}

}

// processStoppedData 处理gdb返回的stopped数据，程序停止到程序的某个位置就会返回stopped data
func (g *GDBDebugger) processStoppedData(m interface{}) {
	stoppedOutput := g.GdbOutputUtil.ParseStoppedEventOutput(m)
	if stoppedOutput == nil {
		return
	}
	// 停留在断点
	if stoppedOutput.reason == constants.StepStopped || stoppedOutput.reason == constants.BreakpointStopped {
		// 返回停留的断点位置
		g.callback(&dap.StoppedEvent{
			Event: *NewEvent(0, "stopped"),
			Body:  dap.StoppedEventBody{Reason: string(stoppedOutput.reason)},
		})
	}
	if stoppedOutput.reason == constants.ExitedNormally {
		// 程序退出
		g.callback(&dap.TerminatedEvent{
			Event: *NewEvent(0, "terminated"),
		})
	}
}

// processRunningData 处理gdb返回的running事件
func (g *GDBDebugger) processRunningData() {
	// 程序执行，如果有需要跳过的continue事件，则跳过
	skipCount := atomic.LoadInt64(&g.skipContinuedEventCount)
	if skipCount > 0 {
		atomic.AddInt64(&g.skipContinuedEventCount, -1)
		return
	}
	g.callback(&dap.ContinuedEvent{
		Event: *NewEvent(0, "continued"),
		Body:  dap.ContinuedEventBody{},
	})
	// 设置用户程序为执行状态
	g.StatusManager.Set(Running)
}

func (g *GDBDebugger) sendWithTimeOut(timeout time.Duration, operation string, args ...string) (map[string]interface{}, error) {
	channel := make(chan map[string]interface{}, 1)

	err := g.GDB.SendAsync(func(obj map[string]interface{}) {
		channel <- obj
	}, operation, args...)
	if err != nil {
		return nil, err
	}
	select {
	case m := <-channel:
		return m, nil
	case <-time.After(timeout):
		return nil, errors.New("GetStackTrace time out")
	}
}
