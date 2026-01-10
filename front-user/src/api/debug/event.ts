export type DebugEventType = string

export const BreakpointEvent: DebugEventType = 'breakpoint'
export const OutputEvent: DebugEventType = 'output'
export const StoppedEvent: DebugEventType = 'stopped'
export const ContinuedEvent: DebugEventType = 'continued'
export const CompileEvent: DebugEventType = 'compile'
export const ExitedEvent: DebugEventType = 'exited'
export const LaunchEvent: DebugEventType = 'launch'

export type BreakpointReasonType = string

export const ChangeType: BreakpointReasonType = 'change'
export const NewType: BreakpointReasonType = 'new'
export const RemovedType: BreakpointReasonType = 'removed'

export type StoppedReasonType = string

export const BreakpointStopped: StoppedReasonType = 'breakpoint'
export const StepStopped: StoppedReasonType = 'step'

export type ConnectedEvent = {}

// 事件类型，通过sse连接异步的方式返回给前端
// CompileEvent 编译事件
export type CompileEvent = {
  event: DebugEventType
  success: boolean
  message: string // 编译产生的信息
}

// LaunchEvent 加载用户代码事件
export type LaunchEvent = {
  event: DebugEventType
  success: boolean
  message: string // 启动 gdb 的消息
}

// BreakpointEvent 断点事件
// 该 event 指示有关断点的某些信息已更改。
export type BreakpointEvent = {
  event: DebugEventType
  reason: BreakpointReasonType
  breakpoints: number[]
}

// OutputEvent
// 该事件表明目标已经产生了一些输出。
export type OutputEvent = {
  event: DebugEventType
  output: string // 输出内容
}

// StoppedEvent
// 该 event 表明，由于某些原因，被调试进程的执行已经停止。
export type StoppedEvent = {
  event: DebugEventType
  reason: StoppedReasonType
  line: number // 停止在某行
}

// ContinuedEvent
// 该 event 表明 debug 的执行已经继续。
// 请注意:debug adapter 不期望发送此事件来响应暗示执行继续的请求，例如启动或继续。
// 它只有在没有先前的 request 暗示这一点时，才有必要发送一个持续的事件。
export type ContinuedEvent = {
  event: DebugEventType
}

// ExitedEvent
// 该 event 表明被调试对象已经退出并返回 exit code。
export type ExitedEvent = {
  event: DebugEventType
  exitCode: number
  message: string
}

// TerminatedEvent
// 程序退出事件
export type TerminatedEvent = {
  event: DebugEventType
}

export type DebugEvent =
  | ConnectedEvent
  | StoppedEvent
  | LaunchEvent
  | ContinuedEvent
  | ExitedEvent
  | TerminatedEvent
  | OutputEvent
  | CompileEvent
