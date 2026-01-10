package debug_core

import "github.com/fansqz/fancode-backend/constants"

// StructVisualQuery 可视化查询的参数
// 结构体导向，会根据结构体去遍历这个结构体的所有数据和指针
// 1. 该方法会找出所有指向该结构体的指针
// 2. 深度遍历所有指定type的结构体数据并返回
type StructVisualQuery struct {
	// 需要查询的结构体名称
	Struct string `json:"struct"`
	// Values 数据域
	Values []string `json:"values"`
	// Points 指针域，指针域如果是一个数组，那么这个数组下所有元素都将作为指针
	Points []string `json:"points"`
}

// StructVisualData 可视化拆查询返回的数据
type StructVisualData struct {
	// Nodes 可视化结构的节点列表
	Nodes []*StructVisualNode `json:"nodes"`
	// Points 变量列表
	Points []*VisualVariable `json:"points"`
}

// StructVisualNode 结构体可视化的一个节点
// 包含所有的数据域和指针域
type StructVisualNode struct {
	Name string `json:"name"`
	// 可以理解为地址
	ID     string            `json:"id"`
	Type   string            `json:"type"`
	Values []*VisualVariable `json:"values"`
	Points []*VisualVariable `json:"points"`
}

type ArrayVisualQuery struct {
	// StructVars 作为结构体的变量
	ArrayName string `json:"arrayName"`
	// PointVars 作为指针的变量
	PointNames []string `json:"pointNames"`
}

type ArrayVisualData struct {
	Array  []*VisualVariable `json:"array"`
	Points []*VisualVariable `json:"points"`
}

type Array2DVisualQuery struct {
	// StructVars 作为结构体的变量
	ArrayName string `json:"arrayName"`
	// 二维数组的指针，行和列
	RowPointNames []string `json:"rowPointNames"`
	ColPointNames []string `json:"colPointNames"`
}

type Array2DVisualData struct {
	Array     [][]*VisualVariable `json:"array"`
	RowPoints []*VisualVariable   `json:"rowPoints"`
	ColPoints []*VisualVariable   `json:"colPoints"`
}

type VisualVariable struct {
	// 变量名称
	Name string `json:"name"`
	// 变量类型
	Type string `json:"type"`
	// 变量的值，在可视化中大部分作为指针使用
	Value string `json:"value"`
}

func NewVisualVariable(variable *Variable) *VisualVariable {
	return &VisualVariable{
		Name:  variable.Name,
		Type:  variable.Type,
		Value: variable.Value,
	}
}

// Breakpoint 表示断点
type Breakpoint struct {
	// Verified 断点是否设置成功
	Verified bool   `json:"verified"`
	Message  string `json:"message,omitempty"`
	Line     int    `json:"line,omitempty"`
}

// StackFrame 栈帧
type StackFrame struct {
	ID   int    `json:"id"`   // 栈帧id
	Name string `json:"name"` // 函数名称
	Path string `json:"path"` // 文件路径
	Line int    `json:"line"`
}

// Scope 作用域
type Scope struct {
	Name      constants.ScopeName
	Reference string // 作用域的引用
}

// Variable 变量
type Variable struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value string `json:"value"`
	// 变量引用
	Reference        int `json:"reference"`
	NamedVariables   int `json:"namedVariables"`
	IndexedVariables int `json:"indexedVariables"`
}

// OutputEvent
// 用户程序输出
type OutputEvent struct {
	Output string // 输出内容
}

func NewOutputEvent(output string) *OutputEvent {
	return &OutputEvent{
		Output: output,
	}
}

// StoppedEvent
// 该event表明，由于某些原因，被调试进程的执行已经停止。
// 这可能是由先前设置的断点、完成的步进请求、执行调试器语句等引起的。
type StoppedEvent struct {
	Reason constants.StoppedReasonType // 停止执行的原因
}

func NewStoppedEvent(reason constants.StoppedReasonType) *StoppedEvent {
	return &StoppedEvent{
		Reason: reason,
	}
}

// ContinuedEvent
// 该event表明debug的执行已经继续。
// 请注意:debug adapter不期望发送此事件来响应暗示执行继续的请求，例如启动或继续。
// 它只有在没有先前的request暗示这一点时，才有必要发送一个持续的事件。
type ContinuedEvent struct {
}

func NewContinuedEvent() *ContinuedEvent {
	return &ContinuedEvent{}
}

// ExitedEvent
// 该event表明被调试对象已经退出并返回exit code。但是并不意味着调试会话结束
type ExitedEvent struct {
	ExitCode int
	Message  string
}

func NewExitedEvent(code int, message string) *ExitedEvent {
	return &ExitedEvent{
		ExitCode: code,
		Message:  message,
	}
}

type TerminalEvent struct {
}

func NewTerminalEvent() *TerminalEvent {
	return &TerminalEvent{}
}

// CompileEvent
// 编译事件
type CompileEvent struct {
	Success bool
	Message string // 编译产生的信息
}

func NewCompileEvent(success bool, message string) *CompileEvent {
	return &CompileEvent{
		Success: success,
		Message: message,
	}
}
