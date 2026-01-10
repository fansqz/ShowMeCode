<template>
  <div class="visual-panel">
    <header class="visual-panel__header">
      <el-text>可视化面板</el-text>
    </header>
    <main class="visual-panel__main" ref="mainContainer">
      <div class="visual-panel__canvas">
        <VisualCanvas ref="visualCanvas" :action="action" :sources="sources" />
      </div>
      <VisualSettings />
    </main>
  </div>
</template>

<script setup lang="ts">
  import useDebugStore from '@/store/modules/debug'
  import useVisualStore from '@/store/modules/visual'
  import { onMounted, ref } from 'vue'
  import { storeToRefs } from 'pinia'
  import VisualCanvas from '@/components/code-visual/visual/index.vue'
  import VisualSettings from '@/components/code-visual/visual-setting/index.vue'
  import { reqVisualData } from '@/components/code-visual/utils/index.ts'
  import { Sources } from 'structv2'
  import { watch } from 'vue'

  // Store instances
  const debugStore = useDebugStore()
  const visualStore = useVisualStore()

  // Store refs
  const { id, status } = storeToRefs(debugStore)
  const { action, descriptionType } = storeToRefs(visualStore)

  // Component refs
  const sources = ref<Sources>()
  const visualCanvas = ref<InstanceType<typeof VisualCanvas> | null>()
  const mainContainer = ref<HTMLElement>()

  /**
   * 获取可视化数据
   */
  const fetchVisualData = async () => {
    if (status.value !== 'stopped') {
      return
    }

    if (action.value) {
      const visualData = await reqVisualData(
        id.value,
        descriptionType.value,
        visualStore.getDescription(descriptionType.value),
      )
      sources.value = {
        visualData: visualData,
      }
    }
  }

  /**
   * 重新设置可视化画布的大小
   */
  const resizeVisualCanvas = () => {
    const width = mainContainer.value?.offsetWidth || 0
    const height = mainContainer.value?.offsetHeight || 0
    visualCanvas.value?.resizeVisualView(width, height)
  }

  // Lifecycle
  onMounted(() => {
    action.value = true

    // 监听调试状态变化
    watch(
      () => status.value,
      (newStatus) => {
        if (newStatus === 'stopped') {
          fetchVisualData()
        }
        if (newStatus === 'terminated') {
          sources.value = null
        }
      },
    )
  })

  // Expose methods
  defineExpose({
    resizeVisualCanvas,
  })
</script>

<style lang="scss" scoped>
  .visual-panel {
    display: flex;
    flex-direction: column;
    height: 100%;
    overflow: hidden;

    &__header {
      flex-shrink: 0;
      height: 35px;
      padding: 0 16px;
      line-height: 35px;
      user-select: none;
      background-color: $base-background-color;
      border-bottom: 1px solid $border-color;
    }

    &__main {
      position: relative;
      flex: 1;
      overflow: hidden;
      user-select: none;

      .visual-panel__canvas {
        position: absolute;
        inset: 0;
        background-color: $base-background-color;
        background-image: radial-gradient(circle, var(--project-dot-color) 1px, transparent 1px);
        background-size: 25px 25px;
      }
    }
  }
</style>
