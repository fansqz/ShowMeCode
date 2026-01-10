<template>
  <el-collapse v-model="activeNames" class="visual-settings">
    <el-collapse-item name="1">
      <template #title>
        <div class="visual-settings__header" @click.stop>
          <el-text class="visual-settings__title">可视化配置</el-text>
          <div class="visual-settings__controls">
            <div class="control-group">
              <el-text class="control-label">可视化</el-text>
              <el-switch
                v-model="action"
                size="small"
                active-color="var(--el-color-primary)"
                inactive-color="var(--el-border-color)"
              />
            </div>
            <div class="control-divider"></div>
            <div class="mode-switch">
              <el-button
                :type="!isAIEnabled ? 'primary' : 'default'"
                size="small"
                @click.stop="toggleCustomMode"
                class="mode-btn"
              >
                <el-icon><Edit /></el-icon>
                编辑
              </el-button>
              <el-button
                :type="isAIEnabled ? 'primary' : 'default'"
                size="small"
                @click.stop="toggleAIMode"
                class="mode-btn"
              >
                <el-icon><Star /></el-icon>
                AI识别
              </el-button>
            </div>
          </div>
        </div>
      </template>
      <div class="visual-settings__content" :class="{ disabled: isAIEnabled }">
        <VisualTemplate />
      </div>
    </el-collapse-item>
  </el-collapse>
</template>

<script setup lang="ts">
  import useVisualStore from '@/store/modules/visual'
  import { ref } from 'vue'
  import { storeToRefs } from 'pinia'
  import VisualTemplate from './setting/index.vue'
  import { Edit, Star } from '@element-plus/icons-vue'

  // Store
  const visualStore = useVisualStore()
  const { action, descriptionType, isAIEnabled } = storeToRefs(visualStore)

  // Local state
  const activeNames = ref<string[]>([])

  /**
   * 切换到自定义编辑模式
   */
  const toggleCustomMode = () => {
    isAIEnabled.value = false
    activeNames.value = ['1']
  }

  /**
   * 切换到AI识别模式
   */
  const toggleAIMode = () => {
    isAIEnabled.value = true
    activeNames.value = []
  }
</script>

<style lang="scss" scoped>
  .visual-settings {
    position: absolute;
    right: 0;
    bottom: 0;
    left: 0;
    z-index: 10;
    width: 100%;
    height: auto;
    background-color: $base-background-color;
    border-radius: 8px;

    :deep(.el-collapse) {
      background-color: $base-background-color;
      border: 1px solid $border-color;
      border-radius: 8px;
    }

    :deep(.el-collapse-item__header) {
      background-color: $base-background-color;
      border-bottom: 1px solid $border-color;
    }

    :deep(.el-collapse-item__content) {
      background-color: $base-background-color;
    }

    &__header {
      display: flex;
      align-items: center;
      justify-content: space-between;
      width: 100%;
      padding: 0 20px;

      .visual-settings__title {
        margin: 0;
      }
    }

    &__controls {
      display: flex;
      gap: 12px;
      align-items: center;

      .control-group {
        display: flex;
        gap: 6px;
        align-items: center;

        .control-label {
          font-size: 12px;
          font-weight: 500;
          color: $text-color-secondary;
        }
      }

      .control-divider {
        width: 1px;
        height: 16px;
        background-color: $border-color-light;
      }

      .mode-switch {
        display: flex;
        gap: 4px;
        align-items: center;
        padding: 2px;
        background-color: $bg-color-secondary;
        border-radius: 6px;

        .mode-btn {
          height: 28px;
          padding: 0 12px;
          font-size: 12px;
          border: none;
          border-radius: 4px;
          transition: all 0.2s ease;

          &:hover {
            box-shadow: $box-shadow-light;
            transform: translateY(-1px);
          }

          .el-icon {
            margin-right: 4px;
            font-size: 12px;
          }
        }
      }
    }

    &__content {
      height: 30%;

      &.disabled {
        pointer-events: none;
        opacity: 0.6;

        :deep(input),
        :deep(textarea),
        :deep(select),
        :deep(button:not(.ai-control-button)) {
          pointer-events: none;
          user-select: none;
        }

        :deep(.el-input__inner),
        :deep(.el-textarea__inner) {
          color: $text-color-placeholder;
          cursor: not-allowed;
          background-color: $bg-color-secondary;
        }
      }
    }
  }
</style>
