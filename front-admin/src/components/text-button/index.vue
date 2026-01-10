<template>
  <pre
    class="button"
    @mouseover="style = 'color:blue'"
    @mouseleave="style = 'color:#4089ef'"
    :style="style"
    >{{ text }}</pre
  >
</template>

<script setup lang="ts">
  import { ref } from 'vue';
  defineProps(['text']);
  let style = ref();
</script>

<style scoped lang="scss">
  .button {
    cursor: pointer;
  }
</style>
