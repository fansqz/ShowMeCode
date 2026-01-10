<template>
  <el-splitter layout="vertical">
    <el-splitter-panel size="70" class="editor-panel">
      <!--选择语言或者主题区域-->
      <EditorSelector class="editor-switcher" @showSavedCode="showSavedCodeModal" />
      <!--代码编辑区域-->
      <Editor class="editor" />
    </el-splitter-panel>
    <el-splitter-panel size="30" class="console-panel">
      <!--控制台-->
      <Console class="console-content" :userInput="false" :userOutput="false" :terminal="true" />
    </el-splitter-panel>
  </el-splitter>

  <!-- 保存代码模态框 -->
  <SavedCodeModal v-model="savedCodeModalVisible" :language="language" :documentID="documentID" />
</template>

<script setup lang="ts">
  import { ref } from 'vue'
  import { ElMessage } from 'element-plus'
  import { storeToRefs } from 'pinia'
  import Editor from '@/components/code-editor/editor/index.vue'
  import EditorSelector from './language-theme-switcher.vue'
  import Console from '@/components/code-editor/console/index.vue'
  import SavedCodeModal from './saved-user-code.vue'
  import useCodingStore from '@/store/modules/coding'
  import useVisualDocumentStore from '@/store/modules/visual-document'

  const codingStore = useCodingStore()
  const visualDocumentStore = useVisualDocumentStore()

  const { language } = storeToRefs(codingStore)
  const { id: documentID } = storeToRefs(visualDocumentStore)

  const savedCodeModalVisible = ref(false)

  /**
   * 显示保存代码模态框
   */
  const showSavedCodeModal = () => {
    savedCodeModalVisible.value = true
  }

  /**
   * 处理选择保存的代码
   */
  const handleSelectSavedCode = (savedCode: any) => {
    codingStore.code = savedCode.code
    ElMessage.success('代码已加载')
  }
</script>

<style scoped lang="scss">
  :deep(.editor-panel) {
    margin-bottom: 4px;
    overflow: hidden;
    border: 1px solid $border-color;
    border-radius: $border-radius-large;
  }

  :deep(.console-panel) {
    margin-top: 4px;
    overflow: hidden;
    border: 1px solid $border-color;
    border-radius: $border-radius-large;
  }

  .editor-switcher {
    position: relative;
    width: 100%;
    height: 35px;
  }

  .editor {
    position: relative;
    width: 100%;
    height: calc(100% - 35px);
  }

  .console-content {
    position: relative;
    height: 100%;
  }
</style>
