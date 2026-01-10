<template>
  <el-form status-icon label-width="auto" class="arrayDescription" size="small">
    <el-form-item label="图节点">
      <el-input v-model="graphDescription.graphNode" />
    </el-form-item>
    <el-form-item label="节点值">
      <el-input v-model="graphDescription.data" />
    </el-form-item>
    <el-form-item label="next指针">
      <el-input v-model="pointNamesStr" />
    </el-form-item>
  </el-form>
</template>

<script setup lang="ts">
  import useVisualStore from '@/store/modules/visual.ts'
  import { descriptions } from '@/constants/description.ts'
  import { onMounted, ref, watch, computed } from 'vue'
  import { storeToRefs } from 'pinia'

  const visualStore = useVisualStore()
  const { descriptions: storeDescriptions } = storeToRefs(visualStore)

  // 使用计算属性获取图描述
  const graphDescription = computed({
    get: () => storeDescriptions.value[descriptions.Graph],
    set: (value) => {
      visualStore.setDescription(descriptions.Graph, value)
    },
  })

  const pointNamesStr = ref('')

  watch(
    () => pointNamesStr.value,
    (val) => {
      let arr: string[] = val.split(',').map((item) => {
        return item.trim()
      })
      graphDescription.value.nexts = arr
    },
  )

  watch(
    () => graphDescription.value.nexts,
    (val) => {
      let arr: string[] = pointNamesStr.value.split(',').map((item) => {
        return item.trim()
      })
      if (JSON.stringify(val) != JSON.stringify(arr)) {
        pointNamesStr.value = val.join(',')
      }
    },
    {
      deep: true,
    },
  )

  onMounted(() => {
    pointNamesStr.value = graphDescription.value.nexts.join(',')
  })
</script>

<style lang="scss" scoped>
  .arrayDescription {
    box-sizing: border-box;
    width: 100%;
    height: 100%;
    padding: 20px;
  }
</style>
