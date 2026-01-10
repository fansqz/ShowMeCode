package c_debugger

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strconv"

	"github.com/fansqz/ShowMeCode/go-debugger/constants"
	. "github.com/fansqz/ShowMeCode/go-debugger/debugger"
	. "github.com/fansqz/ShowMeCode/go-debugger/debugger/gdb_debugger"
	"github.com/fansqz/ShowMeCode/go-debugger/debugger/gdb_debugger/gdb"
	"github.com/fansqz/ShowMeCode/go-debugger/debugger/utils"
	"github.com/google/go-dap"
)

type CDebugger struct {
	// 因为都是gdb调试器，所以使用c调试器即可
	gdbDebugger   *GDBDebugger
	gdbOutputUtil *GDBOutputUtil
	statusManager *utils.StatusManager
	referenceUtil *ReferenceUtil
	gdb           *gdb.Gdb
}

func NewCDebugger() *CDebugger {
	gdbDebugger := NewGDBDebugger(constants.LanguageC)
	d := &CDebugger{
		gdbDebugger:   gdbDebugger,
		gdbOutputUtil: gdbDebugger.GdbOutputUtil,
		statusManager: gdbDebugger.StatusManager,
		referenceUtil: gdbDebugger.ReferenceUtil,
	}
	return d
}

func (c *CDebugger) Start(option *StartOption) error {
	err := c.gdbDebugger.Start(option)
	c.gdb = c.gdbDebugger.GDB
	return err
}

// Run 同步方法，开始运行
func (c *CDebugger) Run() error {
	return c.gdbDebugger.Run()
}

func (c *CDebugger) StepOver() error {
	return c.gdbDebugger.StepOver()
}

func (c *CDebugger) StepIn() error {
	return c.gdbDebugger.StepIn()
}

func (c *CDebugger) StepOut() error {
	return c.gdbDebugger.StepOut()
}

func (c *CDebugger) Continue() error {
	return c.gdbDebugger.Continue()
}

func (c *CDebugger) SetBreakpoints(source dap.Source, breakpoints []dap.SourceBreakpoint) error {
	return c.gdbDebugger.SetBreakpoints(source, breakpoints)
}

func (c *CDebugger) GetStackTrace() ([]dap.StackFrame, error) {
	return c.gdbDebugger.GetStackTrace()
}

func (c *CDebugger) GetScopes(frameId int) ([]dap.Scope, error) {
	return c.gdbDebugger.GetScopes(frameId)
}

func (c *CDebugger) Terminate() error {
	return c.gdbDebugger.Terminate()
}

func (c *CDebugger) GetVariables(reference int) ([]dap.Variable, error) {
	if !c.gdbDebugger.StatusManager.Is(utils.Stopped) {
		return nil, errors.New("程序未暂停变量信息")
	}
	var variables []dap.Variable
	var err error
	// 通过scope引用获取变量列表
	if c.gdbDebugger.ReferenceUtil.CheckIsGlobalScope(reference) {
		variables, err = c.getGlobalScopeVariables()
	} else if c.gdbDebugger.ReferenceUtil.CheckIsLocalScope(reference) {
		variables, err = c.getLocalScopeVariables(reference)
	} else {
		variables, err = c.getVariables(reference)
	}
	return variables, err
}

func (c *CDebugger) getLocalScopeVariables(reference int) ([]dap.Variable, error) {
	variables, err := c.gdbDebugger.GetLocalScopeVariables(reference)
	if err != nil {
		return nil, err
	}
	frameId := c.referenceUtil.GetFrameIDByLocalReference(reference)
	var answer []dap.Variable
	for _, variable := range variables {
		// 结构体类型，设置结构体引用
		if !c.gdbOutputUtil.CheckIsAddress(variable.Value) && variable.IndexedVariables != 0 {
			variable.VariablesReference, err = c.referenceUtil.CreateVariableReference(NewStructReferenceStruct(strconv.Itoa(frameId), variable.Name, variable.Type))
			if err != nil {
				return nil, err
			}
		} else if c.gdbOutputUtil.CheckIsAddress(variable.Value) && variable.IndexedVariables != 0 && variable.Type != "char *" {
			if c.gdbOutputUtil.IsShouldBeFilterAddress(variable.Value) {
				continue
			}
			address := c.gdbOutputUtil.ConvertValueToAddress(variable.Value)
			variable.Value = address
			if !c.gdbOutputUtil.IsNullPoint(address) {
				variable.VariablesReference, err = c.referenceUtil.CreateVariableReference(NewPointReferenceStruct(variable.Type, address))
				if err != nil {
					return nil, err
				}
			}
		}
		// 如果是数组类型，设置value为数组的首地址
		addr, err := c.checkAndSetArrayAddress(variable)
		if err != nil {
			log.Printf("checkAndSetArrayAddress failed: %v\n", err)
		} else if addr != "" {
			variable.Value = addr
		}
		answer = append(answer, variable)
	}
	return answer, nil
}

