<template>
  <div class="container">
    <div class="visual-container" ref="visualContainer">
      <Visaul ref="visual" class="visaul" :action="action" :sources="sources" />
    </div>
    <VisualSettings />
  </div>
</template>

<script setup lang="ts">
  import useDebugStore from '@/store/modules/debug'
  import useVisualStore from '@/store/modules/visual'
  import { onMounted, ref, watch } from 'vue'
  import { storeToRefs } from 'pinia'
  import Visaul from '@/components/code-visual/visual/index.vue'
  import VisualSettings from '@/components/code-visual/visual-setting/index.vue'
  import { reqVisualData } from '@/components/code-visual/utils/index.ts'
  import { Sources } from 'structv2'

  const debugStore = useDebugStore()
  const visualStore = useVisualStore()
  const { id, status } = storeToRefs(debugStore)
  const { action, descriptionType } = storeToRefs(visualStore)
  const sources = ref<Sources>()
  const visual = ref<InstanceType<typeof Visaul> | null>()
  const visualContainer = ref<HTMLElement>()

  const getVisualData = async () => {
    if (status.value != 'stopped') {
      return
    }
    let visualData = await reqVisualData(
      id.value,
      descriptionType.value,
      visualStore.getDescription(descriptionType.value),
    )
    // 设置可视化数据
    sources.value = {
      visaulData: visualData,
    }
  }

  // 重新设置visual的大小
  const resizeVisualView = () => {
    visual.value?.resizeVisualView(
      visualContainer.value?.offsetWidth,
      visualContainer.value?.offsetHeight,
    )
  }

  defineExpose({
    resizeVisualView,
  })

  onMounted(() => {
    // 注册一些事件
    watch(
      () => status.value,
      (val) => {
        if (val == 'stopped') {
          getVisualData()
        }
        if (val == 'terminated') {
          sources.value = null
        }
      },
    )
  })
</script>

<style lang="scss" scoped>
  .container {
    position: relative;
    width: 100%;
    height: 100%;
    margin: 0;

    .visual-container {
      position: absolute;
      flex: 1;
      width: 100%;
      height: 100%;
      background-color: $base-background-color;
      background-image: radial-gradient(circle, var(--project-dot-color) 1px, transparent 1px);
      background-size: 25px 25px;
    }
  }
</style>
