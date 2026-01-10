<template>
  <el-tooltip :content="'可视化调试'" placement="top" :show-after="800">
    <el-button
      tag="i"
      text
      :class="{
        'iconfont icon-caozuo-tiaoshi ing': isDebugging,
        'iconfont icon-caozuo-tiaoshi not-ing': debugStore.isDebugging() == false,
        'debug-button': true,
      }"
      :loading="loading"
      @click="startDebug"
    ></el-button>
  </el-tooltip>
</template>

<script setup lang="ts">
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
  import { ref } from 'vue'
  import { watch } from 'vue'
  import { onUnmounted } from 'vue'

  const codingStore = useCodingStore()
  const debugStore = useDebugStore()
  const { language, code } = storeToRefs(codingStore)
  // 标识的是否在调试中
  // 为什么不使用debugStore.isDebugging()，因为在执行（非调试）的方法中debugStore.isDebugging()也为true
  // 只有通过调试按钮以后isDebugging设置为true，才表示调试中
  let isDebugging = false
  let { status, id } = storeToRefs(debugStore)
  let loading = ref(false)

  const startDebug = async () => {
    // 调试状态，那么关闭调试
    if (debugStore.isDebugging()) {
      terminal()
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
        breakpoints: debugStore.breakpoints,
      }
      setTimeout(async () => {
        await reqStart(startReq)
      }, 1000)
    } else {
      loading.value = true
    }
  }

  const terminal = async () => {
    let result = await reqTerminate(id.value)
    if (result.code != 200) {
      status.value = 'terminated'
    }
    return
  }

  watch(
    () => status.value,
    (val) => {
      if (val == 'compiled' || val == 'terminated') {
        // 编译成功说明可以开始调试
        if (loading.value == true) {
          loading.value = false
          isDebugging = true
        }
      }
    },
  )

  onUnmounted(async () => {
    if (debugStore.isDebugging()) {
      await terminal()
      return
    }
  })
</script>

<style scoped lang="scss">
  .ing {
    color: $error-color;
    transition: color 0.2s;
  }

  .not-ing {
    color: grey;
    transition: color 0.2s;
  }

  .debug-button {
    width: 32px;
    height: 32px;
  }
</style>
