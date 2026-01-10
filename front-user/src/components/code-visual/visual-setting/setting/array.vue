<template>
  <el-form status-icon label-width="auto" class="arrayDescription" size="small">
    <el-form-item label="数组名">
      <el-input v-model="arrayDescription.arrayName" />
    </el-form-item>
    <el-form-item label="取值变量">
      <el-input v-model="pointNamesStr" />
    </el-form-item>
    <el-form-item label="展示方式">
      <el-radio-group v-model="arrayDescription.displayType">
        <el-radio label="array">正常展示</el-radio>
        <el-radio label="array-bar">柱状图</el-radio>
      </el-radio-group>
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

  // 使用计算属性获取数组描述
  const arrayDescription = computed({
    get: () => storeDescriptions.value[descriptions.Array],
    set: (value) => {
      visualStore.setDescription(descriptions.Array, value)
    },
  })

  const pointNamesStr = ref('')

  watch(
    () => pointNamesStr.value,
    (val) => {
      let arr: string[] = val.split(',').map((item) => {
        return item.trim()
      })
      arrayDescription.value.pointNames = arr
    },
  )

  watch(
    () => arrayDescription.value.pointNames,
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
    pointNamesStr.value = arrayDescription.value.pointNames.join(',')
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
