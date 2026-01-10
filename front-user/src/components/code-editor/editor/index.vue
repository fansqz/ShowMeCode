<template>
  <div ref="editor" class="editor" />
</template>

<script setup lang="ts">
  import { ref } from 'vue'
  import { useVsCode } from './hooks/useVSCode'
  import useDebugStore from '@/store/modules/debug'
  import { storeToRefs } from 'pinia'
  import { reqSetBreakpoint } from '@/api/debug'
  import useCodingStore from '@/store/modules/coding'
  import {
    OutputEventDispatcher,
    StoppedEventDispatcher,
    ContinuedEventDispatcher,
    CompileEventDispatcher,
    TerminatedEventDispatcher,
    ConnectedEventDispatcher,
  } from '@/api/debug/debug-event-dispatcher'
  import { OutputEvent, CompileEvent } from '@/api/debug/event'
  import { onMounted } from 'vue'
  import { onUnmounted } from 'vue'

  const debugStore = useDebugStore()
  const codingStore = useCodingStore()
  let { id, outputs, breakpoints, status } = storeToRefs(debugStore)
  let { code } = storeToRefs(codingStore)

  const emits = defineEmits<{
    // 修改文本事件
    (event: 'onChangeValue', value: string, type: 'input' | 'blur'): void
    // 设置断点事件
    (event: 'onSetBP', BP: string[] | number[]): void
    // 更新断点事件
    (event: 'onUpdateBP', BP: number, flag: 'add' | 'del'): void
  }>()
  const onContentChanged = (value: string) => {
    code.value = value
    emits('onChangeValue', value, 'input')
  }

  const onEditorBlur = (value: string) => {
    emits('onChangeValue', value, 'blur')
  }

  const onCtrlS = (value: string) => {
    emits('onChangeValue', value, 'blur')
  }
  const onSetBP = (bps?: number[]) => {
    breakpoints.value = bps
    if (debugStore.isDebugging()) {
      // 断点更新
      emits('onSetBP', bps)
      reqSetBreakpoint(id.value, bps)
    }
  }

  const onUpdateBP = (bps: number[], lineNum: number, mode: 'add' | 'del') => {
    onSetBP(bps)
    emits('onUpdateBP', lineNum, mode)
  }

  // 处理不同调试事件，更新状态
  const onStopped = () => {
    status.value = 'stopped'
  }
  const onContinued = () => {
    status.value = 'running'
  }
  const onCompiled = (data: CompileEvent) => {
    if (!data.success) {
      // 调试失败，设置状态为终止状态
      status.value = 'terminated'
      outputs.value.push({
        type: 'error',
        event: 'compile',
        title: '编译错误',
        message: data.message,
      })
    } else {
      status.value = 'compiled'
      outputs.value.push({
        type: 'success',
        event: 'compile',
        title: '编译成功',
        message: data.message,
      })
    }
  }
  const onOutput = (data: OutputEvent) => {
    outputs.value.push({
      type: 'info',
      event: 'output',
      title: '',
      message: data.output,
    })
  }
  const onTerminated = () => {
    status.value = 'terminated'
    debugStore.lineNum = 0
    outputs.value.push({
      type: 'success',
      event: 'terminated',
      message: '调试结束',
    })
  }
  // 启动一个新的调试管道时候清除原来的input和output
  const onConnected = () => {
    status.value = 'init'
    outputs.value = []
  }

  onMounted(() => {
    // 注册一些事件
    StoppedEventDispatcher.on('stopped', onStopped)
    ContinuedEventDispatcher.on('continued', onContinued)
    OutputEventDispatcher.on('output', onOutput)
    CompileEventDispatcher.on('compile', onCompiled)
    TerminatedEventDispatcher.on('terminated', onTerminated)
    ConnectedEventDispatcher.on('connected', onConnected)
  })

  onUnmounted(() => {
    StoppedEventDispatcher.off('stopped', onStopped)
    ContinuedEventDispatcher.off('continued', onContinued)
    OutputEventDispatcher.off('output', onOutput)
    CompileEventDispatcher.off('compile', onCompiled)
    TerminatedEventDispatcher.off('terminated', onTerminated)
    ConnectedEventDispatcher.off('connected', onConnected)
  })

  const editor = ref<HTMLElement>()
  useVsCode({
    target: editor,
    onContentChanged,
    onEditorBlur,
    onCtrlS,
    onSetBP,
    onUpdateBP,
  })
</script>

<style scoped lang="scss">
  @import url('./style/breakpoint.scss');
  @import url('./style/debug.scss');

  .editor {
    position: absolute;
    width: 100%;
    height: 100%;
    max-height: 100% !important;
    padding: 0;
    margin: 0;
    overflow: hidden;
  }
</style>
