<template>
  <div class="menu">
    <div class="left">
      <!-- 语言选择 -->
      <el-select class="language-select" v-model="language" placeholder="Select" size="small">
        <el-option v-for="item in languages" :key="item" :value="item">
          <div class="language-item">{{ item }}</div>
        </el-option>
      </el-select>
    </div>

    <div class="right">
      <!-- 管理用户代码 -->
      <el-tooltip content="保存/读取 代码记录" placement="bottom">
        <el-icon @click="showSavedCodeModal" class="icon-button"><Document /></el-icon>
      </el-tooltip>
      <!-- 主题选择 -->
      <el-select class="theme-select" v-model="editorTheme" placeholder="Select" size="small">
        <el-option v-for="item in themeList" :key="item" :label="item" :value="item" />
      </el-select>
    </div>
  </div>
</template>

<script setup lang="ts">
  import { ref } from 'vue'
  import { storeToRefs } from 'pinia'
  import { getAllTheme } from '@/components/code-editor/editor/themes'
  import { ElMessage } from 'element-plus'
  import useCodingStore from '@/store/modules/coding'
  import useThemeStore from '@/store/modules/theme'
  import useUserStore from '@/store/modules/user'

  const emit = defineEmits<{
    (e: 'showSavedCode'): void
    (e: 'saveCurrentCode'): void
  }>()

  const codingStore = useCodingStore()
  const themeStore = useThemeStore()
  const userStore = useUserStore()

  // 主题列表
  const themeList = ref(getAllTheme())
  const { language, languages } = storeToRefs(codingStore)
  const { editorTheme } = storeToRefs(themeStore)

  // 显示保存代码模态框
  const showSavedCodeModal = () => {
    // 检查用户是否登录
    if (userStore.isLoggedIn()) {
      emit('showSavedCode')
    } else {
      // 未登录时显示提示
      ElMessage.warning('请先登录以使用代码管理功能')
    }
  }
</script>

<style scoped lang="scss">
  .menu {
    .left {
      position: absolute;
      box-sizing: border-box;
      display: flex;
      align-items: center;
      width: 100%;
      height: 35px;
      background-color: $base-header-background;
      border-bottom: 1px solid $border-color;

      .language-select {
        width: 150px;
        margin-right: 7px;
        margin-left: 30px;

        .language-item {
          user-select: none;
        }
      }
    }

    .right {
      position: absolute;
      right: 10px;
      display: flex;
      gap: 12px;
      align-items: center;
      padding: 5px 0;

      .icon-button {
        font-size: 18px;
        color: $text-color;
        cursor: pointer;
        transition: color 0.2s;

        &:hover {
          color: $primary-color;
        }
      }

      .theme-select {
        width: 120px;
        margin-left: 10px;
      }
    }
  }
</style>
