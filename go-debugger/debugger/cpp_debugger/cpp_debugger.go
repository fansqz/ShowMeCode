package cpp_debugger

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/fansqz/ShowMeCode/go-debugger/constants"
	. "github.com/fansqz/ShowMeCode/go-debugger/debugger"
	. "github.com/fansqz/ShowMeCode/go-debugger/debugger/gdb_debugger"
	"github.com/fansqz/ShowMeCode/go-debugger/debugger/gdb_debugger/gdb"
	"github.com/fansqz/ShowMeCode/go-debugger/debugger/utils"
	"github.com/google/go-dap"
	"github.com/sirupsen/logrus"
)

// CPPDebugger C++调试器，基于GDB实现
// 主要处理C++特有的调试需求，如智能指针、STL容器、访问修饰符等
type CPPDebugger struct {
	gdbDebugger   *GDBDebugger         // 底层GDB调试器
	gdbOutputUtil *GDBOutputUtil       // GDB输出解析工具
	statusManager *utils.StatusManager // 调试状态管理器
	referenceUtil *ReferenceUtil       // 变量引用工具
	gdb           *gdb.Gdb             // GDB实例
}

// NewCPPDebugger 创建新的C++调试器实例
func NewCPPDebugger() *CPPDebugger {
	gdbDebugger := NewGDBDebugger(constants.LanguageCpp)
	return &CPPDebugger{
		gdbDebugger:   gdbDebugger,
		gdbOutputUtil: gdbDebugger.GdbOutputUtil,
		statusManager: gdbDebugger.StatusManager,
		referenceUtil: gdbDebugger.ReferenceUtil,
	}
}

// Start 启动调试器，加载可执行文件
func (c *CPPDebugger) Start(option *StartOption) error {
	err := c.gdbDebugger.Start(option)
	c.gdb = c.gdbDebugger.GDB
	return err
}

// 代理方法 - 直接调用底层GDB调试器
func (c *CPPDebugger) Run() error                               { return c.gdbDebugger.Run() }
func (c *CPPDebugger) StepOver() error                          { return c.gdbDebugger.StepOver() }
func (c *CPPDebugger) StepIn() error                            { return c.gdbDebugger.StepIn() }
func (c *CPPDebugger) StepOut() error                           { return c.gdbDebugger.StepOut() }
func (c *CPPDebugger) Continue() error                          { return c.gdbDebugger.Continue() }
func (c *CPPDebugger) Terminate() error                         { return c.gdbDebugger.Terminate() }
func (c *CPPDebugger) GetStackTrace() ([]dap.StackFrame, error) { return c.gdbDebugger.GetStackTrace() }
func (c *CPPDebugger) GetScopes(frameId int) ([]dap.Scope, error) {
	return c.gdbDebugger.GetScopes(frameId)
}

// SetBreakpoints 设置断点
func (c *CPPDebugger) SetBreakpoints(source dap.Source, breakpoints []dap.SourceBreakpoint) error {
	return c.gdbDebugger.SetBreakpoints(source, breakpoints)
}

// GetVariables 获取变量列表，根据引用类型分发到不同的处理方法
// C++调试需要特殊处理，因为可能存在public、private等访问修饰符
func (c *CPPDebugger) GetVariables(reference int) ([]dap.Variable, error) {
	if !c.statusManager.Is(utils.Stopped) {
		return nil, errors.New("程序未暂停变量信息")
	}

	switch {
	case c.referenceUtil.CheckIsGlobalScope(reference):
		return c.getGlobalScopeVariables()
	case c.referenceUtil.CheckIsLocalScope(reference):
		return c.getLocalScopeVariables(reference)
	default:
		return c.getVariables(reference)
	}
}

// CompileCPPFile 编译C++文件
// 使用G++编译器，启用调试信息和优化选项
func CompileCPPFile(workPath string, code string) (string, error) {
	// 创建工作目录
	if err := os.MkdirAll(workPath, os.ModePerm); err != nil {
		return "", err
	}

	// 保存源代码文件
	codeFile := path.Join(workPath, "main.cpp")
	if err := os.WriteFile(codeFile, []byte(code), 0777); err != nil {
		return "", err
	}

	// 编译选项说明：
	// -g: 生成调试信息
	// -O0: 禁用优化，便于调试
	// -fno-inline-functions: 禁用函数内联
	// -ftrivial-auto-var-init=zero: 自动初始化变量为0
	// -fsanitize=undefined: 启用未定义行为检测
	// -fno-omit-frame-pointer: 保留帧指针
	// -fno-reorder-blocks-and-partition: 禁用块重排序
	// -fvar-tracking-assignments: 启用变量跟踪
	execFile := path.Join(workPath, "main")
	cmd := exec.Command("g++", "-g", "-O0",
		"-fno-inline-functions",
		"-ftrivial-auto-var-init=zero", "-fsanitize=undefined", "-fno-omit-frame-pointer",
		"-fno-reorder-blocks-and-partition", "-fvar-tracking-assignments", codeFile, "-o", execFile)

	if _, err := cmd.Output(); err != nil {
		return "", err
	}
	return execFile, nil
}

