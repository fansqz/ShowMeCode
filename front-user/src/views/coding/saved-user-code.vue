<template>
  <div v-if="visible" class="saved-code-overlay" @click="handleOverlayClick">
    <div class="saved-code-card" @click.stop>
      <div class="card-header">
        <h3 class="card-title">{{ `代码备份 - ${currentLanguage}` }}</h3>
        <el-button type="text" @click="handleClose" class="close-btn">
          <el-icon><Close /></el-icon>
        </el-button>
      </div>

      <div class="card-body">
        <!-- 保存新代码区域 -->
        <div class="save-new-code">
          <el-button
            type="primary"
            @click="showSaveDialog = true"
            :loading="saving"
            size="small"
            plain
          >
            保存当前代码
          </el-button>
        </div>

        <div class="code-list">
          <el-table
            :data="savedCodeList"
            style="width: 100%"
            v-loading="loading"
            :empty-text="loading ? '加载中...' : '暂无保存的代码'"
            height="350"
            max-height="350"
          >
            <el-table-column prop="createdAt" label="创建时间" width="120">
              <template #default="{ row }">
                {{ formatTime(row.createdAt) }}
              </template>
            </el-table-column>

            <el-table-column prop="updatedAt" label="更新时间" width="120">
              <template #default="{ row }">
                {{ formatTime(row.updatedAt) }}
              </template>
            </el-table-column>

            <el-table-column prop="remark" label="备注" min-width="150">
              <template #default="{ row }">
                <span v-if="row.remark">{{ row.remark }}</span>
                <span v-else class="no-remark">无备注</span>
              </template>
            </el-table-column>

            <el-table-column label="操作" width="180" fixed="right">
              <template #default="{ row }">
                <el-button
                  type="primary"
                  size="small"
                  @click.stop="handleApply(row)"
                  :loading="applyingId === row.id"
                  class="apply-btn"
                >
                  应用
                </el-button>
                <el-button
                  type="danger"
                  size="small"
                  @click.stop="handleDelete(row)"
                  :loading="deletingId === row.id"
                  class="delete-btn"
                >
                  删除
                </el-button>
              </template>
            </el-table-column>
          </el-table>
        </div>
      </div>
    </div>
  </div>

  <!-- 保存代码对话框 -->
  <el-dialog
    v-model="showSaveDialog"
    title="保存代码"
    width="500px"
    :before-close="handleSaveDialogClose"
    center
  >
    <div class="save-dialog-content">
      <el-form :model="saveForm" label-width="40px">
        <el-form-item label="语言">
          <el-input v-model="saveForm.language" disabled />
        </el-form-item>
        <el-form-item label="备注">
          <el-input
            v-model="saveForm.remark"
            placeholder="请输入备注信息（可选）"
            type="textarea"
            :rows="3"
          />
        </el-form-item>
      </el-form>
    </div>

    <template #footer>
      <div class="dialog-footer">
        <el-button @click="handleSaveDialogClose">取消</el-button>
        <el-button type="primary" @click="handleSaveCode" :loading="saving"> 保存 </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
  import { ref, computed, watch, onMounted } from 'vue'
  import { storeToRefs } from 'pinia'
  import { ElMessage, ElMessageBox } from 'element-plus'
  import { Close } from '@element-plus/icons-vue'
  import {
    reqGetSavedCodeList,
    reqCreateSavedCode,
    reqGetSavedCodeDetail,
    reqDeleteSavedCode,
  } from '@/api/saved-code'
  import type { SavedCode, CreateSavedCodeRequest } from '@/api/saved-code/type'
  import useCodingStore from '@/store/modules/coding'
  const codingStore = useCodingStore()
  const { code } = storeToRefs(codingStore)

  interface Props {
    modelValue: boolean
    language: string
    documentID?: number
  }

  interface Emits {
    (e: 'update:modelValue', value: boolean): void
  }

  const props = defineProps<Props>()
  const emit = defineEmits<Emits>()

  // 响应式数据
  const visible = computed({
    get: () => props.modelValue,
    set: (value) => emit('update:modelValue', value),
  })

  const currentLanguage = computed(() => props.language)
  const savedCodeList = ref<SavedCode[]>([])
  const loading = ref(false)
  const saving = ref(false)
  const loadingDetail = ref(false)
  const deletingId = ref<number | null>(null)
  const applyingId = ref<number | null>(null)
  const showSaveDialog = ref(false)
  const saveForm = ref<CreateSavedCodeRequest>({
    language: props.language,
    code: code.value,
    remark: '',
  })

  // 格式化时间
  const formatTime = (timeStr: string) => {
    if (!timeStr) return '-'
    const date = new Date(timeStr)
    return date.toLocaleDateString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
    })
  }

  // 获取保存的代码列表
  const loadSavedCodeList = async () => {
    try {
      loading.value = true
      const params = {
        language: props.language,
        documentID: props.documentID,
        page: 1,
        pageSize: 50,
      }
      const response = await reqGetSavedCodeList(params)
      if (response.code === 200) {
        savedCodeList.value = response.data.list
      } else {
        ElMessage.error(response.message || '获取代码列表失败')
        savedCodeList.value = []
      }
    } catch (error) {
      console.error('获取代码列表失败:', error)
      ElMessage.error('获取代码列表失败')
      savedCodeList.value = []
    } finally {
      loading.value = false
    }
  }

  // 应用代码
  const handleApply = async (row: SavedCode) => {
    try {
      applyingId.value = row.id

      // 根据ID获取完整的代码内容
      const response = await reqGetSavedCodeDetail(row.id)

      if (response.code === 200) {
        const codeDetail = response.data
        console.log('应用代码详情:', codeDetail)
        code.value = codeDetail.code
        ElMessage.success('代码已应用到编辑器')
        handleClose()
      } else {
        ElMessage.error(response.message || '获取代码详情失败')
      }
    } catch (error) {
      console.error('应用代码失败:', error)
      ElMessage.error('应用代码失败')
    } finally {
      applyingId.value = null
    }
  }

  // 关闭对话框
  const handleClose = () => {
    visible.value = false
  }

  // 处理遮罩层点击
  const handleOverlayClick = () => {
    handleClose()
  }

  // 处理保存对话框关闭
  const handleSaveDialogClose = () => {
    showSaveDialog.value = false
    // 重置表单
    saveForm.value = {
      language: props.language,
      code: code.value,
      remark: '',
    }
  }

  // 删除代码
  const handleDelete = async (row: SavedCode) => {
    try {
      // 显示确认对话框
      await ElMessageBox.confirm(
        `确定要删除这条代码记录吗？${row.remark ? `\n备注：${row.remark}` : ''}`,
        '确认删除',
        {
          confirmButtonText: '确定删除',
          cancelButtonText: '取消',
          type: 'warning',
          customClass: 'delete-confirm-dialog',
        },
      )

      deletingId.value = row.id

      const response = await reqDeleteSavedCode(row.id)

      if (response.code === 200) {
        ElMessage.success('删除成功')
        // 重新加载列表
        await loadSavedCodeList()
      } else {
        ElMessage.error(response.message || '删除失败')
      }
    } catch (error) {
      if (error !== 'cancel') {
        console.error('删除代码失败:', error)
        ElMessage.error('删除代码失败')
      }
    } finally {
      deletingId.value = null
    }
  }

  // 保存代码
  const handleSaveCode = async () => {
    try {
      saving.value = true

      const saveData: CreateSavedCodeRequest = {
        language: props.language,
        code: code.value,
        remark: saveForm.value.remark || undefined,
        documentID: props.documentID,
      }

      const response = await reqCreateSavedCode(saveData)

      if (response.code === 200) {
        ElMessage.success('代码保存成功')
        handleSaveDialogClose()
        // 重新加载列表
        await loadSavedCodeList()
      } else {
        ElMessage.error(response.message || '保存失败')
      }
    } catch (error) {
      console.error('保存代码失败:', error)
      ElMessage.error('保存代码失败')
    } finally {
      saving.value = false
    }
  }

  // 监听语言变化
  watch(
    () => props.language,
    () => {
      if (visible.value) {
        loadSavedCodeList()
      }
      // 更新保存表单的语言
      saveForm.value.language = props.language
    },
  )

  // 监听代码变化
  watch(
    () => code.value,
    () => {
      // 更新保存表单的代码
      saveForm.value.code = code.value
    },
  )

  // 监听对话框显示
  watch(visible, (newVal) => {
    if (newVal) {
      loadSavedCodeList()
    }
  })

  onMounted(() => {
    if (visible.value) {
      loadSavedCodeList()
    }
  })
