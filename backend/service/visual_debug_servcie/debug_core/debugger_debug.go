package debug_core

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/docker/docker/api/types/container"
	e "github.com/fansqz/fancode-backend/common/error"
	"github.com/fansqz/fancode-backend/common/logger"
	"github.com/fansqz/fancode-backend/constants"
	. "github.com/fansqz/fancode-backend/service/visual_debug_servcie/debug_core/utils"
	"github.com/fansqz/fancode-backend/utils"
	"github.com/fansqz/fancode-backend/utils/gosync"
	"github.com/google/go-dap"
	"github.com/sirupsen/logrus"
)

// docker 调试器
type debugger struct {
	// dap服务端容器
	docker DockerClient
	// dap客户端
	dapCli DapClient

	// 容器映射端口
	port string

	// 启动配置
	option *Option

	// 用户工作目录
	workPath string
	// 用户文件，只支持单文件调试，所以这里是main文件
	compileFile string
	// 调试文件
	execFile string

	// 事件产生时，触发该回调
	callback NotificationCallback

	timeoutManager *TimeoutManager

	// send 发送同步命令的时候使用pending管道返回数据
	pending map[uint]chan interface{}
	// gdb序列化
	sequence uint
	mutex    sync.RWMutex

	// 由于为了防止stepIn操作会进入系统依赖内部的特殊处理
	preAction               constants.StepType // 记录gdb上一个命令
	skipContinuedEventCount int64              //记录需要跳过continue事件的数量，读写时必须加锁
}

// NewDebugger 创建docker容器
func NewDebugger() Debugger {
	return &debugger{
		pending:        make(map[uint]chan interface{}),
		sequence:       1,
		timeoutManager: NewTimeoutManager(),
	}
}

func (d *debugger) Start(ctx context.Context, option *Option) error {
	d.option = option

	if !d.checkLanguage(option.Language) {
		return fmt.Errorf("not support language %s", option.Language)
	}

	// 分配端口
	if err := d.allocatePort(); err != nil {
		logger.WithCtx(ctx).Errorf("allocate port fail, err = %v", err)
		return err
	}

	// 创建docker
	docker, err := NewDockerClient(ctx, &Config{
		ImageName: option.DebuggerImage,
		Memory:    option.MemoryLimit,
		// 限制 CPU 份额
		CPUQuota: option.CPUQuota,
		// 端口映射
		PortMapping: [][]string{{"8080", d.port}},
	})
	if err != nil {
		logger.WithCtx(ctx).Errorf("[Start] NewDockerClient fail, err = %s", err)
		return err
	}
	d.docker = docker

	// 复制代码到容器
	if err := d.saveCode(ctx); err != nil {
		logger.WithCtx(ctx).Errorf("save code fail, err = %v", err)
		return err
	}

	// 设置超时
	d.callback = func(i any) {
		option.Callback(i)
		if _, ok := i.(*ExitedEvent); !ok {
			// 重置计时器
			if d.timeoutManager != nil {
				d.timeoutManager.Reset()
			}
		}
	}

	// 开启计时器，超时未响应关闭调试
	d.timeoutManager.Start(ctx, option.DebugTimeout, func() {
		if err := d.Terminate(ctx); err != nil {
			logger.WithCtx(ctx).Errorf("[timeoutManager] terminate fail, err = %s", err)
		}
		d.callback(NewExitedEvent(-1, "调试超时，程序被关闭"))
	})

	// 启动调试
	gosync.Go(ctx, func(ctx context.Context) {
		d.start(ctx)
	})
	return nil
}

func (d *debugger) checkLanguage(language constants.LanguageType) bool {
	for _, l := range constants.SupportLanguages {
		if l == language {
			return true
		}
	}
	return false
}

func (d *debugger) allocatePort() error {
	port, err := DebugPortManager.GetPort()
	if err != nil {
		return err
	}
	d.port = strconv.Itoa(port)
	return nil
}

func (d *debugger) saveCode(ctx context.Context) error {
	// 生成容器内的工作目录路径
	d.workPath = getExecutePath(d.option.TempDir)

	// 在容器内创建工作目录
	if _, err, _ := d.docker.Exec(ctx, []string{"mkdir", "-p", d.workPath}); err != nil {
		return err
	}

	// 获取文件名
	filename, err := getMainFileNameByLanguage(d.option.Language)
	if err != nil {
		return err
	}

	// 复制代码到容器
	if err := d.docker.CopyToContainer(ctx, []byte(d.option.Code), d.workPath, filename); err != nil {
		return err
	}

	d.compileFile = path.Join(d.workPath, filename)
	return nil
}

