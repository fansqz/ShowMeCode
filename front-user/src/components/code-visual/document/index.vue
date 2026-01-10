<template>
  <div
    class="visaulDescription"
    :class="{ 'dark-mode': isDark }"
    :style="{ backgroundColor: isDark ? '#1d1e1f' : 'transparent' }"
  >
    <el-scrollbar>
      <v-md-preview :text="content" :theme="isDark ? 'dark' : 'light'"></v-md-preview>
    </el-scrollbar>
  </div>
</template>

<script setup lang="ts">
  import { ref, watchEffect } from 'vue'
  import { storeToRefs } from 'pinia'
  import useThemeStore from '@/store/modules/theme'

  let props = defineProps(['content'])
  let content = ref(props.content)

  const themeStore = useThemeStore()
  const { isDark } = storeToRefs(themeStore)

  watchEffect(() => {
    if (props.content != undefined) {
      content.value = props.content
    }
    // 调试信息
    console.log('Markdown component - isDark:', isDark.value)
  })
</script>

<style scoped lang="scss">
  .visaulDescription {
    width: 100%;
    height: 100%;
    background-color: $base-background-color;

    // Dark mode styles for markdown content
    &.dark-mode {
      background-color: #1d1e1f !important;

      // 覆盖el-scrollbar的背景
      :deep(.el-scrollbar__view) {
        background-color: #1d1e1f !important;
      }

      :deep(.v-md-editor-preview) {
        color: #e6e6e6 !important;
        background-color: #1d1e1f !important;
      }

      :deep(.v-md-editor-preview-wrapper) {
        background-color: #1d1e1f !important;
      }

      :deep(.v-md-editor-preview-wrapper .v-md-editor-preview) {
        color: #e6e6e6 !important;
        background-color: #1d1e1f !important;
      }

      // 强制覆盖所有可能的背景
      :deep(div) {
        background-color: #1d1e1f !important;
      }

      :deep(body) {
        background-color: #1d1e1f !important;
      }

      :deep(.v-md-editor-preview) {
        color: #e6e6e6 !important;
        background-color: #1d1e1f !important;

        // Headers
        h1,
        h2,
        h3,
        h4,
        h5,
        h6 {
          color: #fff !important;
          border-bottom-color: #404040 !important;
        }

        // Text
        p {
          color: #e6e6e6 !important;
        }

        // Links
        a {
          color: #4a9eff !important;

          &:hover {
            color: #6bb6ff !important;
          }
        }

        // Code blocks - 恢复代码样式
        pre {
          background-color: #0d1117 !important;
          border: 1px solid #30363d !important;

          code {
            color: #e6e6e6 !important;
            background-color: transparent !important;
          }
        }

        // Blockquotes
        blockquote {
          color: #e6e6e6 !important;
          background-color: #262727 !important;
          border-left: 4px solid #4a9eff !important;
        }

        // Tables
        table {
          background-color: #1d1e1f !important;

          th {
            color: #fff !important;
            background-color: #262727 !important;
            border-color: #404040 !important;
          }

          td {
            color: #e6e6e6 !important;
            border-color: #404040 !important;
          }
        }

        // Lists
        ul,
        ol {
          color: #e6e6e6 !important;
        }

        // Horizontal rule
        hr {
          border-color: #404040 !important;
        }

        // Strong and emphasis
        strong {
          color: #fff !important;
        }

        em {
          color: #e6e6e6 !important;
        }

        // Selection
        ::selection {
          color: #fff !important;
          background-color: #4a9eff !important;
        }
      }
    }
  }
</style>
