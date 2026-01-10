import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import path from 'path'
import prismjs from 'vite-plugin-prismjs'
import monacoEditorPlugin from 'vite-plugin-monaco-editor'

// 引入svg需要用的插件
import { createSvgIconsPlugin } from 'vite-plugin-svg-icons'

export default defineConfig(() => {
  return {
    plugins: [
      vue(),
      monacoEditorPlugin({}),
      // svg插件
      createSvgIconsPlugin({
        iconDirs: [path.resolve(process.cwd(), 'src/assets/icons')],
        symbolId: 'icon-[dir]-[name]',
      }),
      prismjs({
        languages: ['json', 'go', 'c', 'java', 'cpp'],
      }),
    ],
    resolve: {
      alias: {
        '@': path.resolve('./src'),
      },
    },
    css: {
      preprocessorOptions: {
        scss: {
          additionalData: `@import "./src/styles/variables.scss";`,
        },
      },
    },
  }
})