func (d *debugger) start(ctx context.Context) {
	// 进行编译
	compileResult := d.compile(ctx, []string{d.compileFile})
	if !compileResult.Success {
		d.callback(compileResult)
		if err := d.Terminate(ctx); err != nil {
			logger.WithCtx(ctx).Errorf("[timeoutManager] terminate fail, err = %s", err)
		}
		return
	} else {
		d.callback(compileResult)
	}
	logger.WithCtx(ctx).Info("[start] compile success")

	// 启动调试器
	if err := d.startDap(ctx); err != nil {
		logger.WithCtx(ctx).Errorf("[start] startDap fail, err = %s", err)
		d.callback(NewExitedEvent(-1, "调试启动失败"))
		return
	}
	logger.WithCtx(ctx).Info("[start] startDap success")

	// 发送初始化请求
	if err := d.sendInitRequest(); err != nil {
		logger.WithCtx(ctx).Errorf("[start] init fail, err = %s", err)
	}

	// 初始断点
	if _, err := d.SetBreakpoints(ctx, d.option.BreakPoints); err != nil {
		logger.WithCtx(ctx).Errorf("[start] add break point fail, err = %s", err)
	}

	// 配置完成
	if err := d.sendConfigurationDoneRequest(); err != nil {
		logger.WithCtx(ctx).Errorf("[start] send config done fail, err = %s", err)
	}

	// 读取用户输出
	d.processUserOutput(ctx)
}

// compile 编译
// 返回编译以后的结果文件路径
func (d *debugger) compile(ctx context.Context, compileFiles []string) *CompileEvent {
	execFile := path.Join(d.workPath, "main")
	var output string
	var err error
	var exitCode int

	switch d.option.Language {
	case constants.LanguageJava:
		// java语言编译麻烦，特殊处理
		return d.compileJava(ctx, compileFiles)
	case constants.LanguageGo:
		// 2. 完整排查：查看目录下所有文件（含隐藏文件、详细信息）
		output, err, exitCode = d.docker.Exec(ctx, []string{"ls", d.workPath})
		if err != nil {
			logger.WithCtx(ctx).Errorf("[compile] 查看目录文件失败, dir = %s, err = %s", d.workPath, err)
			break
		}
		logger.WithCtx(ctx).Infof("[compile] 目录 %s 下的完整文件列表：%s", d.workPath, output)
		output, err, exitCode = d.docker.Exec(ctx, []string{"go", "-C", d.workPath, "mod", "init", "main"})
		if err != nil {
			logger.WithCtx(ctx).Errorf("[compile] exec fail, err = %s", err)
			break
		}
		output, err, exitCode = d.docker.Exec(ctx, []string{"go", "-C", d.workPath, "mod", "tidy"})
		if err != nil {
			logger.WithCtx(ctx).Errorf("[compile] exec fail, err = %s", err)
			break
		}
		output, err, exitCode = d.docker.Exec(ctx, []string{"go", "-C", d.workPath, "build", "-gcflags", "all=-N -l", "-o", execFile})
	case constants.LanguageC:
		output, err, exitCode = d.docker.Exec(ctx, append([]string{"gcc", "-g", "-o", execFile}, compileFiles...))
	case constants.LanguageCPP:
		output, err, exitCode = d.docker.Exec(ctx, append([]string{"g++", "-g", "-O0", "-o", execFile}, compileFiles...))
	}

	// 在容器中执行编译命令
	if err != nil {
		return NewCompileEvent(false, err.Error())
	}
	if exitCode == 0 {
		d.execFile = execFile
		return NewCompileEvent(true, "编译成功")
	}

	// 正则处理路径
	compileOutput := string(output)
	outputs := strings.Split(compileOutput, "# command-line-arguments\n")
	if len(outputs) > 1 {
		compileOutput = outputs[1]
	}
	re := regexp.MustCompile(`(?:\.\./|/).*?/(main\.go:\d+:\d+):`)
	// 进行替换
	compileOutput = re.ReplaceAllString(compileOutput, "$1:")
	return NewCompileEvent(false, strings.Replace(compileOutput, d.workPath, "", -1))
}

