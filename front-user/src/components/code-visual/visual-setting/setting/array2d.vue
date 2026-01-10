<template>
  <el-form status-icon label-width="auto" class="array2DDescription" size="small">
    <el-form-item label="数组名">
      <el-input v-model="array2DDescription.arrayName" />
    </el-form-item>
    <el-form-item label="取值变量-行">
      <el-input v-model="rowPointNamesStr" />
    </el-form-item>
    <el-form-item label="取值变量-列">
      <el-input v-model="colPointNamesStr" />
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

  // 使用计算属性获取二维数组描述
  const array2DDescription = computed({
    get: () => storeDescriptions.value[descriptions.Array2D],
    set: (value) => {
      visualStore.setDescription(descriptions.Array2D, value)
    },
  })

  const rowPointNamesStr = ref('')
  const colPointNamesStr = ref('')

  watch(
    () => rowPointNamesStr.value,
    (val) => {
      let arr: string[] = val.split(',').map((item) => {
        return item.trim()
      })
      array2DDescription.value.rowPointNames = arr
    },
  )

  watch(
    () => colPointNamesStr.value,
    (val) => {
      let arr: string[] = val.split(',').map((item) => {
        return item.trim()
      })
      array2DDescription.value.colPointNames = arr
    },
  )

  watch(
    () => array2DDescription.value.rowPointNames,
    (val) => {
      let arr: string[] = rowPointNamesStr.value.split(',').map((item) => {
        return item.trim()
      })
      if (JSON.stringify(val) != JSON.stringify(arr)) {
        rowPointNamesStr.value = val.join(',')
      }
    },
    {
      deep: true,
    },
  )

  watch(
    () => array2DDescription.value.colPointNames,
    (val) => {
      let arr: string[] = colPointNamesStr.value.split(',').map((item) => {
        return item.trim()
      })
      if (JSON.stringify(val) != JSON.stringify(arr)) {
        colPointNamesStr.value = val.join(',')
      }
    },
    {
      deep: true,
    },
  )

  onMounted(() => {
    colPointNamesStr.value = array2DDescription.value.colPointNames.join(',')
    rowPointNamesStr.value = array2DDescription.value.rowPointNames.join(',')
  })
</script>

<style lang="scss" scoped>
  .array2DDescription {
    box-sizing: border-box;
    width: 100%;
    height: 100%;
    padding: 20px;
  }
</style>
