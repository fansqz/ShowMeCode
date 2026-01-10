<template>
  <div class="setting">
    <el-dropdown :hide-on-click="false" class="user-dropdown">
      <span class="el-dropdown-link">
        {{ userStore.username }}<el-icon class="el-icon--right"><arrow-down /></el-icon>
      </span>
      <template #dropdown>
        <el-dropdown-menu>
          <el-dropdown-item @click="logout">退出登录</el-dropdown-item>
        </el-dropdown-menu>
      </template>
    </el-dropdown>
    <el-avatar v-if="getAvatar() != ''" :src="getAvatar()" class="avatar" />
    <el-avatar v-if="getAvatar() == ''" class="avatar">s</el-avatar>
  </div>
</template>

<script setup lang="ts">
  import useUserStore from '@/store/modules/user';
  import { useRouter, useRoute } from 'vue-router';
  let userStore = useUserStore();
  let $router = useRouter();
  let $route = useRoute();
  const logout = () => {
    userStore.userLogout();
    $router.push({ path: '/login', query: { redirect: $route.path } });
  };

  const getAvatar = (): string => {
    return userStore.token ? userStore.avatar : '';
  };
</script>

<style scoped>
  .setting {
    position: absolute;
    display: flex;
    height: 100%;
    width: 100%;
    align-items: center;
    flex-direction: row-reverse;
    .user-dropdown {
      padding-right: 60px;
    }
    .avatar {
      margin-right: 10px;
      width: 24px;
      height: 24px;
      border-radius: 50%;
    }
  }
</style>
