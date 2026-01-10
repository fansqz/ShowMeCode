<template>
  <!--修改或添加题库-->
  <el-dialog
    class="dialog"
    v-model="visible"
    :title="type == 'update' ? '知识库更新' : '知识库修改'"
  >
    <el-form>
      <el-form-item label="知识库名称">
        <el-input placeholder="请输入知识库名称" v-model="bankData.name"></el-input>
      </el-form-item>
      <el-form-item label="描述">
        <el-input placeholder="请输入描述" v-model="bankData.description"></el-input>
      </el-form-item>
      <el-form-item label="是否启用" label-width="auto">
        <el-switch
          v-model="bankData.enable"
          :active-value="true"
          :inactive-value="false"
          style="--el-switch-on-color: #13ce66; --el-switch-off-color: #646765"
        />
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="visible = false">取消</el-button>
      <el-button type="primary" @click="submitAddOrUpdateBank">确定</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
  import { onMounted, reactive, watch, computed } from 'vue';
  import {
    reqUpdateVisualDocumentBank,
    reqVisualDocumentBank,
    reqInsertVisualDocumentBank,
  } from '@/api/visual-document-bank';
  import { ElMessage } from 'element-plus';

  // type = update or insert
  const props = defineProps(['type', 'bankID', 'visible']);
  const emit = defineEmits(['update:visible', 'afterSubmit']);
  // 题库数据
  let bankData = reactive<any>({
    name: '',
    description: '',
    enable: false,
  });
  // dialog是否可见
  const visible = computed({
    get() {
      return props.visible;
    },
    set(value) {
      emit('update:visible', value);
    },
  });

  const getVisualDocumentBank = async (bankID: number) => {
    if (props.type == 'update') {
      let result = await reqVisualDocumentBank(bankID);
      if (result.code == 200) {
        let data = result.data;
        bankData.name = data.name;
        bankData.description = data.description;
      } else {
        ElMessage({
          type: 'error',
          message: result.message,
        });
      }
    }
  };

  const submitAddOrUpdateBank = async () => {
    let result: any;
    if (props.type == 'insert') {
      result = await reqInsertVisualDocumentBank(bankData);
    } else {
      result = await reqUpdateVisualDocumentBank({
        id: props.bankID,
        name: bankData.name,
        description: bankData.description,
        enable: bankData.enable,
      });
    }
    if (result.code == 200) {
      ElMessage({
        type: 'success',
        message: result.message,
      });
      visible.value = false;
      emit('afterSubmit');
    } else {
      ElMessage({
        type: 'error',
        message: result.message,
      });
    }
  };

  //组件挂载完毕以后获取数据
  onMounted(() => {
    if (props.type == 'update') {
      getVisualDocumentBank(props.bankID);
    }
  });

  watch(
    () => props.visible,
    () => {
      if (props.type == 'insert') {
        bankData.name = '';
        bankData.description = '';
      } else {
        getVisualDocumentBank(props.bankID);
      }
    },
  );
</script>

<style lang="scss" scoped>
  .dialog {
    width: 300px;
  }
</style>
