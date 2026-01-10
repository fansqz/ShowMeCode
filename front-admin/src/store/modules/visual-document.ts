import { defineStore } from 'pinia';

export type VisualDocument = {
  bankID: number;
  id: number;
  parentID: number;
  title: string;
  content: string;
  enable: boolean;
  codeList: VisualDocumentCode[];
};

const useVisualDocumentStore = defineStore('visual-document', {
  state: (): VisualDocument => ({
    bankID: 0,
    id: 0,
    parentID: 0,
    title: '',
    content: '',
    enable: false,
    codeList: [],
  }),
  persist: true,
});

export type VisualDocumentCode = {
  code: string;
  language: string;
  breakpoints: number[];
};

export default useVisualDocumentStore;