// processVariable 处理单个变量的通用逻辑
// 处理结构体、指针、数组等不同类型的变量
func (c *CPPDebugger) processVariable(variable dap.Variable, frameId string) dap.Variable {
	var err error
	// 结构体类型设置结构体引用
	if !c.gdbOutputUtil.CheckIsAddress(variable.Value) && variable.IndexedVariables != 0 {
		variable.VariablesReference, err = c.referenceUtil.CreateVariableReference(NewStructReferenceStruct(frameId, variable.Name, variable.Type))
		if err != nil {
			logrus.Errorf("processVariable fail, err = %s\n", err)
		}
		variable.Value = "" // 结构体类型清空value，避免显示地址
	} else if c.gdbOutputUtil.CheckIsAddress(variable.Value) && variable.IndexedVariables != 0 {
		// 排除char*类型，避免处理字符串指针
		if variable.Type != "char *" && !c.gdbOutputUtil.IsShouldBeFilterAddress(variable.Value) {
			address := c.gdbOutputUtil.ConvertValueToAddress(variable.Value)
			variable.Value = address
			// 非空指针才创建引用
			if !c.gdbOutputUtil.IsNullPoint(address) {
				variable.VariablesReference, _ = c.referenceUtil.CreateVariableReference(NewPointReferenceStruct(variable.Type, address))
			}
		}
	}

	// 处理数组类型：设置value为数组的首地址，便于数组可视化
	if addr, err := c.checkAndSetArrayAddress(variable); err == nil && addr != "" {
		variable.Value = addr
	}

	return variable
}

// getGlobalScopeVariables 获取全局变量列表
func (c *CPPDebugger) getGlobalScopeVariables() ([]dap.Variable, error) {
	variables, err := c.gdbDebugger.GetGlobalVariables()
	if err != nil {
		return nil, err
	}

	result := make([]dap.Variable, 0, len(variables))
	for _, variable := range variables {
		// 全局变量的frameId为"0"
		result = append(result, c.processVariable(variable, "0"))
	}
	return result, nil
}

// getLocalScopeVariables 获取局部变量列表
func (c *CPPDebugger) getLocalScopeVariables(reference int) ([]dap.Variable, error) {
	variables, err := c.gdbDebugger.GetLocalScopeVariables(reference)
	if err != nil {
		return nil, err
	}

	frameId := c.referenceUtil.GetFrameIDByLocalReference(reference)
	result := make([]dap.Variable, 0, len(variables))
	for _, variable := range variables {
		result = append(result, c.processVariable(variable, strconv.Itoa(frameId)))
	}
	return result, nil
}

// checkAndSetArrayAddress 检查并设置数组地址
// 对于数组类型，获取其首地址用于可视化
func (c *CPPDebugger) checkAndSetArrayAddress(variable dap.Variable) (string, error) {
	// 匹配数组类型模式：如 int[10], char[20] 等
	pattern := `\w+\s*\[\d*\]`
	re, err := regexp.Compile(pattern)
	if err != nil {
		return "", fmt.Errorf("compile regex failed: %w", err)
	}

	if !re.MatchString(variable.Type) {
		return "", nil
	}

	// 使用GDB获取数组的地址
	m, err := c.gdb.SendWithTimeout(OptionTimeout, "data-evaluate-expression", "&"+variable.Name)
	if err != nil {
		return "", fmt.Errorf("get array address failed: %w", err)
	}

	payload := c.gdbOutputUtil.GetInterfaceFromMap(m, "payload")
	return c.gdbOutputUtil.GetStringFromMap(payload, "value"), nil
}