func (d *debugger) compileJava(ctx context.Context, compileFiles []string) *CompileEvent {
	execFile := path.Join(d.workPath, "main")
	// 读取main文件，规定第一个文件时main文件
	mainClass := compileFiles[0][strings.LastIndex(compileFiles[0], "/")+1:]
	mainClass = strings.Split(mainClass, ".")[0]

	// 在容器内创建存放class文件的目录
	classPath := path.Join(d.workPath, "classPath")
	if _, err, _ := d.docker.Exec(ctx, []string{"mkdir", "-p", classPath}); err != nil {
		return NewCompileEvent(false, err.Error())
	}

	// 编译为class文件
	output, err, exitCode := d.docker.Exec(ctx, append([]string{"javac", "-encoding", "UTF-8", "-d", classPath}, compileFiles...))
	if err != nil {
		return NewCompileEvent(false, err.Error())
	}
	if exitCode != 0 {
		d.execFile = execFile
		return NewCompileEvent(false, output)
	}

	// 将 MANIFEST.MF 复制到容器内
	manifestContent := "Manifest-Version: 1.0\nMain-Class: " + mainClass + "\nBuilt-By: fancode\n"
	if err := d.docker.CopyToContainer(ctx, []byte(manifestContent), classPath, "MANIFEST.MF"); err != nil {
		return NewCompileEvent(false, err.Error())
	}

	// 打包成jar包
	output, err, exitCode = d.docker.Exec(ctx, []string{"jar", "cvfm", execFile, path.Join(classPath, "MANIFEST.MF"),
		"-C", classPath, "."})
	if err != nil {
		return NewCompileEvent(false, err.Error())
	}
	if exitCode != 0 {
		return NewCompileEvent(false, output)
	}
	d.execFile = execFile
	return NewCompileEvent(true, "编译成功")
}

// startDebugger 启动调试器
func (d *debugger) startDap(ctx context.Context) error {
	// 启动dap服务端
	options := container.ExecOptions{
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          []string{"sh"},
		Tty:          true,
	}
	thread, err := d.docker.GetClient().ContainerExecCreate(ctx, d.docker.GetContainerID(), options)
	if err != nil {
		return err
	}
	attach, err := d.docker.GetClient().ContainerExecAttach(ctx, thread.ID, container.ExecAttachOptions{Detach: false, Tty: true})
	if err != nil {
		return err
	}
	d.docker.SetDebugAttach(attach)

	// 发送启动调试命令
	var cmd string
	switch d.option.Language {
	case constants.LanguageGo:
		cmd = "dlv --listen=0.0.0.0:8080 --headless=true --api-version=2 --check-go-version=false --only-same-user=false exec " + d.execFile + " --\n"
	default:
		cmd = fmt.Sprintf("go-debugger -port 8080 -language %s -file %s -codeFile %s\n", d.option.Language, d.execFile, d.compileFile)
	}
	if _, err := d.docker.GetDebugAttach().Conn.Write([]byte(cmd)); err != nil {
		return err
	}
	maxRetries := 10
	retryInterval := time.Second
	reader := bufio.NewReader(d.docker.GetDebugAttach().Reader)
	for i := 0; i < maxRetries; i++ {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				time.Sleep(retryInterval)
				continue
			}
			break
		}
		if strings.Contains(line, "listening") || strings.Contains(line, "Started server") {
			break
		}
	}

	// 创建 DAP 客户端
	client, err := NewDapClient(ctx, "127.0.0.1:"+d.port, d.dapEventCallback)
	if err != nil {
		return err
	}
	d.dapCli = client
	return nil
}