</script>

<style scoped lang="scss">
  .saved-code-overlay {
    position: fixed;
    inset: 0;
    z-index: 2000;
    display: flex;
    align-items: center;
    justify-content: center;
    background-color: rgb(0 0 0 / 50%);
  }

  .saved-code-card {
    display: flex;
    flex-direction: column;
    width: 800px;
    max-width: 90vw;
    max-height: 90vh;
    overflow: hidden;
    background-color: $base-background-color;
    border-radius: 8px;
    box-shadow: $box-shadow;
  }

  .card-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 10px 24px;
    background-color: $base-header-background;
    border-bottom: 1px solid $border-color;

    .card-title {
      margin: 0;
      font-size: 16px;
      font-weight: 600;
      color: $text-color-primary;
    }

    .close-btn {
      padding: 4px;
      border-radius: 4px;

      &:hover {
        background-color: rgb(0 0 0 / 10%);
      }
    }
  }

  .card-body {
    flex: 1;
    padding: 24px;
    overflow: hidden;
  }

  .modal-description {
    padding: 12px 16px;
    margin-bottom: 20px;
    font-size: 14px;
    line-height: 1.6;
    color: #666;
    background-color: #e3f2fd;
    border-left: 4px solid #2196f3;
    border-radius: 4px;
  }

  .save-new-code {
    margin-bottom: 12px;
    text-align: right;

    .el-button {
      padding: 6px 12px;
      font-size: 13px;
    }
  }

  .code-list {
    .no-remark {
      font-size: 13px;
      font-style: italic;
      color: #999;
    }
  }

  .save-dialog-content {
    .el-form-item {
      margin-bottom: 20px;
    }
  }

  .dialog-footer {
    padding-top: 16px;
    margin-top: 20px;
    text-align: right;
    border-top: 1px solid #e9ecef;

    .el-button {
      margin-left: 12px;

      &:first-child {
        margin-left: 0;
      }
    }
  }

  :deep(.el-table) {
    overflow: hidden;
    border-radius: 6px;

    .el-table__header {
      background-color: $base-header-background;

      th {
        font-weight: 600;
        color: $text-color-primary;
        background-color: $base-header-background;
        border-bottom: 1px solid $border-color;
      }
    }

    .el-table__row {
      cursor: default;
      transition: background-color 0.2s ease;

      &:hover {
        background-color: $base-header-background;
      }

      td {
        padding: 12px 0;
        border-bottom: 1px solid #f0f0f0;
      }
    }

    .el-radio {
      cursor: pointer;
    }

    .apply-btn {
      height: 28px;
      padding: 4px 12px;
      margin-right: 8px;
      font-size: 12px;
    }

    .delete-btn {
      height: 28px;
      padding: 4px 12px;
      font-size: 12px;
    }
  }

  :deep(.el-dialog) {
    overflow: hidden;
    border-radius: 8px;

    .el-dialog__header {
      padding: 20px 24px;
      background-color: $base-header-background;
      border-bottom: 1px solid $border-color;

      .el-dialog__title {
        font-size: 16px;
        font-weight: 600;
        color: $text-color-primary;
      }
    }

    .el-dialog__body {
      padding: 24px;
    }

    .el-dialog__footer {
      padding: 16px 24px;
      background-color: $base-header-background;
      border-top: 1px solid #e9ecef;
    }
  }

  :deep(.delete-confirm-dialog) {
    .el-message-box__header {
      padding: 20px 24px;
      background-color: $base-header-background;
      border-bottom: 1px solid $border-color;

      .el-message-box__title {
        font-size: 16px;
        font-weight: 600;
        color: $text-color-primary;
      }
    }

    .el-message-box__content {
      padding: 24px;
      line-height: 1.6;
      color: #666;
    }

    .el-message-box__btns {
      padding: 16px 24px;
      background-color: $base-header-background;
      border-top: 1px solid #e9ecef;

      .el-button {
        margin-left: 12px;

        &:first-child {
          margin-left: 0;
        }
      }
    }
  }
</style>