// getVariables 获取结构体或对象的成员变量
func (c *CPPDebugger) getVariables(reference int) ([]dap.Variable, error) {
	variables, err := c.getVariablesForCpp(reference)
	if err != nil {
		return nil, err
	}

	refStruct, err := c.referenceUtil.ParseVariableReference(reference)
	if err != nil {
		return nil, err
	}

	result := make([]dap.Variable, 0, len(variables))
	for _, variable := range variables {
		// 处理结构体成员：如果value不为指针且IndexedVariables不为0，说明是结构体
		if !c.gdbOutputUtil.CheckIsAddress(variable.Value) && variable.IndexedVariables != 0 {
			variable.VariablesReference, _ = c.referenceUtil.CreateVariableReference(GetFieldReferenceStruct(refStruct, variable.Name))
		} else if c.gdbOutputUtil.CheckIsAddress(variable.Value) && variable.IndexedVariables != 0 {
			if variable.Type != "char *" && !c.gdbOutputUtil.IsShouldBeFilterAddress(variable.Value) {
				address := c.gdbOutputUtil.ConvertValueToAddress(variable.Value)
				variable.Value = address

				if !c.gdbOutputUtil.IsNullPoint(address) {
					variable.VariablesReference, _ = c.referenceUtil.CreateVariableReference(NewPointReferenceStruct(variable.Type, address))
				}
			}
		}

		result = append(result, variable)
	}
	return result, nil
}

// getVariablesForCpp 获取C++变量的子成员
// 由于C++的访问修饰符，需要使用特殊方法获取结构体内容
func (c *CPPDebugger) getVariablesForCpp(reference int) ([]dap.Variable, error) {
	refStruct, err := c.gdbDebugger.ReferenceUtil.ParseVariableReference(reference)
	if err != nil {
		return nil, fmt.Errorf("parse variable reference failed: %w", err)
	}

	// 切换到对应的栈帧
	if err = c.gdbDebugger.SelectFrame(refStruct); err != nil {
		return nil, err
	}

	// 创建临时变量用于获取子成员
	targetVar := "structName"
	targetVariable, err := c.gdbDebugger.CreateVar(refStruct, targetVar)
	if err != nil {
		return nil, err
	}
	defer c.gdbDebugger.DeleteVar(targetVar)

	return c.varListChildrenForCpp(refStruct, targetVariable)
}

// varListChildrenForCpp 获取C++变量的子成员列表
// 根据变量类型选择不同的处理方法
func (c *CPPDebugger) varListChildrenForCpp(ref *ReferenceStruct, targetVariable *dap.Variable) ([]dap.Variable, error) {
	if c.checkIsCppArrayType(targetVariable) {
		return c.varListChildrenForCppArray(ref, targetVariable)
	}
	return c.varListChildrenForCppStruct(ref)
}

// varListChildrenForCppStruct 获取结构体的成员变量
// 使用data-evaluate-expression获取结构体内容，避免访问修饰符问题
func (c *CPPDebugger) varListChildrenForCppStruct(ref *ReferenceStruct) ([]dap.Variable, error) {
	exp := c.GetExport(ref)
	_, _ = c.gdbDebugger.GDB.SendWithTimeout(OptionTimeout, "enable-pretty-printing")

	// 获取结构体的字符串表示
	m, err := c.gdbDebugger.GDB.SendWithTimeout(OptionTimeout, "data-evaluate-expression", exp)
	if err != nil {
		return nil, fmt.Errorf("evaluate expression failed: %w", err)
	}

	payload := c.gdbOutputUtil.GetInterfaceFromMap(m, "payload")
	value := c.gdbOutputUtil.GetStringFromMap(payload, "value")
	keys := c.parseObject2Keys(value)

	result := make([]dap.Variable, 0, len(keys))
	for _, key := range keys {
		// 为每个成员创建临时变量
		m, err = c.gdb.SendWithTimeout(OptionTimeout, "var-create", key, "*", fmt.Sprintf("(%s).%s", exp, key))
		if err != nil {
			log.Printf("var-create failed for key %s: %v", key, err)
			continue
		}

		if variable, ok := c.gdbOutputUtil.ParseVarCreate(m); ok {
			result = append(result, *variable)
		}
		// 清理临时变量
		_, _ = c.gdb.SendWithTimeout(OptionTimeout, "var-delete", key)
	}
	return result, nil
}