// dapEventCallback 处理dap协议返回的event数据
func (d *debugger) dapEventCallback(message dap.EventMessage) {
	// 处理event事件
	switch message.(type) {
	case *dap.StoppedEvent:
		// 需要判断当前栈帧停留位置是否在调试文件以外，如果在调试文件以外需要回到调试文件
		sf, err := d.getStackTrace()
		if err != nil {
			logger.WithCtx(context.Background()).Errorf("[dapEventCallback] getStackTrace fail, err = %s", err)
		}
		if len(sf) > 0 && !strings.HasPrefix(sf[0].Source.Path, d.workPath) {
			d.mutex.Lock()
			d.skipContinuedEventCount++
			d.mutex.Unlock()
			if d.preAction == constants.StepIn {
				_ = d.stepOut(context.Background())
			} else {
				_ = d.Continue(context.Background())
			}
			break
		}
		event := message.(*dap.StoppedEvent)
		d.callback(NewStoppedEvent(constants.StoppedReasonType(event.Body.Reason)))
	case *dap.ContinuedEvent:
		d.mutex.Lock()
		if d.skipContinuedEventCount > 0 {
			d.skipContinuedEventCount--
			d.mutex.Unlock()
			break
		}
		d.mutex.Unlock()
		d.callback(NewContinuedEvent())
	case *dap.TerminatedEvent:
		d.callback(NewTerminalEvent())
	default:
		logrus.Errorf("%v", message)
	}
}

// processUserOutput 循环处理调试输出
func (d *debugger) processUserOutput(ctx context.Context) {
	buf := make([]byte, 1024)
	for {
		n, err := d.docker.GetDebugAttach().Conn.Read(buf)
		if err != nil {
			break
		}
		d.callback(NewOutputEvent(string(buf[:n])))
	}
}

// sendInitRequest 发送init请求
func (d *debugger) sendInitRequest() error {
	request := &dap.InitializeRequest{Request: *d.dapCli.NewRequest("initialize")}
	request.Arguments = dap.InitializeRequestArguments{
		PathFormat:             "path",
		LinesStartAt1:          true,
		ColumnsStartAt1:        true,
		SupportsVariableType:   true,
		SupportsVariablePaging: true,
		Locale:                 "en-us",
	}
	_, err := d.dapCli.SendWithTimeout(request, d.option.OptionTimeout)

	return err
}

// sendConfigurationDoneRequest 发送配置完成请求
func (d *debugger) sendConfigurationDoneRequest() error {
	request := &dap.ConfigurationDoneRequest{Request: *d.dapCli.NewRequest("configurationDone")}
	_, err := d.dapCli.SendWithTimeout(request, d.option.OptionTimeout)
	return err
}

func (d *debugger) StepOver(ctx context.Context) error {
	request := &dap.NextRequest{Request: *d.dapCli.NewRequest("next")}
	request.Arguments.ThreadId = 1
	message, err := d.dapCli.SendWithTimeout(request, d.option.OptionTimeout)
	if err != nil {
		logger.WithCtx(ctx).Errorf("[StepOver] next fail, err = %v", err)
		return err
	}
	if resp, ok := message.(*dap.ErrorResponse); ok {
		return fmt.Errorf(resp.Message)
	}
	if resp, ok := message.(*dap.NextResponse); ok && !resp.Success {
		return fmt.Errorf(resp.Message)
	}
	d.preAction = constants.StepOver
	return nil
}

func (d *debugger) StepIn(ctx context.Context) error {
	if err := d.stepIn(ctx); err != nil {
		return err
	}
	d.preAction = constants.StepIn
	return nil
}

func (d *debugger) stepIn(ctx context.Context) error {
	request := &dap.StepInRequest{Request: *d.dapCli.NewRequest("stepIn")}
	request.Arguments.ThreadId = 1
	message, err := d.dapCli.SendWithTimeout(request, d.option.OptionTimeout)
	if err != nil {
		logger.WithCtx(ctx).Errorf("[StepIn] step in fail, err = %v", err)
		return err
	}
	if resp, ok := message.(*dap.ErrorResponse); ok {
		return fmt.Errorf(resp.Message)
	}
	if resp, ok := message.(*dap.StepInResponse); ok && !resp.Success {
		return fmt.Errorf(resp.Message)
	}
	return nil
}

func (d *debugger) StepOut(ctx context.Context) error {
	if err := d.stepOut(ctx); err != nil {
		return err
	}
	d.preAction = constants.StepOut
	return nil
}

func (d *debugger) stepOut(ctx context.Context) error {
	request := &dap.StepOutRequest{Request: *d.dapCli.NewRequest("stepOut")}
	request.Arguments.ThreadId = 1
	message, err := d.dapCli.SendWithTimeout(request, d.option.OptionTimeout)
	if err != nil {
		logger.WithCtx(ctx).Errorf("[StepOut] step out fail, err = %v", err)
		return err
	}
	if resp, ok := message.(*dap.ErrorResponse); ok {
		return fmt.Errorf(resp.Message)
	}
	if resp, ok := message.(*dap.StepOutResponse); ok && !resp.Success {
		return fmt.Errorf(resp.Message)
	}
	return nil
}

