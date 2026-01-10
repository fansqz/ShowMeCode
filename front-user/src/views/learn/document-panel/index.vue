<template>
  <div class="document">
    <el-menu
      :default-active="activeIndex"
      class="menu"
      mode="horizontal"
      :ellipsis="false"
      @select="handleSelect"
    >
      <el-menu-item index="0">文档</el-menu-item>
      <el-menu-item index="1">可视化</el-menu-item>
    </el-menu>
    <div class="main" ref="main">
      <ProblemDescription :content="content" v-if="activeIndex == '0'" />
      <StructVisual ref="structVisual" v-if="activeIndex == '1'" />
    </div>
  </div>
</template>

<script setup lang="ts">
  import { ref, computed } from 'vue'
  import ProblemDescription from '@/components/code-visual/document/index.vue'
  import StructVisual from './visual.vue'
  import useVisualDocumentStore from '@/store/modules/visual-document'
  import { storeToRefs } from 'pinia'

  const visualDocumentStore = useVisualDocumentStore()
  const { content } = storeToRefs(visualDocumentStore)

  const activeIndex = ref('0')
  const main = ref<HTMLElement>()
  const structVisual = ref()

  const handleSelect = (key: string) => {
    activeIndex.value = key
  }

  // 重新设置visual的大小
  const resizeVisualView = () => {
    structVisual.value?.resizeVisualView(main.value?.offsetWidth, main.value?.offsetHeight)
  }

  defineExpose({
    resizeVisualView,
  })
</script>

<style scoped lang="scss">
  .document {
    width: 100%;
    height: 100%;

    .menu {
      height: 35px;
      user-select: none;
      background-color: $base-header-background;
      border-bottom: 1px solid $border-color;
    }

    .main {
      height: calc(100% - 35px);
    }
  }
</style>
