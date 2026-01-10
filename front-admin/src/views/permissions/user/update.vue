<template>
  <!--修改或添加用户-->
  <el-dialog v-model="visible" :title="getDialogTitle()">
    <el-form>
      <el-form-item label="用户id">
        <el-input placeholder="请输入用户id" v-model="userData.loginName"></el-input>
      </el-form-item>
      <el-form-item label="用户名称">
        <el-input placeholder="请输入用户名称" v-model="userData.username"></el-input>
      </el-form-item>
      <el-form-item label="邮箱">
        <el-input placeholder="请输入邮箱" v-model="userData.email"></el-input>
      </el-form-item>
      <el-form-item label="电话">
        <el-input placeholder="请输入电话" v-model="userData.phone"></el-input>
      </el-form-item>
      <el-form-item label="初始密码" v-show="type == 'insert'">
        <el-input placeholder="请输入密码" v-model="userData.password"></el-input>
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="visible = false">取消</el-button>
      <el-button type="primary" @click="submitAddOrUpdateUser">确定</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
  import { reactive, watch, computed, onMounted } from 'vue';
  import { ElMessage } from 'element-plus';
  import { reqInsertUser, reqUpdateUser, reqGetUser } from '@/api/user';

  const props = defineProps(['type', 'visible', 'userID']);
  const emit = defineEmits(['after-submit', 'update:visible']);
  // dialog是否可见
  let visible = computed({
    get() {
      return props.visible;
    },
    set(value) {
      emit('update:visible', value);
    },
  });

  // 修改或更新用户
  let userData = reactive({
    loginName: '',
    username: '',
    email: '',
    phone: '',
    password: '',
  });

  const getDialogTitle = () => {
    if (props.type == 'insert') {
      return '添加用户';
    } else {
      return '修改用户';
    }
  };

  // 读取用户信息
  const getUser = async (userID: string) => {
    if (props.type == 'update') {
      let result = await reqGetUser(userID);
      if (result.code == 200) {
        let data = result.data;
        userData.loginName = data.loginName;
        userData.username = data.username;
        userData.email = data.email;
        userData.phone = data.phone;
      } else {
        ElMessage({
          type: 'error',
          message: result.message,
        });
      }
    }
  };

  // 提交更新获取添加
  const submitAddOrUpdateUser = async () => {
    let result: any;
    if (props.type == 'insert') {
      result = await reqInsertUser(userData);
      if (result.code == 200) {
        visible.value = false;
        emit('after-submit');
      }
    } else {
      result = await reqUpdateUser({
        id: props.userID,
        loginName: userData.loginName,
        username: userData.username,
        email: userData.email,
        phone: userData.phone,
      });
      if (result.code == 200) {
        visible.value = false;
        emit('after-submit');
      }
    }
    ElMessage({
      showClose: true,
      message: result.message,
      type: result.code == 200 ? 'success' : 'error',
    });
  };

  onMounted(() => {
    if (props.type == 'update') {
      getUser(props.userID);
    }
  });

  watch(
    () => props.visible,
    () => {
      if (props.type == 'insert') {
        userData.loginName = '';
        userData.username = '';
        userData.email = '';
        userData.phone = '';
        userData.password = '';
      } else {
        getUser(props.userID);
      }
    },
  );
</script>

<style scoped lang="scss"></style>
