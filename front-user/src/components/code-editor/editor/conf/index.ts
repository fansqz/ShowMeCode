import { editor } from 'monaco-editor'

type Configs = editor.IStandaloneEditorConstructionOptions

export const getConfigs = (configs: Configs) =>
  <Configs>{
    // 自动布局，编辑器会根据可用空间自动调整布局
    automaticLayout: true,

    // 参数提示，启用后编辑器会显示参数文档和类型信息
    parameterHints: { enabled: true },

    // 使用Shadow DOM进行封装，用于隔离编辑器的样式
    useShadowDOM: true,

    // 显示弃用警告，编辑器会标识出不再推荐使用的代码
    showDeprecated: true,

    // 基于文本的建议将仅限于同一语言模式
    wordBasedSuggestionsOnlySameLanguage: true,

    // 括号配对着色，启用后编辑器会为匹配的括号对提供不同的颜色
    // 为每种括号类型提供独立的颜色池
    bracketPairColorization: {
      enabled: true,
      independentColorPoolPerBracketType: true,
    },

    // 粘贴格式化，粘贴内容时自动格式化代码
    formatOnPaste: true,

    // 输入时格式化，每次键入时自动格式化代码
    formatOnType: true,

    // 字形边缘，显示在编辑器左侧的装饰性边缘，常用于断点等
    glyphMargin: true,

    // 行号最小字符数，设置编辑器行号的最小宽度
    lineNumbersMinChars: 3,

    // 启用代码操作小灯泡提示，为快速修复提供图标
    lightbulb: { enabled: true },

    // 括号引导线，显示与括号配对相关的垂直线
    guides: { bracketPairs: true },

    // 制表符补全，启用后编辑器支持使用Tab键自动补全代码
    tabCompletion: 'on',

    // 智能提示相关的配置
    suggest: {
      // 输入时可预览建议的效果
      preview: true,

      // 预览模式为前缀
      previewMode: 'prefix',

      // 共享建议选择，在开启多个编辑器时共享选择状态
      shareSuggestSelections: true,

      // 显示函数建议
      showFunctions: true,

      // 片段阻止快速建议，片段激活时是否禁用快速建议
      snippetsPreventQuickSuggestions: false,
    },

    minimap: {
      enabled: false,
    },

    // 扩展传入的配置对象，覆盖默认设置
    ...configs,
  }
