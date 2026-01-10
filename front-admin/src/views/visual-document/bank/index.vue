<template>
  <el-card class="box-card">
    <!--搜索框-->
    <template #header>
      <!--顶部添加题库-->
      <el-button type="primary" size="default" @click="handlerInsertBank"> 添加知识库 </el-button>
    </template>
    <!--展示题库列表-->
    <el-table style="margin: 10px 0px" border :data="bankList">
      <el-table-column label="序号" width="80px" align="center" type="index"></el-table-column>
      <el-table-column label="知识库名称" align="center">
        <template v-slot="{ row }">
          <el-button @click="handlerVisualDocumentManage(row.id)" type="primary" link>{{
            row.name
          }}</el-button>
        </template>
      </el-table-column>
      <el-table-column label="创建时间" width="160px" align="center">
        <template v-slot="{ row }">
          <pre>{{ row.createdAt }}</pre>
        </template>
      </el-table-column>
      <el-table-column label="创建者" align="center">
        <template v-slot="{ row }">
          <pre>{{ row.creatorName }}</pre>
        </template>
      </el-table-column>
      <el-table-column label="启用" align="center">
        <template v-slot="{ row }">
          <el-text>{{ row.enable == true ? '是' : '否' }}</el-text>
        </template>
      </el-table-column>
      <el-table-column label="操作" align="center" width="120px">
        <template v-slot="{ row }">
          <el-button type="primary" size="small" icon="Edit" @click="handlerUpdateBank(row.id)" />
          <el-popconfirm :title="`顶真要删除吗`" @confirm="handleDeleteBank(row.id)">
            <template #reference>
              <el-button type="danger" size="small" icon="Delete" />
            </template>
          </el-popconfirm>
        </template>
      </el-table-column>
    </el-table>

    <BankUpdate
      v-model:visible="updateDialogData.visible"
      :type="updateDialogData.type"
      :bankID="updateDialogData.bankID"
      @afterSubmit="getAllVisualDocumentBankList"
    />
  </el-card>
</template>

<script setup lang="ts">
  import {
    reqAllVisualDocumentBank,
    reqDeleteVisualDocumentBank,
  } from '@/api/visual-document-bank';
  import BankUpdate from './bank-editor.vue';
  import { ref, onMounted, reactive } from 'vue';
  import { useRouter } from 'vue-router';
  import { ElMessage } from 'element-plus';

  let $router = useRouter();

  let bankList = ref([]);

  // 管理更新或者添加题库的dialog
  let updateDialogData = reactive({
    bankID: '',
    type: 'insert',
    visible: false,
  });

  const getAllVisualDocumentBankList = async () => {
    let result = await reqAllVisualDocumentBank();
    if (result.code == 200) {
      bankList.value = result.data;
    }
  };

  const handleDeleteBank = async (id: string) => {
    let result = await reqDeleteVisualDocumentBank(id);
    ElMessage({
      showClose: true,
      message: result.data,
      type: 'success',
    });
    getAllVisualDocumentBankList();
  };

  // 更新题库
  const handlerUpdateBank = (id: string) => {
    updateDialogData.bankID = id;
    updateDialogData.type = 'update';
    updateDialogData.visible = true;
  };

  // 添加题库
  const handlerInsertBank = () => {
    updateDialogData.bankID = '';
    updateDialogData.type = 'insert';
    updateDialogData.visible = true;
  };

  const handlerVisualDocumentManage = (bankID: string) => {
    $router.push({
      name: 'visual-document',
      params: {
        bankID: bankID,
      },
    });
  };

  //组件挂载完毕以后获取数据
  onMounted(() => {
    getAllVisualDocumentBankList();
  });
</script>

<style scoped>
  .box-card {
    margin: 40px;
    min-width: 800px;
    min-height: calc(100% - 82px);
  }
</style>
