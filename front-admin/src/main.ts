import { createApp } from 'vue';
import App from '@/App.vue';
// 引入element-plus插件
import ElementPlus from 'element-plus';
import 'element-plus/dist/index.css';
// 配置element-plus国际化
// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore
import zhCn from 'element-plus/dist/locale/zh-cn.mjs';
// svg配置
import 'virtual:svg-icons-register';
// 全局组件
import gloalComponent from '@/components';
// 路由
import router from '@/router';
// 仓库
import pinia from '@/store';
// 引入路由鉴权文件
import './premisstion';
// 全局样式
import '@/styles/index.scss';

// markdowm
import VueMarkdownEditor from '@kangc/v-md-editor';
import '@kangc/v-md-editor/lib/style/base-editor.css';
import '@kangc/v-md-editor/lib/style/preview-html.css';
import vuepressTheme from '@kangc/v-md-editor/lib/theme/vuepress.js';
import '@kangc/v-md-editor/lib/theme/style/vuepress.css';
import Prism from 'prismjs';

const app = createApp(App);

// 初始化markdown编辑器
const initMarkdowmEditor = (app: any) => {
  VueMarkdownEditor.use(vuepressTheme, {
    Prism,
  });
  VueMarkdownEditor.xss.extend({
    // 扩展白名单
    whiteList: {
      iframe: ['src', 'width', 'height', 'frameborder', 'allowfullscreen'],
      source: [],
    },
  });
  app.use(VueMarkdownEditor);
};

app.use(ElementPlus, {
  //element-plus国际化
  locale: zhCn,
});
app.use(gloalComponent as any);
app.use(router);
app.use(pinia as any);
app.mount('#app');
initMarkdowmEditor(app);
