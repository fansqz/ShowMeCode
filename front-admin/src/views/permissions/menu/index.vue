<template>
  <el-card class="box-card">
    <el-table :data="menuTree" style="width: 100%; margin-bottom: 20px" row-key="id" border>
      <el-table-column label="名称" prop="name"></el-table-column>
      <el-table-column label="权限值" prop="code"></el-table-column>
      <el-table-column label="描述" prop="description"></el-table-column>
      <el-table-column label="更新时间" prop="updatedAt" width="160px"></el-table-column>
      <el-table-column label="操作" width="160px">
        <!--row:为已有的菜单对象-->
        <template v-slot="{ row }">
          <el-button type="primary" size="small" icon="Plus" @click="handleInsertMenu(row)" />
          <el-button
            type="primary"
            size="small"
            icon="Edit"
            :disabled="row.parentMenuID == 0"
            @click="handleUpdateMenu(row)"
          />
          <el-popconfirm :title="`顶真要删除吗`" @confirm="handleDeleteMenu(row)">
            <template #reference>
              <el-button
                type="danger"
                size="small"
                icon="Delete"
                :disabled="row.parentMenuID == 0"
              />
            </template>
          </el-popconfirm>
        </template>
      </el-table-column>
    </el-table>

    <!--修改或添加-->
    <el-dialog v-model="dialogVisible" :title="getDialogTitle()">
      <el-form>
        <el-form-item label="菜单名称">
          <el-input placeholder="请输入菜单名称" v-model="menuData.name"></el-input>
        </el-form-item>
        <el-form-item label="权限值">
          <el-input placeholder="请输入请求路径" v-model="menuData.code"></el-input>
        </el-form-item>
        <el-form-item label="概述">
          <el-input placeholder="请输入菜单概述" v-model="menuData.description"></el-input>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitAddOrUpdateMenu">确定</el-button>
      </template>
    </el-dialog>
  </el-card>
</template>

<script setup lang="ts">
  import { ref, onMounted, reactive } from 'vue';
  import { reqGetMenuTree, reqInsertMenu, reqUpdateMenu, reqDeleteMenu } from '@/api/menu';
  import { ElMessage } from 'element-plus';

  let menuTree = ref([]);
  // menu数据
  let menuData = reactive({
    id: '',
    name: '',
    code: '',
    description: '',
    parentMenuID: 0,
  });
  // 添加menu/修改menu页面是否显示
  let dialogVisible = ref(false);
  // 用于标识是添加还是修改menu，1标识添加，0表示修改
  let dialogType = ref('insert');
  //请求
  onMounted(() => {
    getMenuTree();
  });

  // 获取menu树
  const getMenuTree = async () => {
    let result = await reqGetMenuTree();
    if (result.code == 200) {
      menuTree.value = result.data;
    }
  };

  // 添加菜单
  const handleInsertMenu = (row: any) => {
    dialogType.value = 'insert';
    dialogVisible.value = true;
    menuData.parentMenuID = row.id;
    menuData.id = '';
    menuData.name = '';
    menuData.code = '';
    menuData.description = '';
  };

  // 修改菜单
  const handleUpdateMenu = (row: any) => {
    dialogType.value = 'update';
    dialogVisible.value = true;
    menuData.id = row.id;
    menuData.name = row.name;
    menuData.code = row.code;
    menuData.description = row.description;
    menuData.parentMenuID = row.parentMenuID;
  };

  // 提交 添加菜单/修改菜单 的请求
  const submitAddOrUpdateMenu = async () => {
    if (dialogType.value == 'insert') {
      const result = await reqInsertMenu(menuData);
      if (result.code == 200) {
        dialogVisible.value = false;
        getMenuTree();
        ElMessage({
          showClose: true,
          message: '菜单添加成功',
          type: 'success',
        });
      } else {
        ElMessage({
          showClose: true,
          message: '菜单添加失败',
          type: 'error',
        });
      }
    } else {
      const result = await reqUpdateMenu(menuData);
      if (result.code == 200) {
        dialogVisible.value = false;
        getMenuTree();
        ElMessage({
          showClose: true,
          message: '菜单更新成功',
          type: 'success',
        });
      } else {
        ElMessage({
          showClose: true,
          message: '菜单更新失败',
          type: 'error',
        });
      }
    }
  };

  const getDialogTitle = () => {
    return dialogType.value == 'insert' ? '添加菜单' : '修改菜单';
  };

  // 删除menu
  const handleDeleteMenu = async (row: any) => {
    let result = await reqDeleteMenu(row.id);
    if (result.code == 200) {
      getMenuTree();
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
</script>

<style scoped lang="scss">
  .box-card {
    min-height: calc(100% - 82px);
    min-width: 800px;
    margin: 40px;
  }
</style>