// varListChildrenForCppArray 获取数组的元素
// 支持std::array、C数组、std::vector等类型
func (c *CPPDebugger) varListChildrenForCppArray(ref *ReferenceStruct, targetVariable *dap.Variable) ([]dap.Variable, error) {
	arrayLength := c.getArrayLength(ref, targetVariable)
	if arrayLength == 0 {
		return nil, nil
	}

	exp := c.gdbDebugger.GetExport(ref)
	result := make([]dap.Variable, 0, arrayLength)

	// 遍历数组的每个元素
	for i := 0; i < arrayLength; i++ {
		m, err := c.gdb.SendWithTimeout(OptionTimeout, "var-create", "arrayNameChildren", "*", fmt.Sprintf("%s[%d]", exp, i))
		if err != nil {
			log.Printf("var-create failed for index %d: %v", i, err)
			continue
		}

		if variable, ok := c.gdbOutputUtil.ParseVarCreate(m); ok {
			variable.Name = strconv.Itoa(i) // 使用索引作为元素名称
			result = append(result, *variable)
		}
		_, _ = c.gdb.SendWithTimeout(OptionTimeout, "var-delete", "arrayNameChildren")
	}
	return result, nil
}

// getArrayLength 获取数组长度
// 支持多种数组类型：std::array、C数组、std::vector
func (c *CPPDebugger) getArrayLength(ref *ReferenceStruct, targetVariable *dap.Variable) int {
	// 尝试从std::array获取长度
	if length := c.extractStdArrayLength(targetVariable.Type); length > 0 {
		return length
	}

	// 尝试从C数组获取长度
	if length := c.extractCArrayLength(targetVariable.Type); length > 0 {
		return length
	}

	// 尝试从std::vector获取长度
	if strings.Contains(targetVariable.Type, "std::vector") {
		return c.getVectorLength(ref, targetVariable)
	}

	return 0
}

// extractStdArrayLength 从std::array类型中提取长度
// 匹配模式：std::array<T, N>
func (c *CPPDebugger) extractStdArrayLength(typeStr string) int {
	pattern := `std::array<[^,]+,\s*(\d+)\s*>`
	re := regexp.MustCompile(pattern)
	if match := re.FindStringSubmatch(typeStr); match != nil {
		if length, err := strconv.Atoi(match[1]); err == nil {
			return length
		}
	}
	return 0
}

// extractCArrayLength 从C数组类型中提取长度
// 匹配模式：T[N]
func (c *CPPDebugger) extractCArrayLength(typeStr string) int {
	pattern := `([\w\s\*\:]+)\[\s*(\d+)\s*\]`
	re := regexp.MustCompile(pattern)
	if match := re.FindStringSubmatch(typeStr); match != nil {
		if length, err := strconv.Atoi(match[2]); err == nil {
			return length
		}
	}
	return 0
}

// getVectorLength 获取std::vector的长度
// 优先使用size()方法，失败时使用sizeof计算
func (c *CPPDebugger) getVectorLength(ref *ReferenceStruct, targetVariable *dap.Variable) int {
	exp := c.gdbDebugger.GetExport(ref)

	// 方法1：通过size()获取长度
	if m, err := c.gdb.SendWithTimeout(OptionTimeout, "data-evaluate-expression", fmt.Sprintf("%s.size()", exp)); err == nil {
		payload := c.gdbOutputUtil.GetInterfaceFromMap(m, "payload")
		value := c.gdbOutputUtil.GetStringFromMap(payload, "value")
		if length, err := strconv.Atoi(value); err == nil && length > 0 {
			return length
		}
	}

	// 方法2：通过sizeof计算长度（兜底方案）
	if typ := c.extractVectorElementType(targetVariable.Type); typ != "" {
		if m, err := c.gdb.SendWithTimeout(OptionTimeout, "data-evaluate-expression", fmt.Sprintf("sizeof(%s)/sizeof(%s)", exp, typ)); err == nil {
			payload := c.gdbOutputUtil.GetInterfaceFromMap(m, "payload")
			value := c.gdbOutputUtil.GetStringFromMap(payload, "value")
			if length, err := strconv.Atoi(value); err == nil {
				return length
			}
		}
	}

	return 0
}

// extractVectorElementType 从std::vector类型中提取元素类型
// 匹配模式：std::vector<T, ...>
func (c *CPPDebugger) extractVectorElementType(typeStr string) string {
	re := regexp.MustCompile(`\bstd::\w+<\s*([^,\s>]+)(?:\s*,|\s*>)`)
	if match := re.FindStringSubmatch(typeStr); match != nil {
		return strings.TrimSpace(match[1])
	}
	return ""
}

