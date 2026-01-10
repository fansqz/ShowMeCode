<template>
  <!--分配接口与菜单-->
  <el-drawer v-model="visible">
    <template #header>
      <h4>分配权限</h4>
    </template>
    <template #default>
      <div>
        <div inlink="true">
          <el-radio size="large" v-model="type" label="apiTree"> 接口管理 </el-radio>
          <el-radio size="large" v-model="type" label="menuTree"> 菜单管理 </el-radio>
        </div>

        <el-tree
          v-show="type == 'apiTree'"
          :data="apiData"
          ref="apiTree"
          show-checkbox
          default-expand-all
          check-strictly
          node-key="id"
          :default-checked-keys="apiSelecteds"
          :props="defaultProps"
        />
        <el-tree
          v-show="type == 'menuTree'"
          :data="menuData"
          ref="menuTree"
          show-checkbox
          default-expand-all
          check-strictly
          node-key="id"
          :default-checked-keys="menuSelecteds"
          :props="defaultProps"
        />
      </div>
    </template>
    <template #footer>
      <div style="flex: auto">
        <el-button @click="visible = false">取消</el-button>
        <el-button type="primary" @click="submitRolePermisstion">提交</el-button>
      </div>
    </template>
  </el-drawer>
</template>

<script setup lang="ts">
  import { ref, watch, computed } from 'vue';
  import { ElMessage } from 'element-plus';
  import { reqGetApiTree } from '@/api/api';
  import { reqGetMenuTree } from '@/api/menu';
  import { reqRoleMenu, reqRoleApi, reqUpdateRoleApi, reqUpdateRoleMenu } from '@/api/role';

  const props = defineProps(['visible', 'roleID']);
  const emit = defineEmits(['update:visible', 'afterSubmit']);
  let visible = computed({
    get() {
      return props.visible;
    },
    set(value) {
      emit('update:visible', value);
    },
  });
  // 展示menu树/ 展示api树
  let menuData = ref([]);
  let apiData = ref([]);
  // 已选
  let apiSelecteds = ref([]);
  let menuSelecteds = ref([]);
  let type = ref('apiTree');
  let menuTree = ref();
  let apiTree = ref();

  const defaultProps = {
    chidren: 'children',
    label: 'name',
  };

  // 提交添加权限请求
  const submitRolePermisstion = async () => {
    // 用户选择的
    let apiIDs = apiTree.value.getCheckedKeys();
    let menuIDs = menuTree.value.getCheckedKeys();
    let result1 = await reqUpdateRoleApi({ roleID: props.roleID, apiIDs: apiIDs });
    let result2 = await reqUpdateRoleMenu({ roleID: props.roleID, menuIDs: menuIDs });
    if (result1.code == 200 && result2.code == 200) {
      ElMessage({
        showClose: true,
        message: '提交成功',
        type: 'success',
      });
      //关闭抽屉
      visible.value = false;
      emit('afterSubmit');
    } else {
      ElMessage({
        showClose: true,
        message: '操作失败',
        type: 'error',
      });
    }
  };

  const loadData = async () => {
    let result = await reqGetApiTree();
    if (result.code == 200) {
      apiData.value = result.data;
    }
    result = await reqGetMenuTree();
    if (result.code == 200) {
      menuData.value = result.data;
    }
    result = await reqRoleApi(props.roleID);
    if (result.code == 200) {
      apiSelecteds.value = result.data;
    }
    result = await reqRoleMenu(props.roleID);
    if (result.code == 200) {
      menuSelecteds.value = result.data;
    }
  };

  watch(
    () => props.visible,
    () => {
      if (props.visible) {
        loadData();
      }
    },
  );
</script>

<style scoped lang="scss"></style>