func (c *CDebugger) getGlobalScopeVariables() ([]dap.Variable, error) {
	variables, err := c.gdbDebugger.GetGlobalVariables()
	if err != nil {
		return nil, err
	}
	var answer []dap.Variable
	// 遍历所有的answer
	for _, variable := range variables {
		// 结构体类型，创建结构体引用
		if !c.gdbOutputUtil.CheckIsAddress(variable.Value) && variable.IndexedVariables != 0 {
			variable.VariablesReference, err = c.referenceUtil.CreateVariableReference(NewStructReferenceStruct("0", variable.Name, variable.Type))
			if err != nil {
				return nil, err
			}
			variable.Value = ""
		} else if c.gdbOutputUtil.CheckIsAddress(variable.Value) && variable.IndexedVariables != 0 && variable.Type != "char *" {
			if c.gdbOutputUtil.IsShouldBeFilterAddress(variable.Value) {
				continue
			}
			address := c.gdbOutputUtil.ConvertValueToAddress(variable.Value)
			variable.Value = address
			if !c.gdbOutputUtil.IsNullPoint(address) {
				variable.VariablesReference, err = c.referenceUtil.CreateVariableReference(NewPointReferenceStruct(variable.Type, address))
				if err != nil {
					return nil, err
				}
			}
		}

		// 如果是数组类型，设置value为数组的首地址
		addr, err := c.checkAndSetArrayAddress(variable)
		if err != nil {
			log.Printf("checkAndSetArrayAddress failed: %v\n", err)
		} else if addr != "" {
			variable.Value = addr
		}
		answer = append(answer, variable)
	}
	return answer, nil
}

func (c *CDebugger) checkAndSetArrayAddress(variable dap.Variable) (string, error) {
	if c.gdb == nil {
		return "", errors.New("gdb instance not initialized")
	}
	pattern := `\w+\s*\[\d*\]`
	re, err := regexp.Compile(pattern)
	if err != nil {
		log.Printf("checkAndSetArrayAddress fail, err = %s\n", err)
		return "", err
	}
	if re.MatchString(variable.Type) {
		// 如果是类型是数组类型，需要设置value为地址，用于数组可视化
		m, err := c.gdb.SendWithTimeout(OptionTimeout, "data-evaluate-expression", "&"+variable.Name)
		if err != nil {
			log.Printf("checkAndSetArrayAddress fail, err = %s\n", err)
			return "", err
		}
		payload := c.gdbOutputUtil.GetInterfaceFromMap(m, "payload")
		return c.gdbOutputUtil.GetStringFromMap(payload, "value"), nil
	}
	return "", nil
}

func (c *CDebugger) getVariables(reference int) ([]dap.Variable, error) {
	variables, err := c.gdbDebugger.GetVariables(reference)
	if err != nil {
		log.Printf("getVariables failed: %v\n", err)
		return nil, err
	}
	refStruct, err := c.referenceUtil.ParseVariableReference(reference)
	if err != nil {
		log.Printf("parseVariableReference failed: %v\n", err)
		return nil, err
	}
	// 解析c语言结构体，并二次处理
	answer := make([]dap.Variable, 0, 10)
	for _, variable := range variables {
		// 如果value不为指针，且chidren不为0说明是结构体类型
		if !c.gdbOutputUtil.CheckIsAddress(variable.Value) && variable.IndexedVariables != 0 {
			variable.VariablesReference, err = c.referenceUtil.CreateVariableReference(GetFieldReferenceStruct(refStruct, variable.Name))
			if err != nil {
				return nil, err
			}
		}
		// value指针且chidren不为0说明是指针类型
		if c.gdbDebugger.GdbOutputUtil.CheckIsAddress(variable.Value) && variable.IndexedVariables != 0 {
			if variable.Type != "char *" {
				if c.gdbOutputUtil.IsShouldBeFilterAddress(variable.Value) {
					continue
				}
				address := c.gdbOutputUtil.ConvertValueToAddress(variable.Value)
				variable.Value = address
				if !c.gdbOutputUtil.IsNullPoint(address) {
					variable.VariablesReference, err = c.referenceUtil.CreateVariableReference(NewPointReferenceStruct(variable.Type, address))
					if err != nil {
						return nil, err
					}
				}
			}
		}
		answer = append(answer, variable)
	}
	return answer, nil
}

// CompileCFile 开始编译文件
func CompileCFile(workPath string, code string) (string, error) {
	// 创建工作目录, 用户的临时文件
	if err := os.MkdirAll(workPath, os.ModePerm); err != nil {
		return "", err
	}

	// 保存待编译文件
	codeFile := path.Join(workPath, "main.c")
	err := os.WriteFile(codeFile, []byte(code), 777)
	if err != nil {
		return "", err
	}
	execFile := path.Join(workPath, "main")

	cmd := exec.Command("gcc", "-g", "-o", execFile, codeFile)
	_, err = cmd.Output()
	if err != nil {
		return "", err
	}
	return execFile, err
}
