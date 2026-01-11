<template>
  <div class="setting">
    <el-avatar v-if="getAvatar()" :src="getAvatar()" :size="28" />
    <el-avatar v-else :size="28">{{ getUsername().charAt(0) }}</el-avatar>
    <el-dropdown>
      <span class="user-trigger">
        {{ getUsername() }}
        <el-icon class="el-icon--right"><arrow-down /></el-icon>
      </span>
      <template #dropdown>
        <el-dropdown-menu>
          <el-dropdown-item @click="changeRoute('accountSetting')">信息修改</el-dropdown-item>
          <el-dropdown-item @click="toggleTheme">
            {{ isDark ? '浅色模式' : '深色模式' }}
          </el-dropdown-item>
          <el-dropdown-item divided @click="logout">退出登录</el-dropdown-item>
        </el-dropdown-menu>
      </template>
    </el-dropdown>
  </div>
</template>

<script setup lang="ts">
  import useUserStore from '@/store/modules/user'
  import useThemeStore from '@/store/modules/theme'
  import { storeToRefs } from 'pinia'
  import { useRouter, useRoute } from 'vue-router'
  let userStore = useUserStore()
  const themeStore = useThemeStore()
  const { isDark } = storeToRefs(themeStore)
  let $router = useRouter()
  let $route = useRoute()

  const getUsername = (): string => {
    return userStore.token ? userStore.username : '未登录'
  }

  const isLogged = (): boolean => {
    return !!userStore.token
  }

  const getAvatar = (): string => {
    return userStore.token ? userStore.avatar : ''
  }

  const logout = () => {
    userStore.userLogout()
    $router.push({ path: '/login', query: { redirect: $route.path } })
  }

  const changeRoute = (routeName: string, params = {}) => {
    if ($route.name === routeName) {
      return
    }
    $router.push({ name: routeName, params })
  }

  const toggleTheme = () => {
    themeStore.toggleTheme()
  }
</script>

<style scoped lang="scss">
  .setting {
    display: flex;
    align-items: center;
    gap: 8px;

    .user-trigger {
      display: flex;
      align-items: center;
      cursor: pointer;
      color: var(--el-text-color-primary);

      &:hover {
        color: var(--el-color-primary);
      }
    }
  }
</style>
