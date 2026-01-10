package c_debugger

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"testing"

	"github.com/google/go-dap"

	"github.com/fansqz/go-debugger/debugger"
	"github.com/fansqz/go-debugger/utils"
	"github.com/stretchr/testify/assert"
)

// testHelper 测试辅助结构体,封装测试所需的通用组件
type testHelper struct {
	t        *testing.T
	workPath string
	debug    *CDebugger
	eventCh  chan dap.EventMessage
}

// newTestHelper 创建新的测试辅助实例
func newTestHelper(t *testing.T) *testHelper {
	workPath := path.Join("/var/fanCode/tempDir", utils.GetUUID())
	debug := NewCDebugger()
	eventCh := make(chan dap.EventMessage, 10)

	return &testHelper{
		t:        t,
		workPath: workPath,
		debug:    debug,
		eventCh:  eventCh,
	}
}

// setup 设置测试环境
func (h *testHelper) setup(cFile string) {
	// 编译测试文件
	execFile, code, err := compileFile(h.workPath, cFile)
	assert.Nil(h.t, err)

	// 启动调试器
	err = h.debug.Start(&debugger.StartOption{
		MainCode: code,
		ExecFile: execFile,
		Callback: func(message dap.EventMessage) {
			h.eventCh <- message
		},
	})
	assert.Nil(h.t, err)
}

// cleanup 清理测试环境
func (h *testHelper) cleanup() {
	os.RemoveAll(h.workPath)
}

// waitForEvent 等待并验证事件
func (h *testHelper) waitForEvent(expectedEvent string) dap.EventMessage {
	event := <-h.eventCh
	assert.Equal(h.t, expectedEvent, event.GetEvent().Event)
	return event
}

// TestDebug 测试普通调试功能
func TestDebug(t *testing.T) {
	helper := newTestHelper(t)
	defer helper.cleanup()

	helper.setup("debug.c")

	// 设置断点
	err := helper.debug.SetBreakpoints(dap.Source{Path: "main.c"}, []dap.SourceBreakpoint{
		{Line: 3},
		{Line: 7},
	})
	assert.Nil(t, err)

	// 启动调试并验证初始事件
	err = helper.debug.Run()
	assert.Nil(t, err)
	helper.waitForEvent("continued")
	helper.waitForEvent("stopped")
	assert.Equal(t, 3, getStoppedLine(helper.debug))

	// 测试单步执行
	err = helper.debug.StepOver()
	assert.Nil(t, err)
	helper.waitForEvent("continued")
	helper.waitForEvent("stopped")
	assert.Equal(t, 5, getStoppedLine(helper.debug))

	// 测试继续执行
	err = helper.debug.Continue()
	assert.Nil(t, err)
	helper.waitForEvent("continued")

	// 模拟用户输入
	r, w, _ := os.Pipe()
	os.Stdin = r
	io.WriteString(w, "10\n")
	helper.waitForEvent("stopped")

	// 验证程序结束
	err = helper.debug.Continue()
	helper.waitForEvent("continued")
	helper.waitForEvent("terminated")
}

// TestVariable 测试变量获取
func TestVariable(t *testing.T) {
	helper := newTestHelper(t)
	defer helper.cleanup()

	helper.setup("variable.c")

	// 设置断点
	err := helper.debug.SetBreakpoints(dap.Source{Path: "main.c"}, []dap.SourceBreakpoint{
		{Line: 55},
		{Line: 76},
	})
	assert.Nil(t, err)

	// 启动调试
	err = helper.debug.Run()
	assert.Nil(t, err)
	helper.waitForEvent("continued")
	helper.waitForEvent("stopped")

	// 获取并验证栈帧信息
	stacks, err := helper.debug.GetStackTrace()
	assert.Nil(t, err)

	// 验证作用域
	scopes, _ := helper.debug.GetScopes(stacks[0].Id)
	assert.Equal(t, []dap.Scope{
		{Name: "Global", VariablesReference: 1001},
		{Name: "Local", VariablesReference: 1002},
	}, scopes)

	// 验证全局变量
	verifyGlobalVariables(t, helper.debug, scopes[0].VariablesReference)

	// 验证局部变量
	verifyLocalVariables(t, helper.debug, scopes[1].VariablesReference)

	// 继续执行到下一个断点
	err = helper.debug.Continue()
	assert.Nil(t, err)
	helper.waitForEvent("continued")
	helper.waitForEvent("stopped")

	verifyLocalVariables2(t, helper.debug, scopes[1].VariablesReference)
}

