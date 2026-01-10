import { App } from 'vue'

import * as ElementPlusIconsVue from '@element-plus/icons-vue'

import SvgIcon from './svgicon/index.vue'
import TextButton from './text-button/index.vue'

// 引入element-plus提供的全部图标组件

// 对外暴露插件对象
export default {
  install(app: App) {
    app.component('SvgIcon', SvgIcon)
    app.component('TextButton', TextButton)
    // 注册element-plus的图标为全局组件
    for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
      app.component(key, component)
    }
  },
}
