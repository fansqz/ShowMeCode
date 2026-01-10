<template>
  <el-card class="box-card">
    <el-table :data="apiTree" style="width: 100%; margin-bottom: 20px" row-key="id" border>
      <el-table-column label="名称" prop="name"></el-table-column>
      <el-table-column label="路径" prop="path"></el-table-column>
      <el-table-column label="请求方式" prop="method" width="100"></el-table-column>
      <el-table-column label="更新时间" prop="updatedAt" width="160"></el-table-column>
      <el-table-column label="操作" width="160px">
        <template v-slot="{ row }">
          <el-button type="primary" size="small" icon="Plus" @click="handleInsertApi(row)" />
          <el-button
            type="primary"
            size="small"
            icon="Edit"
            :disabled="row.parentApiID == 0"
            @click="handleUpdateApi(row)"
          />
          <el-popconfirm :title="`顶真要删除吗`" @confirm="handleDeleteApi(row)">
            <template #reference>
              <el-button
                type="danger"
                size="small"
                icon="Delete"
                :disabled="row.parentApiID == 0"
              />
            </template>
          </el-popconfirm>
        </template>
      </el-table-column>
    </el-table>

    <!--修改或添加-->
    <el-dialog v-model="dialogVisible" :title="getDialogTitle()">
      <el-form>
        <el-form-item label="接口名称">
          <el-input placeholder="请输入接口名称" v-model="apiData.name"></el-input>
        </el-form-item>
        <el-form-item label="请求路径">
          <el-input placeholder="请输入请求路径" v-model="apiData.path"></el-input>
        </el-form-item>
        <el-form-item label="请求方法">
          <el-input placeholder="请输入请求方法" v-model="apiData.method"></el-input>
        </el-form-item>
        <el-form-item label="概述">
          <el-input placeholder="请输入接口概述" v-model="apiData.description"></el-input>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitAddOrUpdateApi">确定</el-button>
      </template>
    </el-dialog>
  </el-card>
</template>

<script setup lang="ts">
  import { ref, onMounted, reactive } from 'vue';
  import { reqGetApiTree, reqInsertApi, reqUpdateApi, reqDeleteApi } from '@/api/api';
  import { ElMessage } from 'element-plus';

  let apiTree = ref([]);
  // api数据
  let apiData = reactive({
    id: '',
    name: '',
    path: '',
    method: '',
    description: '',
    parentApiID: 0,
  });
  // 添加api/修改api页面是否显示
  let dialogVisible = ref(false);
  // 用于标识是添加还是修改api
  let dialogType = ref('insert');

  //请求
  onMounted(() => {
    getApiTree();
  });

  // 获取api树
  const getApiTree = async () => {
    let result = await reqGetApiTree();
    if (result.code == 200) {
      apiTree.value = result.data;
    }
  };

  // 添加接口
  const handleInsertApi = (row: any) => {
    dialogType.value = 'insert';
    dialogVisible.value = true;
    apiData.id = '';
    apiData.name = '';
    apiData.parentApiID = row.id;
    apiData.path = row.path;
    apiData.method = '';
    apiData.description = '';
  };

  // 修改接口
  const handleUpdateApi = (row: any) => {
    dialogType.value = 'update';
    dialogVisible.value = true;
    apiData.id = row.id;
    apiData.name = row.name;
    apiData.path = row.path;
    apiData.method = row.method;
    apiData.description = row.description;
    apiData.parentApiID = row.parentApiID;
  };

  // 删除api
  const handleDeleteApi = async (row: any) => {
    let result = await reqDeleteApi(row.id);
    if (result.code == 200) {
      getApiTree();
      ElMessage({
        showClose: true,
        message: '删除成功',
        type: 'success',
      });
    } else {
      ElMessage({
        showClose: true,
        message: '删除失败',
        type: 'error',
      });
    }
  };

  // 提交 添加接口/修改接口 的请求
  const submitAddOrUpdateApi = async () => {
    if (dialogType.value == 'insert') {
      const result = await reqInsertApi(apiData);
      if (result.code == 200) {
        dialogVisible.value = false;
        getApiTree();
        ElMessage({
          showClose: true,
          message: '接口添加成功',
          type: 'success',
        });
      } else {
        ElMessage({
          showClose: true,
          message: '接口添加失败',
          type: 'error',
        });
      }
    } else {
      const result = await reqUpdateApi(apiData);
      if (result.code == 200) {
        dialogVisible.value = false;
        getApiTree();
        ElMessage({
          showClose: true,
          message: '接口更新成功',
          type: 'success',
        });
      } else {
        ElMessage({
          showClose: true,
          message: '接口更新失败',
          type: 'error',
        });
      }
    }
  };

  const getDialogTitle = () => {
    return dialogType.value == 'insert' ? '添加接口' : '修改接口';
  };
</script>

<style scoped lang="scss">
  .box-card {
    min-height: calc(100% - 82px);
    min-width: 800px;
    margin: 40px;
  }
</style>
