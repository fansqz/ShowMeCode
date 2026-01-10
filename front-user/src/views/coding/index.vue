<template>
  <div class="coding-workspace">
    <el-splitter layout="horizontal">
      <el-splitter-panel
        class="coding-workspace__visual-panel"
        @update:size="handleVisualPanelResize"
      >
        <VisualPanel ref="visualPanel" />
      </el-splitter-panel>

      <el-splitter-panel class="coding-workspace__editor-panel">
        <el-splitter layout="vertical">
          <el-splitter-panel size="70" class="coding-workspace__code-editor">
            <LanguageThemeSelector
              class="coding-workspace__language-selector"
              @showSavedCode="showSavedCodeModal"
            />
            <CodeEditor class="coding-workspace__editor" @onChangeValue="handleCodeChange" />
          </el-splitter-panel>
          <el-splitter-panel size="30" class="coding-workspace__console">
            <Console
              :userInput="false"
              :userOutput="false"
              :terminal="true"
              class="coding-workspace__console-content"
            />
          </el-splitter-panel>
        </el-splitter>
      </el-splitter-panel>
    </el-splitter>
  </div>

  <!-- 保存代码模态框 -->
  <SavedCodeModal v-model="savedCodeModalVisible" :language="language" />

  <!-- 调试指引 -->
  <DebugGuide />
</template>

<script setup lang="ts">
  import { ref, onMounted, watch } from 'vue'
  import 'splitpanes/dist/splitpanes.css'
  import VisualPanel from './visual.vue'
  import CodeEditor from '@/components/code-editor/editor/index.vue'
  import LanguageThemeSelector from './language-theme-switcher.vue'
  import Console from '@/components/code-editor/console/index.vue'
  import SavedCodeModal from './saved-user-code.vue'
  import DebugGuide from './guide.vue'
  import { storeToRefs } from 'pinia'
  import useCodingStore from '@/store/modules/coding.ts'
  import useDebugStore from '@/store/modules/debug'
  import { languageConstants, supportedLanguages } from '@/constants/languages.ts'

  // Store instances
  const codingStore = useCodingStore()
  const debugStore = useDebugStore()

  // Store refs
  const { code, language, languages } = storeToRefs(codingStore)
  const { breakpoints } = storeToRefs(debugStore)

  // Component refs
  const visualPanel = ref<InstanceType<typeof VisualPanel> | null>()
  const savedCodeModalVisible = ref(false)

  // Types
  type CodeConfig = {
    code: string
    language: languageConstants
    breakpoints: number[]
  }

  // State
  const codeConfig: CodeConfig = {
    code: '',
    language: languageConstants.C,
    breakpoints: [],
  }

  // Initialize
  language.value = languageConstants.GO
  languages.value = supportedLanguages

  /**
   * 加载代码配置
   */
  const loadCodeConfig = async () => {
    // 尝试从本地存储读取配置
    const localConfigJson = localStorage.getItem(`code-${language.value}`)
    let configLoaded = false

    if (localConfigJson) {
      try {
        Object.assign(codeConfig, JSON.parse(localConfigJson))
        configLoaded = true
      } catch (error) {
        console.warn('Failed to parse local code config:', error)
      }
    }

    // 如果本地配置无效，从服务器加载
    if (!configLoaded) {
      try {
        const [configsJson, codeText] = await Promise.all([
          fetch('/document/visual-learn-config.json').then((res) => res.text()),
          fetch(`/document/visual-learn-code.${language.value}`).then((res) => res.text()),
        ])

        const configs: CodeConfig[] = JSON.parse(configsJson)
        const targetConfig = configs.find((config) => config.language === language.value)

        if (targetConfig) {
          Object.assign(codeConfig, targetConfig)
        }

        codeConfig.code = codeText
      } catch (error) {
        console.error('Failed to load code config from server:', error)
      }
    }

    // 应用配置
    code.value = codeConfig.code
    breakpoints.value = codeConfig.breakpoints
  }

  /**
   * 保存代码配置到本地存储
   */
  const saveCodeConfig = () => {
    const configToSave = {
      code: code.value,
      language: language.value,
      breakpoints: breakpoints.value,
    }
    localStorage.setItem(`code-${language.value}`, JSON.stringify(configToSave))
  }

  /**
   * 处理代码变化
   */
  const handleCodeChange = (_value: string, _type: string) => {
    saveCodeConfig()
  }

  /**
   * 处理语言切换
   */
  const handleLanguageChange = () => {
    loadCodeConfig()
  }

  /**
   * 处理可视化面板大小变化
   */
  const handleVisualPanelResize = () => {
    visualPanel.value?.resizeVisualCanvas()
  }

  /**
   * 显示保存代码模态框
   */
  const showSavedCodeModal = () => {
    savedCodeModalVisible.value = true
  }

  // Watchers
  watch(
    () => language.value,
    () => {
      handleLanguageChange()
    },
  )

  // Lifecycle
  onMounted(() => {
    loadCodeConfig()
  })
</script>

<style scoped lang="scss">
  .coding-workspace {
    position: relative;
    height: calc(100vh - $header-height);
    padding: 0 8px 8px;
    background-color: $deep-background-color;

    // Splitter panels styling
    :deep(.coding-workspace__visual-panel) {
      margin-right: 4px;
      overflow: hidden;
      border: 1px solid $border-color;
      border-radius: $border-radius-large;
    }

    :deep(.coding-workspace__code-editor) {
      margin-bottom: 4px;
      margin-left: 4px;
      overflow: hidden;
      border: 1px solid $border-color;
      border-radius: $border-radius-large;
    }

    :deep(.coding-workspace__console) {
      margin-top: 4px;
      margin-left: 4px;
      overflow: hidden;
      border: 1px solid $border-color;
      border-radius: $border-radius-large;
    }

    // Component styling
    &__language-selector {
      position: relative;
      width: 100%;
      height: 35px;
    }

    &__editor {
      position: relative;
      width: 100%;
      height: calc(100% - 35px);
    }

    &__console-content {
      position: relative;
      height: 100%;
    }
  }
</style>
