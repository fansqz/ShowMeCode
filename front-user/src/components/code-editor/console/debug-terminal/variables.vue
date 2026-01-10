<template>
  <div class="variables">
    <el-tree
      lazy
      :load="loadNode"
      :data="treeData"
      :props="defaultProps"
      node-key="id"
      :ref="treeRef"
      class="variables-tree"
    />
  </div>
</template>

<script setup lang="ts">
  import { ref, watch, onMounted, onUnmounted } from 'vue'
  import useDebugStore from '@/store/modules/debug'
  import { storeToRefs } from 'pinia'
  import { reqGetFrameVariables, reqGetVariables } from '@/api/debug'
  import { TerminatedEventDispatcher } from '@/api/debug/debug-event-dispatcher'

  const debugStore = useDebugStore()
  const { id, status, currentFrameID } = storeToRefs(debugStore)
  const treeData = ref([]) // 树的数据源
  const props = defineProps(['frameId'])
  const treeRef = ref()

  const defaultProps = {
    children: 'children',
    label: 'name',
    isLeaf: 'isLeaf',
  }

  // 监控栈帧变化，主动触发load
  const loadNode = (node: any, resolve: any) => {
    if (props.frameId != '' && debugStore.isDebugging()) {
      if (node.level == 0) {
        getVariablesByFrameId(props['frameId'], resolve)
      } else {
        getVariables(node.data.reference, resolve)
      }
    }
  }

  const getVariablesByFrameId = async (frameId: number, resolve: any) => {
    let result = await reqGetFrameVariables(id.value, frameId)
    if (result.code == 200) {
      let variables = result.data
      let tree = []
      for (let i in variables) {
        let treeNode = {
          name: `${variables[i].name}(${variables[i].type})`,
          reference: variables[i].reference,
          children: [],
          isLeaf: variables[i].reference == 0,
        }
        if (variables[i].value != '') {
          treeNode.name = `${treeNode.name}:${variables[i].value}`
        }
        tree.push(treeNode)
      }
      resolve(tree)
    }
  }

  const getVariables = async (reference: number, resolve: any) => {
    let result = await reqGetVariables(id.value, reference)
    if (result.code == 200) {
      let variables = result.data
      let tree = []
      for (let i in variables) {
        let treeNode = {
          name: `${variables[i].name}(${variables[i].type})`,
          reference: variables[i].reference,
          children: [],
          isLeaf: variables[i].reference == 0,
        }
        if (variables[i].value != '') {
          treeNode.name = `${treeNode.name}:${variables[i].value}`
        }
        tree.push(treeNode)
      }
      resolve(tree)
    }
  }

  const onExited = () => {
    // 程序退出所有信息
    treeData.value = []
  }

  // 主动更新
  const updateVariables = (frameId: number) => {
    getVariablesByFrameId(frameId, (tree: any) => {
      treeData.value = tree
    })
  }

  defineExpose({
    updateVariables,
  })

  onMounted(() => {
    // 监控栈帧变化，如果变化则重新加载数据
    watch(
      () => currentFrameID.value,
      () => {
        if (currentFrameID.value == -1) {
          treeData.value = []
          return
        }
        if (status.value == 'stopped') {
          updateVariables(currentFrameID.value)
        }
      },
    )
  })
  onUnmounted(() => {
    TerminatedEventDispatcher.off('terminated', onExited)
  })
</script>

<style lang="scss" scoped>
  .variables-tree {
    width: 100%;
    height: 100%;
    background-color: $base-background-color;

    :deep(.el-tree) {
      background-color: $base-background-color;
    }

    :deep(.el-tree-node) {
      background-color: $base-background-color;
    }

    :deep(.el-tree-node__content) {
      background-color: $base-background-color;

      &:hover {
        background-color: rgba($primary-color, 0.1);
      }
    }
  }
</style>