// verifyGlobalVariables 验证全局变量
func verifyGlobalVariables(t *testing.T, debug *CDebugger, ref int) {
	variables, err := debug.GetVariables(ref)
	assert.Nil(t, err)

	// 验证基本类型变量
	assert.Equal(t, []dap.Variable{
		{Name: "globalChar", Value: "65 'A'", Type: "char"},
		{Name: "globalFloat", Value: "3.1400001", Type: "float"},
		{Name: "globalInt", Value: "10", Type: "int"},
		{Name: "globalItem", Value: "", Type: "Item", VariablesReference: 1100, IndexedVariables: 3},
	}, variables[0:4])

	// 验证静态全局变量
	assert.Equal(t, []dap.Variable{{Name: "staticGlobalInt", Value: "20", Type: "int"}}, variables[5:6])

	// 验证指针变量
	assert.Equal(t, "globalItemPtr", variables[4].Name)
	assert.Equal(t, "Item *", variables[4].Type)
	assert.Equal(t, 1101, variables[4].VariablesReference)
}

// verifyLocalVariables 验证局部变量
func verifyLocalVariables(t *testing.T, debug *CDebugger, ref int) {
	variables, err := debug.GetVariables(ref)
	assert.Nil(t, err)

	assert.Equal(t, []dap.Variable{
		{Name: "argint", Value: "2", Type: "int"},
		{Name: "localInt", Value: "5", Type: "int"},
		{Name: "localChar", Value: "71 'G'", Type: "char"},
		{Name: "staticLocalFloat", Value: "6.78000021", Type: "float"},
		{Name: "localItem", Value: "{...}", Type: "Item", VariablesReference: 1102, IndexedVariables: 3},
		{Name: "localColor", Value: "BLUE", Type: "Color"},
		{Name: "localValue", Value: "{...}", Type: "Value", VariablesReference: 1103, IndexedVariables: 3},
	}, variables)
}

// verifyLocalVariables2 测试一下数组
func verifyLocalVariables2(t *testing.T, debug *CDebugger, ref int) {
	variables, _ := debug.GetVariables(ref)
	variables, err := debug.GetVariables(variables[4].VariablesReference)
	assert.Nil(t, err)
	assert.NotEqual(t, 0, len(variables))
}

// TestLink 测试链表算法
func TestLink(t *testing.T) {
	helper := newTestHelper(t)
	defer helper.cleanup()

	helper.setup("link.c")

	// 设置断点
	err := helper.debug.SetBreakpoints(dap.Source{Path: "main.c"}, []dap.SourceBreakpoint{
		{Line: 28},
	})
	assert.Nil(t, err)

	// 启动调试
	err = helper.debug.Run()
	assert.Nil(t, err)
	helper.waitForEvent("continued")
	helper.waitForEvent("stopped")

	// 验证链表结构
	stacks, err := helper.debug.GetStackTrace()
	assert.Nil(t, err)
	scopes, err := helper.debug.GetScopes(stacks[0].Id)
	assert.Nil(t, err)

	// 验证链表节点
	nodes, err := helper.debug.GetVariables(scopes[1].VariablesReference)
	assert.Nil(t, err)
	next, err := helper.debug.GetVariables(nodes[0].VariablesReference)
	assert.Nil(t, err)
	next, err = helper.debug.GetVariables(next[1].VariablesReference)
	assert.Nil(t, err)
	next, err = helper.debug.GetVariables(next[1].VariablesReference)
	assert.NotNil(t, next)
}

