<template>
  <div class="visual-canvas">
    <div class="visual-canvas__container" id="visualContainer" @wheel="handleWheelScroll"> </div>
    <div class="visual-canvas__controls">
      <button
        class="visual-canvas__zoom-btn"
        @mousedown="startContinuousZoom(1)"
        @mouseup="stopContinuousZoom"
        @mouseleave="stopContinuousZoom"
        title="放大"
      >
        <span class="zoom-icon">+</span>
      </button>
      <button
        class="visual-canvas__zoom-btn"
        @mousedown="startContinuousZoom(-1)"
        @mouseup="stopContinuousZoom"
        @mouseleave="stopContinuousZoom"
        title="缩小"
      >
        <span class="zoom-icon">-</span>
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
  import './layouter/link-list.ts'
  import './layouter/graph.ts'
  import './layouter/array.ts'
  import './layouter/array-bar.ts'
  import './layouter/binary-tree.ts'
  import './layouter/array2d.ts'
  import { SV, Sources, Engine } from 'structv2'
  import { onMounted, watch, toRefs, onUnmounted } from 'vue'

  // Props
  const props = defineProps<{
    sources: Sources
    action: boolean
  }>()
  const { sources, action } = toRefs(props)

  // Internal state
  let visualEngine: Engine | null = null
  let containerElement: HTMLElement | null = null
  let zoomTimer: number | null = null
  const ZOOM_INTERVAL = 50 // 缩放间隔时间（毫秒）

  /**
   * 更新可视化视图
   */
  const updateVisualView = (visualSources: Sources) => {
    if (!visualSources) {
      disposeVisualView()
      return
    }

    // 检查是否需要强制更新
    let isForceUpdate = false
    Object.keys(visualSources).forEach((key) => {
      if (visualSources[key].isForce === true) {
        isForceUpdate = true
      }
    })

    // 初始化可视化引擎
    if (!visualEngine) {
      visualEngine = SV(
        containerElement,
        {
          view: {
            groupPadding: 0,
            fitCenter: true,
            centerOffsetYPercent: -0.2,
          },
        },
        isForceUpdate,
      )
    }

    visualEngine.render(visualSources)
  }

  /**
   * 销毁可视化视图
   */
  const disposeVisualView = () => {
    if (visualEngine) {
      visualEngine.destroy()
      visualEngine = null
    }
  }

  /**
   * 调整可视化视图尺寸
   */
  const resizeVisualView = (width: number, height: number) => {
    if (visualEngine) {
      visualEngine.resize(width, height)
    }
  }

  /**
   * 处理鼠标滚轮事件
   */
  const handleWheelScroll = (event: WheelEvent) => {
    // Ctrl + 滚轮进行缩放
    if (event.ctrlKey) {
      event.preventDefault()
      const zoomDirection = event.deltaY > 0 ? 1 : -1
      performZoom(zoomDirection)
    }
  }

  /**
   * 执行缩放操作
   */
  const performZoom = (increment: number) => {
    if (!visualEngine) return

    const currentZoom = visualEngine.getGraphInstance().getZoom()
    const newZoom = currentZoom + increment / 100
    visualEngine.getGraphInstance().zoomTo(newZoom)
  }

  /**
   * 开始持续缩放
   */
  const startContinuousZoom = (increment: number) => {
    // 立即执行一次缩放
    performZoom(increment)

    // 设置定时器持续缩放
    zoomTimer = window.setInterval(() => {
      performZoom(increment)
    }, ZOOM_INTERVAL)
  }

  /**
   * 停止持续缩放
   */
  const stopContinuousZoom = () => {
    if (zoomTimer !== null) {
      clearInterval(zoomTimer)
      zoomTimer = null
    }
  }

  // Lifecycle
  onMounted(() => {
    containerElement = document.getElementById('visualContainer')

    // 监听数据源变化
    watch(
      () => sources.value,
      (newSources) => {
        updateVisualView(newSources)
      },
    )

    // 监听可视化开关状态
    watch(
      () => action.value,
      (isActive) => {
        if (isActive) {
          updateVisualView(sources.value)
        } else {
          updateVisualView(null)
        }
      },
    )
  })

  onUnmounted(() => {
    disposeVisualView()
    stopContinuousZoom()
  })

  // Expose methods
  defineExpose({
    resizeVisualView,
  })
</script>

<style lang="scss" scoped>
  .visual-canvas {
    position: relative;
    width: 100%;
    height: 100%;
    max-height: 100% !important;
    padding: 0;
    margin: 0;

    &__container {
      position: absolute;
      top: 0;
      left: 0;
      width: 100%;
      height: 100%;
    }

    &__controls {
      position: absolute;
      top: 50%;
      right: 20px;
      z-index: 1000;
      display: flex;
      flex-direction: column;
      gap: 8px;
      transform: translateY(-50%);
    }

    &__zoom-btn {
      display: flex;
      align-items: center;
      justify-content: center;
      width: 32px;
      height: 32px;
      font-size: 18px;
      font-weight: 500;
      color: $text-color-primary;
      cursor: pointer;
      user-select: none;
      background: rgba($base-background-color, 0.9);
      backdrop-filter: blur(8px);
      border: 1px solid $border-color-light;
      border-radius: $border-radius-medium;
      box-shadow: $box-shadow-light;
      transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);

      &:hover {
        background: rgba($base-background-color, 1);
        border-color: $border-color;
        box-shadow: $box-shadow;
        transform: translateY(-1px);
      }

      &:active {
        box-shadow: $box-shadow-light;
        transform: translateY(1px);
      }

      &::before {
        position: absolute;
        inset: 0;
        padding: 1px;
        content: '';
        background: linear-gradient(
          135deg,
          rgba($base-background-color, 0.4),
          rgba($base-background-color, 0.1)
        );
        border-radius: $border-radius-medium;
        mask:
          linear-gradient(var(--el-color-white, #fff) 0 0) content-box,
          linear-gradient(var(--el-color-white, #fff) 0 0);
        mask-composite: xor;
        mask-composite: exclude;
      }

      .zoom-icon {
        position: relative;
        z-index: 1;
      }
    }
  }
</style>
