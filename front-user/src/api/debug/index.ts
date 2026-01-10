import request from '@/utils/request'

import {
  CreateDebugSessionResponse,
  StartRequset,
  StartResponse,
  SendToConsoleResponse,
  StepInResponse,
  StepOutResponse,
  StepOverResponse,
  ContinueResponse,
  SetBreakpointResponse,
  TerminatedResponse,
  GetStackTraceResponse,
  GetFrameVariablesResponse,
  GetVariablesResponse,
} from './type'

enum API {
  CreateSessionURL = '/debug/session/create',
  StartURL = '/debug/start',
  SseURL = '/debug/sse',
  NextURL = '/debug/next',
  StepInURL = '/debug/step/in',
  StepOutURL = '/debug/step/out',
  StepOverURL = '/debug/step/over',
  ContinueURL = '/debug/continue',
  SendToConsoleURL = '/debug/sendToConsole',
  SetBreakpointsURL = '/debug/setBreakpoints',
  GetStackTraceURL = '/debug/stackTrace',
  GetFrameVariablesURL = '/debug/frame/variables',
  GetVariablesURL = '/debug/variables',
  TerminateURL = '/debug/terminate',
}

// reqCreateDebugSession 创建一个调试session
export const reqCreateDebugSession = (): Promise<CreateDebugSessionResponse> => {
  return request.post(API.CreateSessionURL)
}

// reqListenDebugEvent 创建sse连接，返回一个eventsource
export const reqListenDebugEvent = (id: string): EventSource => {
  const baseApi = import.meta.env.VITE_APP_BASE_API
  const serve = import.meta.env.VITE_SERVE
  const debugURL = new URL(`${baseApi}${API.SseURL}/${id}`, serve)
  const source = new EventSource(debugURL)
  return source
}

// reqStart 开始调试
export const reqStart = (req: StartRequset): Promise<StartResponse> => {
  return request.post(API.StartURL, req)
}

// reqSendToConsole 发送用户输入到目标程序
export const reqSendToConsole = (id: string, input: string): Promise<SendToConsoleResponse> => {
  return request.post(API.SendToConsoleURL, {
    input: input,
    id: id,
  })
}

// reqStepIn 单步调试，会进入函数内部
export const reqStepIn = (id: string): Promise<StepInResponse> => {
  return request.post(API.StepInURL, { id: id })
}

// reqStepOut 单步调试，会跳出函数外部
export const reqStepOut = (id: string): Promise<StepOutResponse> => {
  return request.post(API.StepOutURL, { id: id })
}

// reqStepOver 单步调试，不会进入函数内部
export const reqStepOver = (id: string): Promise<StepOverResponse> => {
  return request.post(API.StepOverURL, { id: id })
}

// reqContinue continue程序继续执行直到遇到下一个断点
export const reqContinue = (id: string): Promise<ContinueResponse> => {
  return request.post(API.ContinueURL, { id: id })
}

// reqSetBreakpoint 设置断点
export const reqSetBreakpoint = (
  id: string,
  breakpoints: number[],
): Promise<SetBreakpointResponse> => {
  return request.post(API.SetBreakpointsURL, {
    id: id,
    breakpoints: breakpoints,
  })
}

// reqTerminate 关闭调试session
export const reqTerminate = (id: string): Promise<TerminatedResponse> => {
  return request.post(API.TerminateURL, { id: id })
}

// reqGetStackTrace 根据debugid获取当前线程栈帧
export const reqGetStackTrace = (id: string): Promise<GetStackTraceResponse> => {
  return request.post(API.GetStackTraceURL, { id: id })
}

// reqGetFrameVariables 根据栈帧id获取某个栈帧的所有局部变量
export const reqGetFrameVariables = (
  id: string,
  frameId: number,
): Promise<GetFrameVariablesResponse> => {
  return request.post(API.GetFrameVariablesURL, {
    id: id,
    frameID: frameId,
  })
}

// reqGetVariables 根据引用获取某个变量的内容
export const reqGetVariables = (id: string, reference: number): Promise<GetVariablesResponse> => {
  return request.post(API.GetVariablesURL, {
    id: id,
    reference: reference,
  })
}
