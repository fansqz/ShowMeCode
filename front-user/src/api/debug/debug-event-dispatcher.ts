import {
  ConnectedEvent,
  StoppedEvent,
  LaunchEvent,
  ContinuedEvent,
  ExitedEvent,
  OutputEvent,
  CompileEvent,
  DebugEvent,
  TerminatedEvent,
} from '@/api/debug/event'
import { EventDispatcher } from '@/utils/event-dispatcher'

// 设置不同的事件分发器
export const ConnectedEventDispatcher = new EventDispatcher<ConnectedEvent>()

export const StoppedEventDispatcher = new EventDispatcher<StoppedEvent>()

export const LaunchEventDispatcher = new EventDispatcher<LaunchEvent>()

export const ContinuedEventDispatcher = new EventDispatcher<ContinuedEvent>()

export const ExitedEventDispatcher = new EventDispatcher<ExitedEvent>()

export const TerminatedEventDispatcher = new EventDispatcher<TerminatedEvent>()

export const OutputEventDispatcher = new EventDispatcher<OutputEvent>()

export const CompileEventDispatcher = new EventDispatcher<CompileEvent>()

export const DebugEventDispatcher = new EventDispatcher<DebugEvent>()
