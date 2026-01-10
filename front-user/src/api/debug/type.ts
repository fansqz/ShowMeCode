import { BaseResponse } from '../type'

export type CreateDebugSessionResponse = BaseResponse & {
  // sessionid
  data: string
}

export type StartRequset = {
  id: string
  code: string
  language: string
  breakpoints: number[]
}

export type StartResponse = BaseResponse & {}

export type SendToConsoleResponse = BaseResponse & {}

export type StepInResponse = BaseResponse & {}

export type StepOutResponse = BaseResponse & {}

export type StepOverResponse = BaseResponse & {}

export type ContinueResponse = BaseResponse & {}

export type SetBreakpointResponse = BaseResponse & {
  data: Breakpoint[]
}

export type RemoveBreakpointResponse = BaseResponse & {}

export type TerminatedResponse = BaseResponse & {}

// 断点
export type Breakpoint = {
  verified: boolean
  message: string
  line: number
}
// 栈帧
export type StackFrame = {
  id: number // 栈帧id
  name: string // 函数名称
  path: string // 文件路径
  line: number
}

export type GetStackTraceResponse = BaseResponse & {
  data: StackFrame[]
}

// 变量
export type Variable = {
  name: string
  type: string
  value: string
  reference: number //变量引用
  namedVariables: number // 元素数量（结构体）
  indexedVariables: number // 元素数量（数据结构）
}

export type GetFrameVariablesResponse = BaseResponse & {
  data: Variable[]
}

export type GetVariablesResponse = BaseResponse & {
  data: Variable[]
}
