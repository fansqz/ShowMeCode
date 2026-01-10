<template>
  <!--分配角色-->
  <el-drawer v-model="visible">
    <template #header>
      <h4>分配角色</h4>
    </template>
    <template #default>
      <el-form>
        <el-form-item label="用户名称："> {{ username }} </el-form-item>
        <el-form-item lable="角色列表">
          <el-checkbox v-model="checkAll" :indeterminate="indeterminate" @change="handleSelectAll">
            全选
          </el-checkbox>
          <el-checkbox-group v-model="selectedRoleIDs" @change="handleRoleSelectChange">
            <el-checkbox v-for="(role, index) in roles" :key="index" :label="role.id">
              {{ role.name }}
            </el-checkbox>
          </el-checkbox-group>
        </el-form-item>
      </el-form>
    </template>
    <template #footer>
      <div style="flex: auto">
        <el-button @click="visible = false">取消</el-button>
        <el-button type="primary" @click="submitUpdateUserRole">提交</el-button>
      </div>
    </template>
  </el-drawer>
</template>

<script setup lang="ts">
  import { ref, computed, watch } from 'vue';
  import { ElMessage } from 'element-plus';
  import { reqUpdateUserRole, reqGetUser, reqSimpleRoleList, reqUserRole } from '@/api/user';
  const props = defineProps(['userID', 'visible']);
  const emit = defineEmits(['update:visible', 'afterSubmit']);
  // drawer是否可见
  let visible = computed({
    get() {
      return props.visible;
    },
    set(value) {
      emit('update:visible', value);
    },
  });
  let username = ref('');
  // 是否选中所有
  let checkAll = ref(false);
  let indeterminate = ref<boolean>(true);
  // 可选用户列表
  let roles = ref<any[]>([]);
  // 已选用户id
  let selectedRoleIDs = ref<string[]>([]);

  const handleSelectAll = (val: boolean) => {
    selectedRoleIDs.value = val ? roles.value.map((role: any) => role.id) : [];
    indeterminate.value = false;
  };

  const handleRoleSelectChange = () => {
    if (selectedRoleIDs.value == roles.value.map((role: any) => role.id)) {
      indeterminate.value = false;
    } else {
      indeterminate.value = true;
    }
  };

  // 提交添加权限请求
  const submitUpdateUserRole = async () => {
    let result = await reqUpdateUserRole(props.userID, selectedRoleIDs.value);
    if (result.code == 200) {
      visible.value = false;
      emit('afterSubmit');
    }
    ElMessage({
      showClose: true,
      message: result.message,
      type: result.code == 200 ? 'success' : 'error',
    });
  };

  watch(
    () => props.userID,
    async () => {
      let result = await reqGetUser(props.userID);
      if (result.code == 200) {
        username.value = result.data.username;
      }
    },
  );

  watch(
    () => props.visible,
    async () => {
      let result = await reqSimpleRoleList();
      if (result.code == 200) {
        roles.value = result.data;
      }
      result = await reqUserRole(props.userID);
      if (result.code == 200) {
        selectedRoleIDs.value = result.data;
      }
    },
  );
</script>

<style scoped lang="scss"></style>
