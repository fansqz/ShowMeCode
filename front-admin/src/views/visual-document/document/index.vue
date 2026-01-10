<template>
  <div class="container">
    <div class="root-node">
      <el-button
        class="back iconfont icon-back"
        link
        size="large"
        @click="gotoBankList"
      ></el-button>
      <el-text class="bank-name" size="large">{{ bank.name }}</el-text>
      <el-button link size="samll" class="button" @click="clickAddVisualDocument(0)"
        >+</el-button
      ></div
    >
    <el-scrollbar class="directory-tree">
      <el-tree
        :data="direcotry"
        :props="props"
        highlight-current="true"
        :expand-on-click-node="false"
        :current-node-key="id"
        draggable
        node-key="id"
        indent="8"
        @node-click="clickEditorVisualDocument"
        @node-drag-end="handleUpdateDocumentParent"
      >
        <template #default="scope">
          <div class="tree_item">
            <el-text>{{ scope.node.label }}</el-text>
            <el-button
              link
              size="samll"
              class="button"
              @click="clickAddVisualDocument(scope.node.data)"
              >+</el-button
            >
          </div>
        </template>
      </el-tree>
    </el-scrollbar>
    <div v-if="!notVisualDocumentData" class="editor">
      <VisualDocumentEditor
        :key="id"
        @updatedTitle="getDirectory()"
        @handleDocumentDelete="selectFirstDocument()"
        @handleDocumentNotFound="selectFirstDocument()"
      />
    </div>
    <div v-if="notVisualDocumentData" class="not-record">
      <el-text>无数据</el-text>
    </div>

    <el-dialog v-model="visualDocumentForAdd.visible" title="新建可视化文档">
      <el-input v-model="visualDocumentForAdd.title" placeholder="未命名文档" />
      <template #footer>
        <div class="dialog-footer">
          <el-button
            @click="
              () => {
                visualDocumentForAdd.visible = false;
              }
            "
            >取消</el-button
          >
          <el-button @click="addVisualDocument()" type="primary">提交</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
  import {
    reqVisualDocumentDirectory,
    reqInsertVisualDocument,
    reqUpdateVisualDocumentDirectory,
  } from '@/api/visual-document';
  import { onMounted, reactive, ref } from 'vue';
  import VisualDocumentEditor from './document-editor/index.vue';
  import { defaultDocument } from './document-editor/document-editor/default-document';
  import useVisualDocumentStore from '@/store/modules/visual-document';
  import { reqVisualDocumentBank } from '@/api/visual-document-bank/index.ts';
  import { useRoute } from 'vue-router';
  import { useRouter } from 'vue-router';
  import { storeToRefs } from 'pinia';
  const $route = useRoute();
  let $router = useRouter();
  const visualDocumentStore = useVisualDocumentStore();

  const convertToNumber = (input: string | string[]): number => {
    let value: string;
    if (Array.isArray(input)) {
      value = input[0];
    } else {
      value = input;
    }
    return Number(value);
  };

  let bankID = convertToNumber($route.params.bankID);
  const { id } = storeToRefs(visualDocumentStore);
  const direcotry = ref<any[]>([]);
  const bank = reactive({
    bankID: bankID,
    name: '',
  });
  const props = {
    children: 'children',
    label: 'title',
  };
  // 标记是否没有可视化文档数据
  const notVisualDocumentData = ref(false);
  // 添加可视化文档的数据
  const visualDocumentForAdd = reactive({
    bankID: bankID,
    visible: false,
    parentID: 0,
    title: '',
    content: '',
    enable: false,
  });

  const clickEditorVisualDocument = (node: any) => {
    id.value = node.id;
  };

  const clickAddVisualDocument = (node: any) => {
    visualDocumentForAdd.parentID = node.id;
    visualDocumentForAdd.title = '未命名';
    visualDocumentForAdd.content = defaultDocument;
    visualDocumentForAdd.enable = false;
    visualDocumentForAdd.visible = true;
  };

  // 更新文档父节点
  const handleUpdateDocumentParent = async (draggingNode: any, dropNode: any, dropType: any) => {
    if (draggingNode.data.id == dropNode.data.id) {
      return;
    }
    await reqUpdateVisualDocumentDirectory(draggingNode.data.id, dropNode.data.id, dropType);
    // 更新目录
    getDirectory();
  };

  // 更新目录并选择第一篇文档
  const selectFirstDocument = async () => {
    await getDirectory();
    if (direcotry.value.length != 0) {
      id.value = direcotry.value[0].id;
    } else {
      // 目录无数据
      notVisualDocumentData.value = true;
    }
  };

  // 添加可视化文档
  const addVisualDocument = async () => {
    let result = await reqInsertVisualDocument({
      bankID: visualDocumentForAdd.bankID,
      parentID: visualDocumentForAdd.parentID,
      title: visualDocumentForAdd.title,
      content: visualDocumentForAdd.content,
      enable: visualDocumentForAdd.enable,
    });
    if ((result.code = 200)) {
      // 如果一开始是无数据状态，设置为有数据
      if (notVisualDocumentData.value == true) {
        notVisualDocumentData.value = false;
      }
      // 获取新的目录
      await getDirectory();
      id.value = result.data;
    }
    visualDocumentForAdd.visible = false;
  };

  // gotoBankList 返回题库列表
  const gotoBankList = () => {
    $router.push({ name: 'visual-document-bank', params: {} });
  };

  // 获取目录
  const getDirectory = async () => {
    let result = await reqVisualDocumentDirectory(bankID);
    if (result.code == 200) {
      direcotry.value = result.data;
    }
  };

  onMounted(async () => {
    await getDirectory();
    // 如果id为0，或者store中存储的文档所属题库和当前页面不同，那么重新加载文档
    if (id.value == 0 || bankID != visualDocumentStore.bankID) {
      if (direcotry.value.length != 0) {
        id.value = direcotry.value[0].id;
      } else {
        // 目录无数据
        notVisualDocumentData.value = true;
      }
    }
    // 读取读取知识库库
    let result = await reqVisualDocumentBank(bankID);
    if (result.code == 200) {
      bank.name = result.data.name;
    }
  });
</script>

<style scoped lang="scss">
  .container {
    position: relative;
    height: calc(100% - 30px);
    width: calc(100% - 20px);
    box-shadow: 0px 0px 10px rgb(220, 220, 220);
    border-radius: 5px;
    min-width: 800px;
    .root-node {
      position: relative;
      top: 0px;
      width: 220px;
      height: 40px;
      border-right: 1px solid $base-border-color;
      border-bottom: 1px solid $base-border-color;
      box-sizing: border-box;
      display: flex;
      align-items: center; /* 竖直居中 */
      .back {
        width: 30px !important;
      }
      .button {
        position: absolute;
        right: 6px;
      }
    }
    .directory-tree {
      position: absolute;
      top: 40px;
      width: 220px;
      height: calc(100% - 40px);
      border-right: 1px solid $base-border-color;
      box-sizing: border-box;
      .tree_item {
        width: 100%;
        height: 26px;
        display: flex;
        flex-direction: row;
        justify-content: space-between;
        .button {
          margin-right: 6px;
        }
      }
    }
    .editor {
      position: absolute;
      width: calc(100% - 220px);
      height: 100%;
      top: 0px;
      left: 220px;
      margin: 0px;
    }
    .not-record {
      position: absolute;
      width: calc(100% - 220px);
      height: 100%;
      top: 0px;
      left: 220px;
      margin: 0px;
      display: flex;
      flex-direction: column;
      justify-content: center;
    }
  }
</style>
