<template>
  <el-card class="box-card">
    <template #header>
      <div class="search">
        <el-row :gutter="20">
          <el-col :span="4">
            <el-input v-model="listQuery.loginName" clearable placeholder="请输入登录id" />
          </el-col>
          <el-col :span="4">
            <el-input v-model="listQuery.username" clearable placeholder="请输入用户名称" />
          </el-col>
          <el-col :span="4">
            <el-input v-model="listQuery.email" clearable placeholder="请输入邮箱" />
          </el-col>
          <el-col :span="4">
            <el-input v-model="listQuery.phone" clearable placeholder="请输入电话" />
          </el-col>
          <el-col :span="4">
            <el-button type="primary" @click="search"> 查询 </el-button>
            <el-button @click="reset"> 重置 </el-button>
          </el-col>
        </el-row>
      </div>
    </template>
    <el-button type="primary" @click="handleInsertUser">添加用户</el-button>
    <el-table border :data="userList" style="margin: 10px 0px">
      <el-table-column label="序号" width="80px" align="center" type="index"></el-table-column>
      <el-table-column label="登录id" prop="loginName" align="center"></el-table-column>
      <el-table-column label="用户名称" prop="username" align="center"></el-table-column>
      <el-table-column label="邮箱" prop="email" align="center"></el-table-column>
      <el-table-column label="电话" prop="phone" align="center"></el-table-column>
      <el-table-column label="角色列表" prop="roles" align="center"></el-table-column>
      <el-table-column label="更新时间" prop="updatedAt" align="center"></el-table-column>
      <el-table-column label="操作" width="220px" align="center">
        <!--row:为已有角色对象-->
        <template v-slot="{ row }">
          <el-button type="primary" size="small" icon="Plus" @click="handleSetUserRole(row.id)">
            关联角色
          </el-button>
          <el-button type="primary" size="small" icon="Edit" @click="handleUpdateUser(row.id)" />
          <el-popconfirm :title="`顶真要删除吗`" @confirm="handleDeleteUser(row.id)">
            <template #reference>
              <el-button type="danger" size="small" icon="Delete" />
            </template>
          </el-popconfirm>
        </template>
      </el-table-column>
    </el-table>
    <!--分页器组件-->
    <el-pagination
      @current-change="changePageNo"
      @size-change="changePageSize"
      v-model:current-page="listQuery.page"
      v-model:page-size="listQuery.pageSize"
      :page-sizes="[10, 20, 30, 50]"
      :background="true"
      layout="prev, pager, next, jumper, ->,sizes, total"
      :total="total"
      style="margin: 0px 3%"
    />
  </el-card>

  <RoleAssignDrawer
    v-model:visible="roleAssignDrawerData.visible"
    :userID="roleAssignDrawerData.userID"
    @afterSubmit="getUserList"
  />

  <UpdateDialog
    :userID="updateDialogData.userID"
    v-model:visible="updateDialogData.visible"
    :type="updateDialogData.type"
    @afterSubmit="getUserList"
  />
</template>

<script setup lang="ts">
  import { ref, onMounted, reactive } from 'vue';
  import { reqUserList, reqDeleteUser } from '@/api/user';
  import { ElMessage } from 'element-plus';
  import UpdateDialog from './update.vue';
  import RoleAssignDrawer from './role-assign.vue';

  let listQuery = reactive({
    username: '',
    loginName: '',
    email: '',
    phone: '',
    page: 0,
    pageSize: 10,
  });
  // 存储用户列表
  let total = ref(0);
  let userList = ref([]);

  // 更新用户的dialog的数据
  let updateDialogData = reactive({
    userID: '',
    type: 'insert',
    visible: false,
  });

  // 给用户分配角色的drawer的数据
  let roleAssignDrawerData = reactive({
    userID: '',
    visible: false,
  });

  //组件挂载完毕以后获取数据
  onMounted(() => {
    getUserList();
  });

  const getUserList = async () => {
    let result = await reqUserList(listQuery);
    if (result.code == 200) {
      total.value = result.data.total;
      userList.value = result.data.list;
    }
  };

  // 添加用户
  const handleInsertUser = () => {
    updateDialogData.type = 'insert';
    updateDialogData.visible = true;
  };

  // 修改用户
  const handleUpdateUser = (id: string) => {
    updateDialogData.userID = id;
    updateDialogData.type = 'update';
    updateDialogData.visible = true;
  };

  // 点击设置用户的角色
  const handleSetUserRole = (id: string) => {
    roleAssignDrawerData.userID = id;
    roleAssignDrawerData.visible = true;
  };

  // 删除角色
  const handleDeleteUser = async (id: string) => {
    let result = await reqDeleteUser(id);
    ElMessage({
      showClose: true,
      message: result.message,
      type: result.code == 200 ? 'success' : 'error',
    });
    getUserList();
  };

  const search = () => {
    listQuery.page = 1;
    getUserList();
  };

  const reset = () => {
    listQuery.username = '';
    listQuery.loginName = '';
    listQuery.email = '';
    listQuery.phone = '';
    listQuery.page = 1;
    getUserList();
  };

  //页码改变时触发
  const changePageNo = () => {
    // 页面变化时，页面归一
    getUserList();
  };

  const changePageSize = () => {
    listQuery.page = 1;
    getUserList();
  };
</script>

<style scoped lang="scss">
  .box-card {
    min-height: calc(100% - 82px);
    min-width: 800px;
    margin: 40px;
  }
</style>
