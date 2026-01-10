import * as monaco from 'monaco-editor'
import { wireTmGrammars } from 'monaco-editor-textmate'
import { Registry } from 'monaco-textmate'

import { EditorInstance } from '../types'

import { scopeNameMap, tmGrammarJsonMap, codeThemeList } from './config'

let hasGetAllWorkUrl = false

// 主题缓存，避免重复加载
const themeCache = new Map<string, any>()

export const initWorker = async () => {
  const originalGetWorkerUrl = window.MonacoEnvironment.getWorkerUrl
  window.MonacoEnvironment.getWorkerUrl = new Proxy(originalGetWorkerUrl, {
    apply(target, thisArg, args) {
      hasGetAllWorkUrl = true
      return target.apply(thisArg, args)
    },
  })
}

/**
 * 按需加载单个主题
 * @param theme 主题名称
 * @returns Promise<boolean> 加载是否成功
 */
export const loadTheme = async (theme: string): Promise<boolean> => {
  // 如果主题已经加载过，直接返回成功
  if (themeCache.has(theme)) {
    console.log(`Theme ${theme} already loaded`)
    return true
  }

  try {
    console.log(`Loading theme: ${theme}`)
    const themeData = await (await fetch(`/themes/${theme}.json`)).json()
    monaco.editor.defineTheme(theme, themeData)
    themeCache.set(theme, themeData)
    console.log(`Theme ${theme} loaded successfully`)
    return true
  } catch (error) {
    console.error(`Failed to load theme: ${theme}`, error)
    return false
  }
}

/**
 * 初始化默认主题（可选，用于预加载常用主题）
 */
export const initTheme = async (editorTheme?: string) => {
  // 确定要加载的主题
  const themeToLoad = editorTheme || 'BlulocoLight' // 使用传入的主题或默认主题

  console.log('Initializing theme:', themeToLoad)
  await loadTheme(themeToLoad)
}

/**
 * https://www.cnblogs.com/wanglinmantan/p/15345204.html
 * https://github.com/brijeshb42/monaco-themes
 */
export const wire = async (languageId, editor) => {
  if (!scopeNameMap[languageId]) {
    return
  }
  // 语言id到作用域名称的映射
  const grammars = new Map()
  grammars.set(languageId, scopeNameMap[languageId])

  // 创建一个注册表，可以从作用域名称创建语法
  const registry = new Registry({
    getGrammarDefinition: async (scopeName: string, _dependentScope: string) => {
      const jsonMap = tmGrammarJsonMap[scopeName]

      try {
        const content = await (await fetch(`/grammars/${jsonMap}`)).text()
        return {
          format: 'json',
          content,
        }
      } catch (error) {
        console.error(`Failed to load grammar: ${jsonMap}`, error)
        return {
          format: 'json',
          content: '',
        }
      }
    },
  })
  // 循环检测
  const loop = () => {
    if (hasGetAllWorkUrl) {
      Promise.resolve().then(async () => {
        try {
          await wireTmGrammars(monaco, registry, grammars, editor)
        } catch (error) {
          console.error('Failed to wire grammars:', error)
        }
      })
    } else {
      setTimeout(() => {
        loop()
      }, 100)
    }
  }
  loop()
}

/**
 * @param theme 主题名称
 * @param editorInstance 编辑器实例
 */
export const changeTheme = async (theme: string, editorInstance: EditorInstance) => {
  // 先尝试加载主题
  const success = await loadTheme(theme)
  if (success) {
    editorInstance.updateOptions({ theme })
  } else {
    console.warn(`Failed to load theme: ${theme}, keeping current theme`)
  }
}

/**
 * 获取所有的主题列表
 * @returns 所有的主题列表
 */
export const getAllTheme = (): string[] => {
  return codeThemeList
}
