import { storeToRefs } from 'pinia'

import useDebugStore from '@/store/modules/debug'

import {
  CompileEventDispatcher,
  ConnectedEventDispatcher,
  ContinuedEventDispatcher,
  DebugEventDispatcher,
  OutputEventDispatcher,
  StoppedEventDispatcher,
  TerminatedEventDispatcher,
} from './debug-event-dispatcher'

const debugStore = useDebugStore()

const { id } = storeToRefs(debugStore)

export const listenDebugEvent = (debugId: string, eventSource: EventSource) => {
  eventSource.onmessage = function (event) {
    if (debugId != id.value) {
      return
    }
    if (event.data == 'connect success') {
      ConnectedEventDispatcher.dispatch('connected', event)
      DebugEventDispatcher.dispatch('connected', event)
      return
    }
    const data = JSON.parse(event.data)
    if (data.event == 'stopped') {
      StoppedEventDispatcher.dispatch('stopped', data)
      DebugEventDispatcher.dispatch('stopped', data)
    }
    if (data.event == 'continued') {
      ContinuedEventDispatcher.dispatch('continued', data)
      DebugEventDispatcher.dispatch('continued', data)
    }
    if (data.event == 'output') {
      OutputEventDispatcher.dispatch('output', data)
      DebugEventDispatcher.dispatch('output', data)
    }
    if (data.event == 'compile') {
      CompileEventDispatcher.dispatch('compile', data)
      DebugEventDispatcher.dispatch('compile', data)
    }
    if (data.event == 'terminated') {
      CompileEventDispatcher.dispatch('terminated', data)
      TerminatedEventDispatcher.dispatch('terminated', data)
      // 关闭sse
      eventSource.close()
    }
  }
}
