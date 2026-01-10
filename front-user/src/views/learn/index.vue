<template>
  <div class="container">
    <el-splitter class="default-theme main">
      <!--左侧面板：目录树和题目展示-->
      <el-splitter-panel
        @resized="resizeVisualPane"
        class="left-panel"
        size="60"
        @update:size="resizeVisualPane"
      >
        <div class="left-content">
          <!--目录树-->
          <div class="directory-section">
            <DirectoryTree />
          </div>

          <!--题目展示-->
          <div class="document-section">
            <DocumentPanel ref="documentPanel" />
          </div>
        </div>
      </el-splitter-panel>

      <!--右侧面板：coding区域-->
      <el-splitter-panel class="coding-panel" size="40">
        <CodingPanel />
      </el-splitter-panel>
    </el-splitter>
  </div>
</template>

<script setup lang="ts">
  import { onMounted, onUnmounted, ref, watch } from 'vue'
  import DirectoryTree from './directory-tree.vue'
  import DocumentPanel from './document-panel/index.vue'
  import CodingPanel from './coding-panel/index.vue'
  import { reqVisualDocument } from '@/api/visual-document/index.ts'
  import { storeToRefs } from 'pinia'
  import useCodingStore from '@/store/modules/coding.ts'
  import useVisualDocumentStore from '@/store/modules/visual-document'
  import useVisualStore from '@/store/modules/visual'
  import useDebugStore from '@/store/modules/debug'
  import { DebugEventDispatcher } from '@/api/debug/debug-event-dispatcher'
  import { reqVisualData } from '@/components/code-visual/utils/index.ts'

  const documentPanel = ref()

  let codingStore = useCodingStore()
  let visualDocumentStore = useVisualDocumentStore()
  let visualStore = useVisualStore()
  let debugStore = useDebugStore()

  let { code, languages, language } = storeToRefs(codingStore)
  let { id, content, codeList } = storeToRefs(visualDocumentStore)
  let { action, descriptionType } = storeToRefs(visualStore)
  let { breakpoints } = storeToRefs(debugStore)

  let firstVisual = true

  const load = async () => {
    let result = await reqVisualDocument(id.value)
    if (result.code == 200) {
      let data = result.data
      content.value = data.content
      // 设置codeList
      let list: any[] = []
      for (let i = 0; i < result.data.codeList.length; i++) {
        const codeItem = result.data.codeList[i]
        list.push({
          code: codeItem.code,
          language: codeItem.language,
          breakpoints: codeItem.breakpoints,
        })
      }
      codeList.value = list
      // 设置可选的语言
      languages.value = []
      for (let code of codeList.value) {
        languages.value.push(code.language)
      }
      if (codeList.value.length > 0) {
        let firstCode = codeList.value[0]
        // 设置第一个语言
        language.value = languages.value[0]
        // 设置代码
        code.value = firstCode.code
        // 设置断点
        breakpoints.value = firstCode.breakpoints
      }
    }
  }

  const handleLanguageChange = () => {
    // 切换语言，则切换所有代码内容
    for (let visualCode of codeList.value) {
      if (visualCode.language == language.value) {
        action.value = true
        code.value = visualCode.code
        breakpoints.value = visualCode.breakpoints
      }
    }
  }

  const onStopped = async () => {
    // 判断是否需要跳转至可视化页面，如果开启可视化，且从当前调试未开启过可视化调试，并且有可视化数据，则跳转到可视化页面
    if (action.value && firstVisual) {
      let visualData = await reqVisualData(
        debugStore.id,
        descriptionType.value,
        visualStore.getDescription(descriptionType.value),
      )
      if (visualData.data.length != 0) {
        // 跳转至可视化调试页面
        firstVisual = false
      }
    }
  }

  const onLaunch = async () => {
    // 调试开始前清理可视化数据
    firstVisual = true
  }

  /**
   * 当拉伸/缩放可视化面板时调用
   * - 当可视化面板尺寸发生变化时，需重新调整可视化视图
   */
  const resizeVisualPane = () => {
    documentPanel.value?.resizeVisualView()
  }

  onMounted(() => {
    load()

    watch(
      () => id.value,
      () => {
        load()
      },
    )

    watch(
      () => language.value,
      () => {
        handleLanguageChange()
      },
    )

    // 注册一些事件
    DebugEventDispatcher.on('stopped', onStopped)
    DebugEventDispatcher.on('launch', onLaunch)
  })
  onUnmounted(() => {
    DebugEventDispatcher.off('stopped', onStopped)
    DebugEventDispatcher.on('launch', onLaunch)
  })
</script>

<style scoped lang="scss">
  .container {
    position: relative;
    width: 100%;
    height: calc(100vh - $header-height);
    padding: 0 8px 8px;
    margin: 0;
    background-color: $deep-background-color;
  }

  :deep(.left-panel) {
    margin-right: 4px;
    overflow: hidden;
    border: 1px solid $border-color;
    border-radius: $border-radius-large;
  }

  .left-content {
    display: flex;
    width: 100%;
    height: 100%;
  }

  .directory-section {
    width: 200px;
    overflow: hidden;
  }

  .document-section {
    width: calc(100% - 200px);
    overflow: hidden;
  }

  :deep(.coding-panel) {
    margin-left: 4px;
    overflow: hidden;
    border-radius: $border-radius-large;
  }

  .main {
    position: relative;
    height: 100%;
  }
</style>
