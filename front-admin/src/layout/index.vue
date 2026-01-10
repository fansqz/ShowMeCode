<template>
  <div class="layout-container">
    <!--左侧菜单-->
    <div class="layout-slider">
      <div class=""></div>
      <el-scrollbar class="scrollbar">
        <el-menu :default-active="$route.path">
          <Menu :menuList="menu"></Menu>
        </el-menu>
      </el-scrollbar>
    </div>

    <!--顶部tabbar-->
    <div class="layout-header">
      <Header></Header>
    </div>

    <!--展示栏-->
    <div class="layout-main">
      <Main></Main>
    </div>
  </div>
</template>

<script setup lang="ts">
  import { reactive } from 'vue';
  import { useRoute } from 'vue-router';
  import Menu from './menu/index.vue';
  // 获取用户仓库
  import useUserStore from '@/store/modules/user';
  import Main from './main/index.vue';
  import Header from './header/index.vue';

  let userStore = useUserStore();
  let $route = useRoute();
  let menu: any[] = reactive([]);
  userStore.menuRoutes.forEach((element: any) => {
    if ((element.name == 'layout' || element.name == 'manage') && element.children) {
      menu = [...menu, ...element.children];
    }
  });
</script>

<style scoped lang="scss">
  .layout-container {
    width: 100vw;
    height: 100vh;

    .layout-slider {
      box-sizing: border-box;
      position: relative;
      width: $base-menu-width;
      height: calc(100vh - $base-header-height);
      background: $base-menu-background;
      top: $base-header-height;
      border-right: 1px solid $base-border-color;
      .scrollbar {
        position: fixed;
        width: $base-menu-width - 1px;
        height: calc(100% - $base-header-height);
        bottom: 0%;
      }
    }
    .layout-header {
      position: fixed;
      width: 100%;
      top: 0;
    }
    .layout-main {
      box-sizing: border-box;
      position: fixed;
      width: calc(100vw - $base-menu-width);
      height: calc(100vh - $base-header-height);
      top: $base-header-height;
      left: $base-menu-width;
      overflow: auto;
    }
  }
</style>
