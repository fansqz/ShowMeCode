package judger

import (
	"bytes"
	"context"
	"errors"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/fansqz/fancode-backend/constants"
	"github.com/fansqz/fancode-backend/service/user_coding_service/judger/cgroup"
	"github.com/fansqz/fancode-backend/utils"
)

type JudgeCore struct {
}

func NewJudgeCore() *JudgeCore {
	return &JudgeCore{}
}

// Compile 编译，编译时在容器外进行编译的
// compileFiles第个文件是main文件
func (j *JudgeCore) Compile(compileFiles []string, outFilePath string, options *CompileOptions) (*CompileResult, error) {
	result := &CompileResult{
		Compiled:         false,
		ErrorMessage:     "",
		CompiledFilePath: "",
	}

	// 进行编译
	var cmd *exec.Cmd
	var ctx context.Context
	language := j.getLanguage(options)
	switch language {
	case constants.LanguageGo:
		compileFiles = append([]string{"build", "-gcflags", "-l -N", "-o", outFilePath}, compileFiles...)
		var cancel context.CancelFunc
		if options != nil && options.LimitTime != 0 {
			ctx, cancel = context.WithTimeout(context.Background(), time.Duration(options.LimitTime))
			defer cancel()
		} else {
			ctx = context.Background()
		}
		cmd = exec.CommandContext(ctx, "go", compileFiles...)
	default:
		result.ErrorMessage = "不支持该语言\n"
		return result, nil
	}
	cmd.Stdout = &bytes.Buffer{}
	cmd.Stderr = &bytes.Buffer{}

	var err error
	if err = cmd.Start(); err == nil {
		err = cmd.Wait()
	}
	if err != nil {
		return j.setErrMessageForCompileResult(ctx, cmd, err, result, options)
	}

	result.Compiled = true
	result.CompiledFilePath = outFilePath
	return result, nil
}

func (j *JudgeCore) getLanguage(options *CompileOptions) constants.LanguageType {
	language := constants.LanguageC
	if options != nil && options.Language != "" {
		language = options.Language
	}
	return language
}

func (j *JudgeCore) setErrMessageForCompileResult(ctx context.Context, cmd *exec.Cmd, err error, result *CompileResult, options *CompileOptions) (*CompileResult, error) {
	if err != nil {
		errBytes := cmd.Stderr.(*bytes.Buffer).Bytes()
		errMessage := string(errBytes)
		if len(options.ExcludedPaths) != 0 {
			errMessage = j.maskPath(string(errBytes), options.ExcludedPaths, options.ReplacementPath)
		}
		// 如果是由于超时导致的错误，则返回自定义错误
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			result.ErrorMessage = "编译超时\n" + errMessage
			return result, nil
		}
		if len(errBytes) != 0 {
			result.ErrorMessage = errMessage
			return result, nil
		}
	}
	return nil, err
}

// maskPath 函数用于屏蔽错误消息中的路径信息
func (j *JudgeCore) maskPath(errorMessage string, excludedPaths []string, replacementPath string) string {
	if errorMessage == "" {
		return ""
	}
	// 遍历需要屏蔽的敏感路径
	for _, excludedPath := range excludedPaths {
		// 如果 excludedPath 是绝对路径，但 errorMessage 中含有相对路径 "./"，则将 "./" 替换为绝对路径
		if filepath.IsAbs(excludedPath) && filepath.IsAbs("./") {
			relativePath := "." + string(filepath.Separator)
			absolutePath := filepath.Join(excludedPath, relativePath)
			errorMessage = strings.Replace(errorMessage, relativePath, absolutePath, -1)
		}

		// 构建正则表达式，匹配包含敏感路径的错误消息
		pattern := regexp.QuoteMeta(excludedPath)
		re := regexp.MustCompile(pattern)
		errorMessage = re.ReplaceAllString(errorMessage, replacementPath)
	}

	return errorMessage
}

