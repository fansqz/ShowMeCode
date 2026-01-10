<template>
  <svg :style="{ width, height }">
    <use :xlink:href="prefix + name" :fill="color"></use>
  </svg>
</template>

<script setup lang="ts">
  defineProps({
    prefix: {
      type: String,
      default: '#icon-',
    },
    // 提供使用图片的名称
    name: String,
    // 接收颜色
    color: {
      type: String,
      default: '',
    },
    // 接收宽高
    width: {
      typee: String,
      default: '16px',
    },
    height: {
      type: String,
      default: '16px',
    },
  });
</script>