// TestStruct 测试结构体
func TestStruct(t *testing.T) {
	helper := newTestHelper(t)
	defer helper.cleanup()

	helper.setup("struct.c")

	// 设置断点
	err := helper.debug.SetBreakpoints(dap.Source{Path: "main.c"}, []dap.SourceBreakpoint{
		{Line: 39},
	})
	assert.Nil(t, err)

	// 启动调试
	err = helper.debug.Run()
	assert.Nil(t, err)
	helper.waitForEvent("continued")
	helper.waitForEvent("stopped")

	// 验证结构体
	stacks, err := helper.debug.GetStackTrace()
	assert.Nil(t, err)
	scopes, err := helper.debug.GetScopes(stacks[0].Id)
	assert.Nil(t, err)

	// 验证全局和局部结构体
	globals, err := helper.debug.GetVariables(scopes[0].VariablesReference)
	assert.Nil(t, err)
	locals, err := helper.debug.GetVariables(scopes[1].VariablesReference)
	assert.Nil(t, err)

	verifyStudentStruct(t, helper.debug, globals[0])
	verifyStudentStruct(t, helper.debug, globals[1])
	verifyStudentStruct(t, helper.debug, locals[0])
	verifyStudentStruct(t, helper.debug, locals[1])
}

// verifyStudentStruct 验证学生结构体
func verifyStudentStruct(t *testing.T, debug *CDebugger, variable dap.Variable) {
	children, err := debug.GetVariables(variable.VariablesReference)
	assert.Nil(t, err)

	// 验证基本字段
	assert.Equal(t, "name", children[0].Name)
	assert.Equal(t, "id", children[1].Name)
	assert.Equal(t, "birthdate", children[2].Name)

	// 验证名字数组
	studentName, err := debug.GetVariables(children[0].VariablesReference)
	assert.Equal(t, 50, len(studentName))

	// 验证生日结构体
	studentBirthdate, err := debug.GetVariables(children[2].VariablesReference)
	assert.Equal(t, []dap.Variable{
		{Name: "year", Value: "2005", Type: "int"},
		{Name: "month", Value: "3", Type: "int"},
		{Name: "day", Value: "15", Type: "int"},
	}, studentBirthdate)
}

// getStoppedLine 获取当前停止的行号
func getStoppedLine(gdb debugger.Debugger) int {
	stackTrace, _ := gdb.GetStackTrace()
	if len(stackTrace) != 0 {
		return stackTrace[0].Line
	}
	return 0
}

// compileFile 编译C源文件
func compileFile(workPath string, cFile string) (string, string, error) {
	// 创建工作目录
	if err := os.MkdirAll(workPath, os.ModePerm); err != nil {
		return "", "", err
	}

	// 保存源文件
	code, err := os.ReadFile(path.Join("./test_file", cFile))
	if err != nil {
		return "", "", err
	}
	execFile, err := CompileCFile(workPath, string(code))

	return execFile, string(code), err
}

// TestMatrix 二维数组测试
func TestMatrix(t *testing.T) {
	helper := newTestHelper(t)
	defer helper.cleanup()
	helper.setup("matrix.c")
	// 设置断点
	err := helper.debug.SetBreakpoints(dap.Source{Path: "main.c"}, []dap.SourceBreakpoint{
		{Line: 20},
	})
	assert.Nil(t, err)
	err = helper.debug.Run()
	assert.Nil(t, err)
	helper.waitForEvent("continued")
	helper.waitForEvent("stopped")

	// 验证全局和局部结构体
	matrix, err := getTargetLocalsVariable(t, helper.debug, "matrix")
	assert.Nil(t, err)
	matrixVariables, err := helper.debug.GetVariables(matrix.VariablesReference)
	fmt.Sprintln(matrixVariables)
	assert.Nil(t, err)
	array1, err := helper.debug.GetVariables(matrixVariables[0].VariablesReference)
	fmt.Sprintln(array1)
}

func getTargetLocalsVariable(t *testing.T, debugger debugger.Debugger, name string) (dap.Variable, error) {
	stacks, err := debugger.GetStackTrace()
	assert.Nil(t, err)
	scopes, err := debugger.GetScopes(stacks[0].Id)
	assert.Nil(t, err)
	locals, err := debugger.GetVariables(scopes[1].VariablesReference)
	assert.Nil(t, err)
	for _, v := range locals {
		if v.Name == name {
			return v, nil
		}
	}
	return dap.Variable{}, errors.New("not found")
}
