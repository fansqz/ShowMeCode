<template>
  <div class="container">
    <div class="setting">
      <el-avatar v-if="getAvatar() != ''" :src="getAvatar()" class="avatar" />
      <el-avatar v-if="getAvatar() == ''" class="avatar">s</el-avatar>
      <el-dropdown class="user_dropdown">
        <span class="el-dropdown-link">
          {{ getUsername() }}<el-icon class="el-icon--right"><arrow-down /></el-icon>
        </span>
        <template #dropdown v-if="isLogged()">
          <el-dropdown-menu>
            <!-- <el-dropdown-item @click="changeRoute('myprofile')">个人中心</el-dropdown-item> -->
            <el-dropdown-item @click="changeRoute('accountSetting')">信息修改</el-dropdown-item>
            <el-dropdown-item @click="logout">退出登录</el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </div>
  </div>
</template>

<script setup lang="ts">
  import useUserStore from '@/store/modules/user'
  import { useRouter, useRoute } from 'vue-router'
  let userStore = useUserStore()
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
</script>

<style scoped>
  .container {
    position: relative;
    width: 100%;
    height: 100%;
    margin: 0;

    .setting {
      position: absolute;
      display: flex;
      flex-direction: row;
      align-items: center;
      justify-content: end;
      width: 100%;
      height: 100%;
      margin-right: 30px;

      .avatar {
        width: 24px;
        height: 24px;
        margin-right: 10px;
        border-radius: 50%;
      }

      .user_dropdown {
        cursor: pointer;
      }
    }
  }
</style>