func (d *debugger) Continue(ctx context.Context) error {
	request := &dap.ContinueRequest{Request: *d.dapCli.NewRequest("continue")}
	request.Arguments.ThreadId = 1
	message, err := d.dapCli.SendWithTimeout(request, d.option.OptionTimeout)
	if err != nil {
		logger.WithCtx(ctx).Errorf("[Continue] continue fail, err = %v", err)
	}
	if resp, ok := message.(*dap.ErrorResponse); ok {
		return fmt.Errorf(resp.Message)
	}
	resp, ok := message.(*dap.ContinueResponse)
	if !ok {
		return fmt.Errorf("unknown continue response")
	}
	if resp.Success {
		d.callback(NewContinuedEvent())
	}
	return err
}

func (d *debugger) SetBreakpoints(ctx context.Context, breakpoints []int) ([]*Breakpoint, error) {
	// 断点去重
	breakpointsSet := utils.List2set(breakpoints)
	breakpoints = []int{}
	for _, v := range breakpointsSet.Values() {
		breakpoints = append(breakpoints, v.(int))
	}
	request := &dap.SetBreakpointsRequest{Request: *d.dapCli.NewRequest("setBreakpoints")}
	request.Arguments = dap.SetBreakpointsArguments{
		Source: dap.Source{
			Name: filepath.Base(d.compileFile),
			Path: d.compileFile,
		},
		Breakpoints: make([]dap.SourceBreakpoint, len(breakpoints)),
	}
	for i, l := range breakpoints {
		request.Arguments.Breakpoints[i].Line = l
	}
	message, err := d.dapCli.SendWithTimeout(request, d.option.OptionTimeout)
	if err != nil {
		logger.WithCtx(ctx).Errorf("[SetBreakpoint] set breakpoint fail, err = %v", err)
		return nil, err
	}
	if resp, ok := message.(*dap.ErrorResponse); ok {
		return nil, fmt.Errorf(resp.Message)
	}
	dapResponse, ok := message.(*dap.SetBreakpointsResponse)
	if !ok {
		return nil, fmt.Errorf("unknown set breakpoints response")
	}
	answer := []*Breakpoint{}
	for _, bp := range dapResponse.Body.Breakpoints {
		answer = append(answer, &Breakpoint{
			Verified: bp.Verified,
			Line:     bp.Line,
			Message:  bp.Message,
		})
	}
	return answer, nil
}

func (d *debugger) Send(ctx context.Context, input string) error {
	_, err := d.docker.GetDebugAttach().Conn.Write([]byte(input))
	if err != nil {
		logger.WithCtx(ctx).Errorf("[Send] err = %v", err)
	}
	return err
}

func (d *debugger) Terminate(ctx context.Context) error {
	// 关闭dap连接
	if d.dapCli != nil {
		d.dapCli.Close()
		d.dapCli = nil
	}
	// 容器直接关闭
	if d.docker != nil {
		if err := d.docker.RemoveContainer(ctx); err != nil {
			logger.WithCtx(ctx).Errorf("remove docker fail, docker = %v, err = %v", d.docker, err)
		}
		d.docker = nil
	}
	// 释放端口
	port, err := strconv.Atoi(d.port)
	if err != nil {
		logger.WithCtx(ctx).Errorf("[Terminate] port fail, err = %v", err)
	} else {
		DebugPortManager.ReleasePort(port)
	}

	// 重置工作路径（文件在容器内，容器删除后自动清理）
	d.workPath = ""
	d.compileFile = ""
	d.execFile = ""

	// 重置timeManager
	if d.timeoutManager != nil {
		d.timeoutManager.Cancel()
	}

	d.callback = nil
	d.pending = make(map[uint]chan interface{})
	d.sequence = 1
	return nil
}

// getExecutePath 给用户的此次运行生成一个临时目录
func getExecutePath(tempPath string) string {
	uuid := utils.GetUUID()
	executePath := path.Join(tempPath, uuid)
	return executePath
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
