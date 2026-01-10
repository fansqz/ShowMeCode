<template>
  <el-container class="visual-description">
    <el-main class="main">
      <div class="template-select">
        <div
          v-for="item in descriptionTypeList"
          :key="item.type"
          class="template-item"
          :class="{ active: item.type === descriptionType }"
          @click="handleCurrentChange(item)"
        >
          <i v-if="item.type == descriptionType" class="selected-icon iconfont icon-Other-10" />
          <i v-if="item.type != descriptionType" class="selected-icon iconfont icon-Other-11" />
          {{ item.name }}
        </div>
      </div>
      <div class="setting">
        <ArraySetting v-if="descriptionType == descriptions.Array" />
        <Array2DSetting v-if="descriptionType == descriptions.Array2D" />
        <BinaryTreeSetting v-if="descriptionType == descriptions.BinaryTree" />
        <LinkList v-if="descriptionType == descriptions.LinkList" />
        <Graph v-if="descriptionType == descriptions.Graph" />
      </div>
    </el-main>
  </el-container>
</template>

<script setup lang="ts">
  import useVisualStore from '@/store/modules/visual.ts'
  import useDebugStore from '@/store/modules/debug'
  import { reqGetVisualDescription } from '@/api/visual'
  import { onMounted, ref, watch } from 'vue'
  import { storeToRefs } from 'pinia'
  import { descriptions } from '@/constants/description.ts'
  import ArraySetting from './array.vue'
  import Array2DSetting from './array2d.vue'
  import BinaryTreeSetting from './binary-tree.vue'
  import LinkList from './link-list.vue'
  import Graph from './graph.vue'
  import { CompileEventDispatcher } from '@/api/debug/debug-event-dispatcher'
  import { CompileEvent } from '@/api/debug/event'

  const visualStore = useVisualStore()
  const debugStore = useDebugStore()

  const { isAIEnabled, descriptionType } = storeToRefs(visualStore)

  const descriptionTypeList = ref([
    {
      name: '一维数组',
      type: descriptions.Array,
    },
    {
      name: '二维数组',
      type: descriptions.Array2D,
    },
    {
      name: '二叉树',
      type: descriptions.BinaryTree,
    },
    {
      name: '链表',
      type: descriptions.LinkList,
    },
    {
      name: '图',
      type: descriptions.Graph,
    },
  ])

  // 选择了其他可视化模板
  const handleCurrentChange = async (val) => {
    descriptionType.value = val.type
  }

  // 编译完成以后读取可视化描述
  const onCompiled = async (data: CompileEvent) => {
    if (data.success && isAIEnabled.value) {
      // 拉取ai配置数据
      let result = await reqGetVisualDescription(debugStore.id)
      if (result.code == 200) {
        descriptionType.value = result.data.visualType
        visualStore.setDescription(descriptionType.value, result.data.description)
      }
    }
  }

  // 监控isAIEnabled的变化
  watch(isAIEnabled, async (newVal) => {
    if (newVal) {
      console.log('isAIEnabled', newVal)
      // 如果isAIEnabled为true，并且正在调试中，则拉取ai配置数据
      if (debugStore.isDebugging()) {
        // 拉取ai配置数据
        let result = await reqGetVisualDescription(debugStore.id)
        if (result.code == 200) {
          descriptionType.value = result.data.visualType
          visualStore.setDescription(descriptionType.value, result.data.description)
        }
      }
    }
  })

  // 监控调试状态
  onMounted(() => {
    // 注册一些事件
    CompileEventDispatcher.on('compile', onCompiled)
  })
</script>

<style lang="scss" scoped>
  .visual-description {
    background-color: $base-background-color;

    .title {
      height: 30px;
      border-bottom: 1px solid $border-color;
    }

    .main {
      box-sizing: border-box;
      display: flex;
      width: 100%;
      height: 100%;
      padding: 0;

      .template-select {
        width: 20%;
        height: 150px;
        margin-right: 10px;
        margin-left: 20px;
        overflow: hidden auto;
        background-color: $base-background-color;
        border-right: 1px solid $border-color;

        // 自定义滚动条样式
        &::-webkit-scrollbar {
          width: 6px;
        }

        &::-webkit-scrollbar-track {
          background: $fill-color-extra-light;
          border-radius: 3px;
        }

        &::-webkit-scrollbar-thumb {
          background: $border-color;
          border-radius: 3px;

          &:hover {
            background: $border-color-light;
          }
        }

        .template-item {
          display: flex;
          flex-shrink: 0;
          align-items: center;
          height: 28px;
          padding: 6px 12px;
          color: $text-color-primary;
          cursor: pointer;
          background-color: $base-background-color;
          transition: all 0.2s ease;

          &:hover {
            background-color: $fill-color-light;
          }

          &.active {
            font-weight: 500;
            color: $primary-color;
            background-color: $fill-color;
          }

          .selected-icon {
            margin-right: 8px;
            font-size: 14px;
            color: $primary-color;
          }
        }
      }

      .setting {
        width: calc(80% - 100px);
        height: 150px;
      }

      .option {
        display: flex;
        flex-flow: column;
        width: 40px;
        height: 150px;
        margin-right: 20px;
        margin-left: 10px;
      }
    }
  }
</style>
