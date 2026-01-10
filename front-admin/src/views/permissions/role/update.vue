<template>
  <!--修改或添加角色-->
  <el-dialog v-model="visible" :title="getDialogTitle()">
    <el-form>
      <el-form-item label="角色名称">
        <el-input placeholder="请输入角色名称" v-model="roleData.name"></el-input>
      </el-form-item>
      <el-form-item label="介绍">
        <el-input placeholder="请输入角色简介" v-model="roleData.description"></el-input>
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="visible = false">取消</el-button>
      <el-button type="primary" @click="submitAddOrUpdateRole">确定</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
  import { reactive, computed, watch } from 'vue';
  import { ElMessage } from 'element-plus';
  import { reqInsertRole, reqUpdateRole, reqGetRole } from '@/api/role';

  const props = defineProps(['type', 'visible', 'roleID']);
  const emit = defineEmits(['update:visible', 'afterSubmit']);
  let visible = computed({
    get() {
      return props.visible;
    },
    set(value) {
      emit('update:visible', value);
    },
  });
  let roleData = reactive({
    name: '',
    description: '',
  });

  const submitAddOrUpdateRole = async () => {
    let result;
    if (props.type == 'insert') {
      result = await reqInsertRole(roleData);
    } else {
      result = await reqUpdateRole({
        ...roleData,
        id: props.roleID,
      });
    }
    if (result.code == 200) {
      ElMessage({
        showClose: true,
        message: props.type == 'insert' ? '角色添加成功' : '角色更新成功',
        type: 'success',
      });
      visible.value = false;
      emit('afterSubmit');
    } else {
      ElMessage({
        showClose: true,
        message: props.type == 'insert' ? '角色添加失败' : '角色更新失败',
        type: 'error',
      });
    }
  };

  const getDialogTitle = (): string => {
    return props.type == 'insert' ? '角色添加' : '角色更新';
  };

  const getRole = async (roleID: string) => {
    let result = await reqGetRole(roleID);
    if (result.code == 200) {
      roleData.name = result.data.name;
      roleData.description = result.data.description;
    }
  };

  watch(
    () => props.visible,
    () => {
      if (props.visible) {
        if (props.type == 'insert') {
          roleData.name = '';
          roleData.description = '';
        } else {
          getRole(props.roleID);
        }
      }
    },
  );
</script>

<style scoped lang="scss"></style>
