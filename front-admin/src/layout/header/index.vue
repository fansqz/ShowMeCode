<template>
  <div class="tabbar">
    <Logo class="logo"></Logo>
    <div class="tabbar-right">
      <Setting></Setting>
    </div>
  </div>
</template>

<script setup lang="ts">
  import Logo from '../logo/index.vue';
  import Setting from './setting/index.vue';
</script>

<style scoped lang="scss">
  .tabbar {
    position: relative;
    height: $base-header-height;
    background-color: $base-header-background;
    box-sizing: border-box;
    border-bottom: 1px solid $base-border-color;
    width: 100%;
    top: 0px;
    .logo {
      display: block;
      height: $base-header-height;
      width: 200px;
      margin: 0px 20px;
    }
    .tabbar-right {
      position: absolute;
      height: $base-header-height;
      width: 30%;
      right: 0%;
      top: 0%;
    }
  }
</style>
