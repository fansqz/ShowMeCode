<template>
  <div class="document-editor-container">
    <v-md-editor
      v-model="content"
      height="100%"
      :disabled-menus="[]"
      @upload-image="handleUploadImage"
      left-toolbar="undo redo clear | h bold italic strikethrough quote | ul ol table hr | link image code"
    ></v-md-editor>
  </div>
</template>

<script setup lang="ts">
  import { storeToRefs } from 'pinia';
  import useVisualDocumentStore from '@/store/modules/visual-document.ts';
  import { reqUploadVisualDocumentFile } from '@/api/visual-document/index.ts';
  import { reqGetURL } from '@/api/common';
  const visualDocumentStore = useVisualDocumentStore();
  const { content } = storeToRefs(visualDocumentStore);

  const handleUploadImage = async (
    _event: any,
    insertImage: (arg0: { url: any }) => void,
    files: any[],
  ) => {
    let result = await reqUploadVisualDocumentFile({
      file: files[0],
    });
    if (result.code == 200) {
      // 读取url
      let result2 = await reqGetURL(result.data);
      if (result2.code == 200) {
        insertImage({ url: result2.data });
      }
    }
  };
</script>

<style scoped lang="scss">
  .document-editor-container {
    position: relative;
    height: 100%;
    width: 100%;
    .v-md-editor {
      box-shadow: 0 0px 0px #ccc !important;
    }
  }
</style>
