import { defineStore } from 'pinia'

import { supportedLanguages, languageConstants } from '@/constants/languages.ts'

type CodingState = {
  // 编辑器代码
  code: string
  // 判断当前标签页编程题目的语言类型
  language: languageConstants
  // 可选的语言
  languages: languageConstants[]
  // 输入框数据
  input: string
}

// 使用 Pinia 创建一个状态存储
const useCodingStore = defineStore('coding', {
  state: (): CodingState => ({
    code: '',
    language: languageConstants.GO,
    languages: supportedLanguages,
    input: '',
  }),
  actions: {},
  getters: {},
})

export default useCodingStore
