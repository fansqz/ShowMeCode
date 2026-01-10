<template>
  <div class="directory-container">
    <div class="bank-name">
      <el-text class="bank-name">{{ bank.name }}</el-text>
    </div>
    <el-scrollbar class="directory-tree">
      <el-tree
        :data="data"
        :props="props"
        :highlight-current="true"
        :current-node-key="currentId"
        :expand-on-click-node="false"
        node-key="id"
        :indent="15"
        @node-click="handleNodeClick"
      >
        <template #default="scope">
          <div class="tree_item">
            <el-text>{{ scope.node.label }}</el-text>
          </div>
        </template>
      </el-tree>
    </el-scrollbar>
  </div>
</template>

<script setup lang="ts">
  import { ref, onMounted, watch, computed } from 'vue'
  import { reqVisualDocumentDirectory } from '@/api/visual-document/index.ts'
  import { VisualDocumentDirectory } from '@/api/visual-document/type.ts'
  import { reqVisualDocumentBank } from '@/api/visual-document-bank/index.ts'
  import { useRouter, useRoute } from 'vue-router'
  import useVisualDocument from '@/store/modules/visual-document'
  import { storeToRefs } from 'pinia'

  const $router = useRouter()
  const $route = useRoute()
  const visualDocumentStore = useVisualDocument()
  const { id } = storeToRefs(visualDocumentStore)

  const data = ref<VisualDocumentDirectory[]>([])
  const props = {
    children: 'children',
    label: 'title',
  }
  let bankID: number
  let bank = ref({
    name: '',
  })

  const currentId = computed(() => id.value)

  const convertToNumber = (input: string | string[]): number => {
    if (input == undefined) {
      return 0
    }
    let value: string
    if (Array.isArray(input)) {
      value = input[0]
    } else {
      value = input
    }
    return Number(value)
  }

  const getDirectory = async (bankID: number) => {
    let result = await reqVisualDocumentDirectory(bankID)
    if (result.code === 200) {
      data.value = result.data
    }
  }

  const getBank = async (bankID: number) => {
    let result = await reqVisualDocumentBank(bankID)
    if (result.code === 200) {
      bank.value = result.data
    }
  }

  const handleNodeClick = (node: any) => {
    $router.push({
      name: 'learn',
      params: {
        bankID: bankID,
      },
      query: {
        documentID: node.id,
      },
    })
  }

  const loadDirectory = async () => {
    bankID = convertToNumber($route.params.bankID)
    let documentID = convertToNumber($route.query.documentID)
    getBank(bankID)
    await getDirectory(bankID)
    if (documentID != 0) {
      id.value = documentID
    } else if (data.value.length > 0) {
      id.value = data.value[0].id
    }
  }

  onMounted(async () => {
    loadDirectory()
    watch([() => $route.params, () => $route.query], () => {
      loadDirectory()
    })
  })
</script>

<style scoped lang="scss">
  .directory-container {
    position: relative;
    width: 100%;
    height: 100%;
    background-color: $base-background-color;
    border-right: 1px solid $border-color;

    .bank-name {
      position: relative;
      box-sizing: border-box;
      display: flex;
      align-items: center;
      width: 100%;
      height: 35px;
      padding-left: 10px;

      .el-text {
        font-size: large;
        font-weight: bold;
      }
    }

    .directory-tree {
      position: relative;
      width: 100%;
      height: calc(100% - 35px);

      ::v-deep .el-tree-node__content {
        height: 35px;

        .el-text {
          font-size: medium;
          font-weight: bold;
        }
      }
    }
  }
</style>
