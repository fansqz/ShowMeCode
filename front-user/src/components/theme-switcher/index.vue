<template>
  <div class="theme-switcher">
    <el-tooltip :content="isDark ? '切换到浅色模式' : '切换到暗黑模式'" placement="bottom">
      <div class="theme-toggle" @click="toggleTheme">
        <div class="toggle-track" :class="{ active: isDark }">
          <div class="toggle-thumb" :class="{ active: isDark }">
            <el-icon class="toggle-icon">
              <Sunny v-if="!isDark" />
              <Moon v-else />
            </el-icon>
          </div>
        </div>
      </div>
    </el-tooltip>
  </div>
</template>

<script setup lang="ts">
  import { storeToRefs } from 'pinia'
  import { Sunny, Moon } from '@element-plus/icons-vue'
  import useThemeStore from '@/store/modules/theme'

  const themeStore = useThemeStore()
  const { isDark } = storeToRefs(themeStore)

  const toggleTheme = () => {
    themeStore.toggleTheme()
  }
</script>

<style lang="scss" scoped>
  .theme-switcher {
    display: flex;
    align-items: center;
    justify-content: center;

    .theme-toggle {
      cursor: pointer;
      user-select: none;

      .toggle-track {
        position: relative;
        width: 48px;
        height: 24px;
        background: linear-gradient(135deg, #e8e8e8, #f0f0f0);
        border: 1px solid rgb(0 0 0 / 10%);
        border-radius: 12px;
        box-shadow:
          inset 0 1px 3px rgb(0 0 0 / 5%),
          0 1px 4px rgb(0 0 0 / 10%);
        transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);

        &.active {
          background: linear-gradient(135deg, #667eea, #764ba2);
          border-color: rgb(102 126 234 / 30%);
          box-shadow:
            inset 0 2px 4px rgb(255 255 255 / 20%),
            0 0 0 2px rgb(102 126 234 / 20%),
            0 4px 12px rgb(102 126 234 / 30%);
        }

        &:hover {
          transform: scale(1.02);

          &.active {
            box-shadow:
              inset 0 2px 4px rgb(255 255 255 / 20%),
              0 0 0 3px rgb(102 126 234 / 30%),
              0 6px 16px rgb(102 126 234 / 40%);
          }
        }
      }

      .toggle-thumb {
        position: absolute;
        top: 1px;
        left: 1px;
        width: 20px;
        height: 20px;
        background: linear-gradient(135deg, #fff, #f8f9fa);
        border: 1px solid rgb(0 0 0 / 10%);
        border-radius: 50%;
        box-shadow:
          0 1px 3px rgb(0 0 0 / 15%),
          0 1px 2px rgb(0 0 0 / 10%);
        transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);

        &.active {
          background: linear-gradient(135deg, #fff, #f0f4ff);
          border-color: rgb(102 126 234 / 30%);
          box-shadow:
            0 2px 4px rgb(102 126 234 / 20%),
            0 1px 2px rgb(0 0 0 / 10%);
          transform: translateX(24px);
        }

        .toggle-icon {
          position: absolute;
          top: 50%;
          left: 50%;
          font-size: 10px;
          color: #666;
          transition: all 0.3s ease;
          transform: translate(-50%, -50%);

          &.active {
            color: #667eea;
          }
        }
      }
    }
  }
</style>
