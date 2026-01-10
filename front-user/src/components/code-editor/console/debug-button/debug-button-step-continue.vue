<template>
  <div
    :class="{ button: true, 'button-show': buttonShow }"
    @mouseover="mouseover"
    @mouseout="mouseout"
  >
    <el-tooltip :content="tooltipContent" placement="top" :show-after="800" :disabled="!able">
      <el-button
        :class="{
          iconfont: true,
          'icon-continue1': type == 'continue',
          'icon-step-in': type == 'step-in',
          'icon-step-out': type == 'step-out',
          'icon-step': type == 'step-over',
          able: able && type != 'continue',
          cable: able && type == 'continue',
          unable: !able,
        }"
        @click="handleStep"
        link
      ></el-button>
    </el-tooltip>
  </div>
</template>

<script setup lang="ts">
  import { reqContinue, reqStepIn, reqStepOut, reqStepOver } from '@/api/debug/index.ts'
  import useDebugStore from '@/store/modules/debug'
  import { ref, toRefs, watch, computed } from 'vue'
  import { storeToRefs } from 'pinia'

  let props = defineProps<{
    type: string
  }>()
  let { type } = toRefs(props)
  const debugStore = useDebugStore()
  let { status } = storeToRefs(debugStore)

  // 按钮提示文本
  const tooltipContent = computed(() => {
    switch (type.value) {
      case 'continue':
        return '继续执行'
      case 'step-in':
        return '单步进入'
      case 'step-out':
        return '单步退出'
      case 'step-over':
        return '单步跳过'
      default:
        return ''
    }
  })

  // 添加点击间隔的状态和时间常量
  const clickInterval = 100
  let lastClicked = 0
  // 判断按钮是出于可执行还是不可执行状态
  let able = ref(false)
  watch(
    () => status.value,
    () => {
      able.value = status.value == 'stopped'
    },
  )
  const handleStep = async () => {
    // 检查距离上次点击的时间是否小于设定的间隔
    const now = Date.now()
    if (!able.value || now - lastClicked < clickInterval) {
      return
    }
    lastClicked = now
    if (type.value == 'continue') {
      reqContinue(debugStore.id)
    } else if (type.value == 'step-in') {
      reqStepIn(debugStore.id)
    } else if (type.value == 'step-out') {
      reqStepOut(debugStore.id)
    } else if (type.value == 'step-over') {
      reqStepOver(debugStore.id)
    }
  }

  let buttonShow = ref(false)
  const mouseover = () => {
    // 鼠标经过显示按钮
    buttonShow.value = able.value
  }
  const mouseout = () => {
    // 鼠标移出，隐藏按钮
    buttonShow.value = false
  }
</script>

<style scoped lang="scss">
  .button {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 20px;
    height: 20px;
    background-color: $base-background-color;
    border-radius: $border-radius-small;

    .cable {
      color: $success-color;
      transition: color 0.2s;
    }

    .able {
      color: $success-color;
      transition: color 0.2s;
    }

    .unable {
      color: $info-color;
      transition: color 0.2s;
    }
  }

  .button-show {
    background-color: $border-color-light;
  }
</style>
