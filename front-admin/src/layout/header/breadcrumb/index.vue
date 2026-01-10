<template>
  <div class="menu-dev">
    <el-breadcrumb class="menu-path" separator-icon="ArrowRight">
      <el-breadcrumb-item v-for="(item, index) in $route.matched" :key="index">
        <span>{{ item.meta.title }}</span>
      </el-breadcrumb-item>
    </el-breadcrumb>
  </div>
</template>

<script setup lang="ts">
  import { useRoute } from 'vue-router';
  let $route = useRoute();
</script>

<style scoped>
  .menu-dev {
    position: relative;
    width: 100%;
    height: 100%;
    .menu-path {
      position: absolute;
      display: flex;
      align-items: center;
      position: absolute;
      bottom: 20%;
    }
  }
</style>
