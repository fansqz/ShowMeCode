package constants

type DebugMessageType string

const (
	RequestMessage  DebugMessageType = "request"
	ResponseMessage DebugMessageType = "response"
	EventMessage    DebugMessageType = "event"
)

// DebugOptionType 调试请求操作类型
type DebugOptionType string

const (
	// StartDebug 开始调试过程，返回可能出现的错误。
	StartDebug DebugOptionType = "start"
	// SendToConsole 输入数据到控制台，返回可能出现的错误。
	SendToConsole DebugOptionType = "sendToConsole"
	// Next 执行下一步操作，但不会进入函数内部，返回可能出现的错误。
	Next DebugOptionType = "next"
	// Step 执行下一步操作，会进入函数内部（如有调用函数，则步入函数），返回可能出现的错误。
	Step DebugOptionType = "step"
	// Continue：继续执行程序，直到遇到下一个断点或程序结束，返回可能出现的错误。
	Continue DebugOptionType = "continue"
	// AddBreakpoints 添加断点，接受文件源和断点列表，返回添加成功的断点和可能出现的错误。
	AddBreakpoints DebugOptionType = "addBreakpoints"
	// RemoveBreakpoints 移除断点，接受文件源和断点列表，返回移除成功的断点和可能出现的错误。
	RemoveBreakpoints DebugOptionType = "removeBreakpoints"
	// Terminate 终止当前的调试会话，之后可以重新调用 Launch 方法开始新的会话，返回可能出现的错误。
	Terminate DebugOptionType = "terminate"
)

type DebugEventType string

const (
	BreakpointEvent DebugEventType = "breakpoint"
	OutputEvent     DebugEventType = "output"
	StoppedEvent    DebugEventType = "stopped"
	ContinuedEvent  DebugEventType = "continued"
	CompileEvent    DebugEventType = "compile"
	ExitedEvent     DebugEventType = "exited"
	TerminatedEvent DebugEventType = "terminated"
	LaunchEvent     DebugEventType = "launch"
)

// BreakpointReasonType 断点改变类型
type BreakpointReasonType string

const (
	ChangeType  BreakpointReasonType = "change"
	NewType     BreakpointReasonType = "new"
	RemovedType BreakpointReasonType = "removed"
)

// StoppedReasonType 程序停止类型
type StoppedReasonType string

const (
	BreakpointStopped StoppedReasonType = "breakpoint"
	StepStopped       StoppedReasonType = "step"
	ExitedNormally    StoppedReasonType = "exited-normally"
)

// StepType 单步调试类型
type StepType string

const (
	StepIn   StepType = "stepIn"
	StepOut  StepType = "stepOut"
	StepOver StepType = "stepOver"
)

// ScopeName 作用域名称
type ScopeName string

// Local: 函数或当前代码块内的局部变量。这是最常访问的作用域，包含了当前栈帧中的局部变量和参数。
// Global: 整个程序的全局变量。这包括在函数、类或文件范围之外声明的所有变量。
// Static: 存在于静态存储区域的静态变量，生命周期贯穿程序执行的整个过程。它通常包括静态成员变量。
// Class: 当前类级别的作用域，包含了类的成员变量。
// Object: 针对特定对象的作用域，包含了对象的属性。
// Closure: 函数闭包的作用域。在一些编程语言中，如果一个内部函数访问了其外部函数的变量，则一个闭包会被创建。
// Module: 模块级别的作用域，通常指向程序中一个模块或命名空间内的变量。
// Register (Debugger Internal): 寄存器级别的作用域，通常用于底层编程或汇编层面的调试，引用硬件寄存器的内容。
const (
	ScopeLocal    ScopeName = "local"
	ScopeGlobal   ScopeName = "common"
	ScopeStatic   ScopeName = "static"
	ScopeClass    ScopeName = "class"
	ScopeObject   ScopeName = "object"
	ScopeClosure  ScopeName = "closure"
	ScopeModule   ScopeName = "module"
	ScopeRegister ScopeName = "register"
)