// checkIsCppArrayType 检查是否为C++数组类型
// 包括std::array、C数组、std::vector
func (c *CPPDebugger) checkIsCppArrayType(targetVariable *dap.Variable) bool {
	return c.extractStdArrayLength(targetVariable.Type) > 0 ||
		c.extractCArrayLength(targetVariable.Type) > 0 ||
		strings.Contains(targetVariable.Type, "std::vector")
}

// parseObject2Keys 从结构体字符串表示中解析成员名称
// 匹配模式：member = value
func (c *CPPDebugger) parseObject2Keys(inputStr string) []string {
	re := regexp.MustCompile(`(\w+)\s*=`)
	matches := re.FindAllStringSubmatch(inputStr, -1)

	result := make([]string, 0, len(matches))
	for _, match := range matches {
		if key := match[1]; key != "\000" { // 过滤空字符
			result = append(result, key)
		}
	}
	return result
}

// 智能指针相关方法

// isSmartPointer 判断是否为智能指针类型
func (c *CPPDebugger) isSmartPointer(typeStr string) bool {
	return strings.Contains(typeStr, "std::unique_ptr<") ||
		strings.Contains(typeStr, "std::shared_ptr<") ||
		strings.Contains(typeStr, "std::weak_ptr<")
}

// getSmartPointerType 获取智能指针类型
func (c *CPPDebugger) getSmartPointerType(typeStr string) string {
	switch {
	case strings.Contains(typeStr, "std::unique_ptr<"):
		return "unique_ptr"
	case strings.Contains(typeStr, "std::shared_ptr<"):
		return "shared_ptr"
	case strings.Contains(typeStr, "std::weak_ptr<"):
		return "weak_ptr"
	default:
		return "unknown"
	}
}

// extractBaseType 从智能指针类型中提取基本类型
// 处理嵌套模板，如std::unique_ptr<std::vector<int>>
func (c *CPPDebugger) extractBaseType(typeStr string) string {
	smartPtrTypes := []string{"std::unique_ptr<", "std::shared_ptr<", "std::weak_ptr<"}

	for _, prefix := range smartPtrTypes {
		if strings.Contains(typeStr, prefix) {
			return c.extractTypeFromTemplate(typeStr, prefix)
		}
	}
	return typeStr
}

// extractTypeFromTemplate 从模板类型中提取基本类型
// 处理嵌套模板和默认参数
func (c *CPPDebugger) extractTypeFromTemplate(typeStr, prefix string) string {
	baseType := strings.TrimPrefix(typeStr, prefix)
	depth := 0

	for i, char := range baseType {
		switch char {
		case '<':
			depth++
		case '>':
			depth--
			if depth == 0 {
				part := baseType[:i+1]
				// 处理默认参数，如std::unique_ptr<T, Deleter>
				if commaIdx := strings.Index(part, ","); commaIdx != -1 {
					return strings.TrimSpace(part[:commaIdx])
				}
				return strings.TrimSpace(part)
			}
		}
	}
	return typeStr
}

// getSmartPointerExpression 获取智能指针的表达式
// 用于访问智能指针指向的对象
func (c *CPPDebugger) getSmartPointerExpression(varName, fieldPath, smartPtrType string) string {
	baseType := c.extractBaseType(varName)

	switch smartPtrType {
	case "weak_ptr":
		// weak_ptr需要先lock()再get()
		if fieldPath == "" {
			return fmt.Sprintf("*(%s *)(%s.lock().get())", baseType, varName)
		}
		return fmt.Sprintf("(*(%s *)(%s.lock().get()))->%s", baseType, varName, fieldPath)
	default: // unique_ptr, shared_ptr, unknown
		// unique_ptr和shared_ptr直接使用get()
		if fieldPath == "" {
			return fmt.Sprintf("*(%s *)(%s.get())", baseType, varName)
		}
		return fmt.Sprintf("(*(%s *)(%s.get()))->%s", baseType, varName, fieldPath)
	}
}

// GetExport 通过ReferenceStruct获取变量表达式
// 用于构建GDB命令中的变量访问表达式
func (c *CPPDebugger) GetExport(ref *ReferenceStruct) string {
	var exp string
	switch ref.Type {
	case "v": // 变量类型
		exp = ref.VariableName
	case "p": // 指针类型
		exp = fmt.Sprintf("*(%s)%s", ref.VariableType, ref.Address)
	}

	// 添加字段路径，但智能指针需要特殊处理
	if ref.FieldPath != "" && !c.isSmartPointer(ref.VariableType) {
		exp = fmt.Sprintf("(%s)%s", exp, ref.FieldPath)
	}
	return exp
}
