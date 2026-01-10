<template>
  <div class="visual-document-container">
    <el-menu
      :default-active="activeIndex"
      class="menu"
      mode="horizontal"
      :append-to-body="false"
      @select="handleSelect"
    >
      <el-menu-item class="text-edit-space"></el-menu-item>
      <el-menu-item index="document">文本编辑</el-menu-item>
      <el-menu-item index="code">代码编辑</el-menu-item>
    </el-menu>

    <div class="title-editor" disabled>
      <div class="text" v-if="!inEditTitle">
        <el-text>{{ title }}</el-text>
        <el-button class="iconfont icon-edit1" link @click.native="handleClickEditTitle()" />
      </div>
      <div class="text-edit" v-if="inEditTitle">
        <el-input size="small" v-model="title">
          <template #suffix>
            <el-button size="small" type="primary" @click="UpdateTitle()" link>确定</el-button>
          </template>
        </el-input>
      </div>
    </div>

    <div class="enable-button">
      <el-switch class="button" v-model="enable" />
    </div>

    <div class="save-button">
      <el-button v-if="!hasChange" type="info" size="small" disabled>保存</el-button>
      <el-button v-if="hasChange" type="primary" size="small" @click="saveVisualDocument"
        >保存</el-button
      >
    </div>

    <div class="delete-button">
      <el-popconfirm :title="`顶真要删除吗`" @confirm="deleteVisualDocument()">
        <template #reference>
          <el-button type="danger" size="small">删除</el-button>
        </template>
      </el-popconfirm>
    </div>

    <div class="editor-main">
      <DocumentEditor v-if="activeIndex == 'document'" />
      <CodeEditor v-if="activeIndex == 'code'" />
    </div>
  </div>
</template>

<script setup lang="ts">
  import { onMounted, watch, ref } from 'vue';
  import {
    reqVisualDocument,
    reqUpdateVisualDocument,
    reqDeleteVisualDocument,
  } from '@/api/visual-document';
  import DocumentEditor from './document-editor/index.vue';
  import CodeEditor from './code-editor/index.vue';
  import useVisualDocumentStore from '@/store/modules/visual-document';
  import { storeToRefs } from 'pinia';
  import { ElMessage } from 'element-plus';

  const visualDocumentStore = useVisualDocumentStore();

  const hasChange = ref(false);
  const emit = defineEmits(['updatedTitle', 'handleDocumentNotFound', 'handleDocumentDelete']);

  const activeIndex = ref('document');
  const { bankID, id, parentID, title, content, enable, codeList } =
    storeToRefs(visualDocumentStore);
  const inEditTitle = ref(false);

  // 获取document
  const getVisualDocument = async () => {
    let result = await reqVisualDocument(id.value);
    if (result.code == 10000) {
      // 记录为找到
      emit('handleDocumentNotFound');
    }
    if (result.code == 200) {
      bankID.value = result.data.bankID;
      parentID.value = result.data.parentID;
      title.value = result.data.title;
      content.value = result.data.content;
      enable.value = result.data.enable;
      let list: any[] = [];
      for (let i = 0; i < result.data.codeList.length; i++) {
        const codeItem = result.data.codeList[i];
        list.push({
          code: codeItem.code,
          language: codeItem.language,
          breakpoints: codeItem.breakpoints,
        });
      }
      codeList.value = list;
    }
  };

  const handleSelect = (key: string) => {
    activeIndex.value = key;
  };

  // 更新标题
  const UpdateTitle = async () => {
    inEditTitle.value = false;
  };

  // 删除可视化文档
  const deleteVisualDocument = async () => {
    let result = await reqDeleteVisualDocument(id.value);
    if (result.code == 200) {
      ElMessage({
        showClose: true,
        message: '删除成功',
        type: 'success',
      });
      emit('handleDocumentDelete');
    }
  };

  // 保存可视化文档
  const saveVisualDocument = async () => {
    let codeListReq = [];
    for (let data of codeList.value) {
      codeListReq.push({
        code: data.code,
        language: data.language,
        breakpoints: data.breakpoints,
      });
    }
    let result = await reqUpdateVisualDocument({
      id: id.value,
      parentID: parentID.value,
      title: title.value,
      content: content.value,
      enable: enable.value,
      codeList: codeListReq,
    });
    if (result.code == 200) {
      ElMessage({
        showClose: true,
        message: '保存成功',
        type: 'success',
      });
      // 更新目录标题
      emit('updatedTitle');
      hasChange.value = false;
    }
  };

  const handleClickEditTitle = () => {
    inEditTitle.value = true;
  };

  onMounted(async () => {
    await getVisualDocument();

    watch(
      () => visualDocumentStore,
      () => {
        hasChange.value = true;
      },
      {
        deep: true,
      },
    );
  });
</script>

<style scoped lang="scss">
  .container {
    position: relative;
    height: 100%;
    width: 100%;
    .menu {
      top: 0px;
      height: 40px;
      width: 100%;
      .text-edit-space {
        width: 150px;
        margin: 0px 20px;
      }
    }
    .title-editor {
      position: absolute;
      top: 0px;
      display: flex;
      align-self: center;
      width: 150px;
      height: 40px;
      margin: 0px 20px;
      .text {
        display: flex;
        justify-content: space-between;
        width: 150px;
        height: 40px;
      }
      .text-edit {
        margin: auto 0px;
        width: 150px;
      }
    }
    .enable-button {
      position: absolute;
      top: 0px;
      right: 140px;
      width: 40px;
      height: 40px;
      display: flex;
      flex-direction: column;
      justify-content: center;
    }
    .save-button {
      position: absolute;
      top: 0px;
      right: 80px;
      height: 40px;
      display: flex;
      flex-direction: column;
      justify-content: center;
    }
    .delete-button {
      position: absolute;
      top: 0px;
      right: 20px;
      height: 40px;
      display: flex;
      flex-direction: column;
      justify-content: center;
    }
    .editor-main {
      position: absolute;
      top: 40px;
      left: 0px;
      height: calc(100% - 40px);
      width: 100%;
    }
  }
</style>
