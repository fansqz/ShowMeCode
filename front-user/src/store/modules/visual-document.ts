import { defineStore } from 'pinia'

import { languageConstants } from '@/constants/languages.ts'

type VisualDocumentState = {
  id: number
  content: string
  codeList: VisualDocumentCode[]
}

const useVisualDocumentStore = defineStore('visual-document', {
  state: (): VisualDocumentState => ({
    id: 0,
    content: '',
    codeList: [],
  }),
})

export type VisualDocumentCode = {
  code: string
  language: languageConstants
  breakpoints: number[]
}

export default useVisualDocumentStore
