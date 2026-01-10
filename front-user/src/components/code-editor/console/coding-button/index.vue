<template>
  <ul class="horizontal-list">
    <li class="li">
      <DebugButton v-if="debug" />
    </li>
    <li class="li">
      <ExecuteButton v-if="execute" />
    </li>
  </ul>
</template>

<script setup lang="ts">
  defineProps<{
    debug: boolean
    execute: boolean
  }>()

  import DebugButton from '../debug-button/debug-button.vue'
  import ExecuteButton from './execute-button.vue'
</script>

<style lang="scss" scoped>
  .horizontal-list {
    position: relative; /* 设置层级需定位 */
    box-sizing: border-box;
    display: flex;
    align-items: center;
    justify-content: flex-end;
    padding: 0;

    .li {
      position: relative; /* 设置层级需定位 */
      display: flex;
      align-items: center;
      justify-content: center;
      margin-right: 18px;
    }
  }
</style>
