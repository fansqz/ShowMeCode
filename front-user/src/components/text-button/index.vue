<template>
  <pre
    class="button"
    @mouseover="isHovered = true"
    @mouseleave="isHovered = false"
    :class="{ hovered: isHovered }"
    >{{ text }}</pre
  >
</template>

<script setup lang="ts">
  import { ref } from 'vue'
  defineProps(['text'])
  const isHovered = ref(false)
</script>

<style scoped lang="scss">
  .button {
    color: $text-color-primary;
    cursor: pointer;
  }

  .hovered {
    color: $primary-color;
  }
</style>
