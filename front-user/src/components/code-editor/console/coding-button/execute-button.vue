<template>
  <el-button
    class="button-submit"
    type="primary"
    :loading="loading"
    size="small"
    @click="startDebug"
    >执行</el-button
  >
</template>

<script setup lang="ts">
  // 执行暂时和调试一致，只是不设置断点
  import useCodingStore from '@/store/modules/coding'
  import {
    reqCreateDebugSession,
    reqStart,
    reqListenDebugEvent,
    reqTerminate,
  } from '@/api/debug/index.ts'
  import { storeToRefs } from 'pinia'
  import useDebugStore from '@/store/modules/debug'
  import { listenDebugEvent } from '@/api/debug/debug-event-listen.ts'
  import { ref, watch } from 'vue'

  const codingStore = useCodingStore()
  const debugStore = useDebugStore()
  const { language, code } = storeToRefs(codingStore)
  let { status, id } = storeToRefs(debugStore)
  let loading = ref(false)

  const startDebug = async () => {
    // 调试状态，那么关闭调试
    if (debugStore.isDebugging()) {
      let result = await reqTerminate(id.value)
      if (result.code != 200) {
        status.value = 'terminated'
      }
      return
    }
    loading.value = true
    // 非调试状态，启动调试
    let result = await reqCreateDebugSession()
    if (result.code == 200) {
      // 开启loading
      status.value = 'init'
      loading.value = true
      id.value = result.data
      // 启动监控管道
      let eventSource = reqListenDebugEvent(id.value)
      listenDebugEvent(id.value, eventSource)
      // 发送启动调试命令
      let startReq = {
        id: id.value,
        code: code.value,
        language: language.value,
        breakpoints: [],
      }
      setTimeout(async () => {
        await reqStart(startReq)
      }, 1000)
    } else {
      loading.value = true
    }
  }

  watch(
    () => status.value,
    (val) => {
      if (val == 'compiled' || val == 'terminated') {
        // 编译成功说明可以开始调试
        loading.value = false
      }
    },
  )
</script>
