<template>
  <div class="container">
    <CodeEditor
      class="code-editor"
      v-model:code="model.code"
      v-model:language="model.language"
      v-model:breakpoints="model.breakpoints"
    />
  </div>
</template>

<script setup lang="ts">
  import CodeEditor from './code/index.vue';
  import { defineModel } from 'vue';

  const model = defineModel<{
    code: string;
    language: string;
    breakpoints: number[];
  }>({
    default: {
      code: '',
      language: 'c',
      breakpoints: [],
    },
  });
</script>

<style scoped lang="scss">
  .container {
    position: relative;
    height: 100%;
    width: 100%;
    .code-editor {
      position: relative;
      height: 100%;
      width: 100%;
    }
  }
</style>
