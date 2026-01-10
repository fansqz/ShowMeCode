import { WatchStopHandle, onBeforeUnmount, onMounted, watch } from 'vue'

import debounce from 'lodash.debounce'
import { editor } from 'monaco-editor'

import useCodingStore from '@/store/modules/coding'
import useDebugStore from '@/store/modules/debug'
import useThemeStore from '@/store/modules/theme'

import { registerAllProvider, unregisterAllProvider } from '../provider/index.ts'
import { initWorker, initTheme, changeTheme, wire } from '../themes'
import { EditorInstance, VsCode } from '../types'
import { setBreakPoint, initBP, getAllBreakpoint } from '../utils/breakpoint'
import { debug, cancelHighlightLine, highlightLine, scrollIntoView } from '../utils/debugUtils'
import editUtils from '../utils/editUtils'

// TODO: 调试时修改代码改变高亮行，输入改变时更新断点
/**
 *
 * @param { VsCode } params
 */
export const useVsCode = (vscode: VsCode) => {
  const debugStore = useDebugStore()
  const codingStore = useCodingStore()
  const themeStore = useThemeStore()
  initWorker()
  return new Promise((resolve) => {
    const { target, onContentChanged, onEditorBlur, onCtrlS, onUpdateBP } = vscode

    let stopValueWatch: WatchStopHandle
    let stopLanguageWatch: WatchStopHandle
    let stopThemeWatch: WatchStopHandle
    let stopLineNumWatch: WatchStopHandle
    let stopErrorLocationWatch: WatchStopHandle
    let stopBreakpointsWatch: WatchStopHandle
    let editorInstance: EditorInstance
    let isUnmounted = false

    onMounted(async () => {
      if (!target.value) {
        throw new TypeError('target值必须是一个正确的vue ref')
      }
      // 必须先初始化主题再创建实例
      await initTheme(themeStore.editorTheme)

      const { getConfigs } = await import('../conf')

      // 注册代码提示
      registerAllProvider()

      // 创建editor实例
      editorInstance = editor.create(
        target.value,
        getConfigs({
          value: codingStore.code,
          readOnly: false,
          theme: themeStore.editorTheme,
          language: codingStore.language,
        }),
      )

      // 设置语言
      wire(codingStore.language, editorInstance)

      // 初始化断点
      initBP(editorInstance, debugStore.breakpoints)

      // 监控value的变化
      stopValueWatch = watch(
        () => codingStore.code,
        () => {
          if (codingStore.code != editorInstance.getValue()) {
            editorInstance?.setValue(codingStore.code)
          }
        },
      )

      // 监控语言变化
      stopLanguageWatch = watch(
        () => codingStore.language,
        async (val) => {
          const model = editorInstance.getModel()
          if (model) {
            //设置语言
            editor.setModelLanguage(model, val)
            wire(val, editorInstance)
          }
        },
      )

      // 监控主题变化
      stopThemeWatch = watch(
        () => themeStore.editorTheme,
        async (val) => {
          changeTheme(val, editorInstance)
        },
      )

      // 调试
      stopLineNumWatch = watch(
        () => debugStore.lineNum,
        async (val, oldVal) => {
          const model = editorInstance.getModel()
          debug({ model, currentline: val, preline: oldVal, instance: editorInstance })
        },
        { immediate: true },
      )

      // 错误定位
      stopErrorLocationWatch = watch(
        () => debugStore.currentErrorLocation,
        async (val: number[]) => {
          // 设置高亮行
          const currentLine = Number(val[0])
          const model = editorInstance.getModel()
          // 如果模型有效，那么高亮错误行
          if (model) {
            highlightLine(model, currentLine)
            // 滚动编辑器视图确保错误行在视野中
            scrollIntoView(editorInstance, currentLine)
            // 添加鼠标点击事件，用于在用户点击时取消高亮显示
            editorInstance.onMouseDown(function () {
              // 获取当前行的所有装饰（高亮）
              const decorations = model?.getLineDecorations(currentLine)
              // 获取包含'debug-line'类的装饰的ID数组
              const decorationsIds = decorations
                .filter((item) => item.options.className?.includes('debug-line'))
                .map((item) => item.id)
              // 调用取消高亮行函数，移除之前的高亮显示
              cancelHighlightLine(model, decorationsIds)
            })
          }
        },
        // 深度监听
        { deep: true },
      )

      // 监控断点变化
      stopBreakpointsWatch = watch(
        () => debugStore.breakpoints,
        () => {
          if (!editorInstance) {
            return
          }
          const bps = getAllBreakpoint(editorInstance)
          let bps2 = debugStore.breakpoints
          if (bps2 && bps2.length > 0) {
            bps2 = bps2.slice().sort()
            if (!isEqual(bps, bps2)) {
              initBP(editorInstance, debugStore.breakpoints)
            }
          }
        },
      )

      function isEqual(arr1: any[], arr2: any[]) {
        if (!arr1 || !arr2) {
          return arr1 == arr2
        }
        if (arr1.length !== arr2.length) {
          return false
        }
        return !arr1.some((item) => !arr2.includes(item))
      }

      // 监控断点情况
      setBreakPoint(editorInstance, onUpdateBP)

      // 成功回调
      resolve(editorInstance)

      // 监控内容修改，执行回调
      editorInstance.onDidChangeModelContent(
        debounce(async () => {
          if (isUnmounted) return
          // 使用防抖，在不输入内容的时候进行保存
          onContentChanged?.(editorInstance.getValue())
          // 更新断点
          const bps = getAllBreakpoint(editorInstance)
          const bps2 = debugStore.breakpoints
          if (bps2 && bps2.length > 0) {
            for (const breakpoint of bps2) {
              if (!bps || !bps.includes(breakpoint)) {
                vscode.onUpdateBP(bps, breakpoint, 'del')
              }
            }
            if (!isEqual(bps, bps2)) {
              initBP(editorInstance, debugStore.breakpoints)
            }
          }
        }, 600),
      )

      // 监听编辑器文本区域失去焦点的事件
      editorInstance.onDidBlurEditorText(() => {
        onEditorBlur?.(editorInstance.getValue())
      })

      // 添加对ctl + s的监听
      editUtils.addSaveCodeListener(editorInstance, onCtrlS)
    })

    // 在组件即将卸载前执行，做些销毁动作
    onBeforeUnmount(() => {
      isUnmounted = true
      stopValueWatch && stopValueWatch()
      stopLanguageWatch && stopLanguageWatch()
      stopThemeWatch && stopThemeWatch()
      stopLineNumWatch && stopLineNumWatch()
      stopErrorLocationWatch && stopErrorLocationWatch()
      stopBreakpointsWatch && stopBreakpointsWatch()
      const model = editorInstance.getModel()
      model?.dispose()
      editorInstance?.dispose()
      // 解注册代码提示
      unregisterAllProvider()
      console.log('离开咯', editorInstance)
    })
  })
}