// Execute 运行
func (j *JudgeCore) Execute(execFile string, inputCh <-chan []byte, outputCh chan<- ExecuteResult, exitCh <-chan string, options *ExecuteOptions) error {
	//language := constants.LanguageC
	//if options != nil && options.Language != "" {
	//	language = options.Language
	//}
	//// 根据扩展名设置执行命令
	//cmdName := ""
	//cmdArg := []string{}
	//switch language {
	//case constants.LanguageC:
	//	cmdName = execFile
	//case constants.LanguageJava:
	//	cmdName = "java"
	//	cmdArg = []string{"-jar", execFile}
	//case constants.LanguageGo:
	//	cmdName = execFile
	//default:
	//	return fmt.Errorf("不支持该语言")
	//}
	//
	//// 创建cgroup限制资源
	//cgroup, err := j.initCGroup(options)
	//if err != nil {
	//	return err
	//}
	//
	//go func() {
	//	defer func() {
	//		if err := recover(); err != nil {
	//			fmt.Println("defer", err)
	//		}
	//	}()
	//	for {
	//		select {
	//		case inputItem := <-inputCh:
	//
	//			// 设置超时上下文
	//			var ctx context.Context
	//			var cancel context.CancelFunc
	//			if options != nil && options.LimitTime != 0 {
	//				ctx, cancel = context.WithTimeout(context.Background(), time.Duration(options.LimitTime)+100)
	//				defer cancel()
	//			} else {
	//				ctx = context.Background()
	//			}
	//
	//			// 创建子进程，并将其加入cgroup
	//			cmd2 := exec.CommandContext(ctx, cmdName, cmdArg...)
	//			result := ExecuteResult{}
	//
	//			cmd2.Stdin = bytes.NewReader(inputItem)
	//			cmd2.Stdout = &bytes.Buffer{}
	//			cmd2.Stderr = &bytes.Buffer{}
	//			cmd2.SysProcAttr = &syscall.SysProcAttr{
	//				Setpgid: true,
	//			}
	//
	//			beginTime := time.Now()
	//			if err = cmd2.Start(); err != nil {
	//				result.Executed = false
	//				result.ErrorMessage = err.Error() + "\n"
	//				outputCh <- result
	//				break
	//			}
	//
	//			// 将进程写入cgroup组
	//			if err = cgroup.AddPID(cmd2.Process.Pid); err != nil {
	//				result.Executed = false
	//				result.ErrorMessage = err.Error() + "\n"
	//				outputCh <- result
	//				break
	//			}
	//
	//			// 等待程序执行
	//			cmd2.Wait()
	//			// 读取使用cpu和内存，以及执行时间
	//			rusage := cmd2.ProcessState.SysUsage().(*syscall.Rusage)
	//			result.UsedCpuTime = rusage.Utime.Sec*1000 + rusage.Utime.Usec/1000
	//			result.UsedMemory = rusage.Maxrss * 1024
	//			result.UsedTime = int64(time.Now().Sub(beginTime))
	//
	//			// 输出的错误信息
	//			errMessage := string(cmd2.Stderr.(*bytes.Buffer).Bytes())
	//			if options != nil && len(options.ExcludedPaths) != 0 {
	//				errMessage = j.maskPath(errMessage, options.ExcludedPaths, options.ReplacementPath)
	//			}
	//			outMessage := cmd2.Stdout.(*bytes.Buffer).Bytes()
	//			// 检测内存占用，cpu占用，以及执行时间
	//			if options != nil && options.LimitTime < result.UsedTime {
	//				result.Executed = false
	//				result.ErrorMessage = "运行超时\n"
	//			} else if options != nil && options.MemoryLimit < result.UsedMemory {
	//				result.Executed = false
	//				result.ErrorMessage = "内存超出限制\n"
	//			} else if options != nil && options.CPUQuota < result.UsedCpuTime {
	//				result.Executed = false
	//				result.ErrorMessage = "cpu超出限制\n"
	//			} else if len(errMessage) != 0 {
	//				result.Executed = false
	//				result.ErrorMessage = errMessage
	//			} else {
	//				result.Executed = true
	//				result.Output = outMessage
	//			}
	//			outputCh <- result
	//		case <-exitCh:
	//			cgroup.Release()
	//			return
	//		}
	//	}
	//}()

	return nil
}

func (j *JudgeCore) initCGroup(options *ExecuteOptions) (*cgroup.CGroup, error) {
	// 创建cgroup限制资源
	cgroup, err := cgroup.NewCGroup(utils.GetUUID())
	if err != nil {
		return nil, err
	}
	if options != nil && options.MemoryLimit != 0 {
		if err = cgroup.SetMemoryLimit(options.MemoryLimit); err != nil {
			return nil, err
		}
	}
	if options != nil && options.CPUQuota != 0 {
		if err = cgroup.SetCPUQuota(options.CPUQuota); err != nil {
			return nil, err
		}
	}
	return cgroup, nil
}
