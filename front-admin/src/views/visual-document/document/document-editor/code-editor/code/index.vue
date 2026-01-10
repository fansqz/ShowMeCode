<template>
  <div ref="editorEl" class="editor" />
</template>

<script lang="ts" setup>
  import { ref, watch, onMounted, defineModel } from 'vue';
  import { editor } from 'monaco-editor';
  import { setBreakPoint, initBreakPoint } from './breakpoint';
  import { getConfigs } from './conf';
  import debounce from 'lodash.debounce';

  let editorInstance: editor.IStandaloneCodeEditor;
  let code = defineModel<string>('code', { default: '' });
  let language = defineModel<string>('language', { default: 'c' });
  let breakpoints = defineModel<number[]>('breakpoints', { default: [] });

  const editorEl = ref<HTMLElement>();

  onMounted(() => {
    if (editorEl.value == undefined) {
      return;
    }
    // 创建editor实例
    editorInstance = editor.create(
      editorEl.value,
      getConfigs({
        language: language.value,
        value: code.value,
        readOnly: false,
        minimap: {
          enabled: false,
        },
      }),
    );

    // 初始化断点
    initBreakPoint(editorInstance, breakpoints.value);

    watch(
      () => language.value,
      async (val) => {
        const model = editorInstance.getModel();
        if (model) {
          //设置语言
          editor.setModelLanguage(model, val);
        }
      },
    );

    watch(
      () => code.value,
      (val) => {
        if (code.value != editorInstance.getValue()) {
          editorInstance?.setValue(val);
        }
      },
    );

    // 监控内容修改
    editorInstance.onDidChangeModelContent(
      debounce(async () => {
        // 使用防抖，在不输入内容的时候进行保存
        code.value = editorInstance.getValue();
      }, 600),
    );

    setBreakPoint(editorInstance, function (bps: number[], _lineNum: number, _mode: 'add' | 'del') {
      // 设置断点
      breakpoints.value = bps;
    });
  });
</script>

<style scoped lang="scss">
  @import './style/breakPoint.scss';
  @import './style/debug.scss';
  .editor {
    position: absolute;
    width: 100%;
    height: 100%;
    max-height: 100% !important;
    margin: 0;
    padding: 0;
    overflow: hidden;
  }
</style>
