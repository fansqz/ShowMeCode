<template>
  <div class="container">
    <!-- 自定义菜单 -->
    <div class="custom-menu">
      <el-scrollbar class="scrollbar">
        <div
          v-for="(item, index) in stackFrames"
          :key="index"
          :class="{ 'menu-item': true, active: index === parseInt(defaultActive) }"
          @click="handleSelect(index.toString())"
        >
          <el-text class="menu-item-text">{{ item.name }}</el-text>
        </div>
      </el-scrollbar>
    </div>
  </div>
</template>

<script setup lang="ts">
  import { onMounted, ref, watch } from 'vue'
  import { storeToRefs } from 'pinia'
  import { reqGetStackTrace } from '@/api/debug'
  import useDebugStore from '@/store/modules/debug'

  const debugStore = useDebugStore()
  let { id, status, currentFrameID, lineNum } = storeToRefs(debugStore)
  const emit = defineEmits(['selectFrame'])
  const defaultActive = ref('0')
  // 栈帧
  const stackFrames = ref<any[]>([])

  const handleSelect = (key: string) => {
    // 点击某个栈帧，则回调selectFrame
    defaultActive.value = key
    let frame = stackFrames.value[parseInt(key)]
    currentFrameID.value = frame.id
    lineNum.value = frame.line
  }

  onMounted(() => {
    watch(
      () => status.value,
      async (val) => {
        if (val == 'stopped') {
          currentFrameID.value = -1
          // 程序暂停时需要更新栈帧
          let result = await reqGetStackTrace(id.value)
          if (result.code == 200) {
            stackFrames.value = result.data
            handleSelect('0')
          }
        }
        if (val == 'terminated' || val == 'init' || val == 'compiled') {
          stackFrames.value = []
          // 设置为-1，清空变量列表
          currentFrameID.value = -1
        }
      },
    )
  })
</script>

<style lang="scss" scoped>
  .container {
    position: relative;
    display: flex;
    flex-flow: column;
    width: 100%;
    height: 100%;
    padding: 0;
    margin: 0;

    .custom-menu {
      width: 100%;
      height: 100%;

      .scrollbar {
        width: 100%;
        height: 100%;

        .menu-item {
          display: flex;
          align-items: center;
          height: 25px;
          padding-left: 15px;
          cursor: pointer;

          &:hover {
            background-color: rgba($primary-color, 0.1);
          }

          &.active {
            font-weight: 500;
            color: $primary-color;
            background-color: rgba($primary-color, 0.2);
          }

          .menu-item-text {
            overflow: hidden;
            text-overflow: ellipsis;
            white-space: nowrap;
          }
        }
      }
    }

    .no-data-show {
      display: flex;
      align-items: center;
      justify-content: center;
      width: 100%;
      height: 100%;
      text-align: center;

      .text {
        margin: auto;
      }
    }
  }
</style>
