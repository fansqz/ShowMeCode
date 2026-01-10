import { defineConfig } from 'vite';
import vue from '@vitejs/plugin-vue';
import path from 'path';
import prismjs from 'vite-plugin-prismjs';

// 引入svg需要用的插件
import { createSvgIconsPlugin } from 'vite-plugin-svg-icons';

export default defineConfig(() => {
  return {
    plugins: [
      vue(),
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
    // scss全局变量配置
    css: {
      preprocessorOptions: {
        scss: {
          javascriptEnables: true,
          additionalData: '@import "@/styles/variable.scss";',
        },
      },
    },
  };
});
