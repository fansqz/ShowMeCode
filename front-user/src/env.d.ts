/// <reference types="vite/client" />
/// <reference types="vue/dist/vue.d.ts" />

declare module '*.vue' {
  import type { DefineComponent } from 'vue'
  const component: DefineComponent<{}, {}, any>
  export default component
}

declare module 'particles.vue3'

declare module 'splitpanes'

// Pinia 类型声明
declare module 'pinia'

// Element Plus 类型声明
declare module 'element-plus'
declare module 'element-plus/dist/locale/zh-cn.mjs'
declare module '@element-plus/icons-vue'

declare module 'ace-builds/src-noconflict/*'
declare module 'ace-builds/src-noconflict/mode-*'
declare module 'ace-builds/src-noconflict/theme-*'
declare module 'ace-builds/src-noconflict/worker-*'
declare module 'ace-builds/src-noconflict/snippets/*'
declare module 'ace-builds/src-noconflict/ext-*'

// Markdown 编辑器
declare module '@kangc/v-md-editor'
declare module '@kangc/v-md-editor/lib/preview'
declare module '@kangc/v-md-editor/lib/theme/vuepress.js'

// 代码高亮
declare module 'prismjs'

// 图表库
declare module 'echarts'

// 虚拟模块
declare module 'virtual:svg-icons-register'
