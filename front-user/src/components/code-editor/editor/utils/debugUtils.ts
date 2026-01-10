import { editor, Range } from 'monaco-editor'

import { DebugInfo, EditorInstance, Model } from '../types'

// 高亮指定行
export const highlightLine = (model: Model | null, lineNum: number) => {
  if (!model) {
    return
  }
  model.deltaDecorations(
    [],
    [
      {
        range: new Range(lineNum, 1, lineNum, 1),
        options: {
          isWholeLine: true, // 整行都需要高亮
          className: 'debug-line', // 应用于高亮行的CSS类
          glyphMarginClassName: 'debug-icon iconfont', // 添加到行左边距的glyph图标的CSS类
          stickiness: editor.TrackedRangeStickiness.NeverGrowsWhenTypingAtEdges,
        },
      },
    ],
  )
}

// 取消高亮指定行
export const cancelHighlightLine = (
  model: Model | null,
  decorationsIdentifier: string | string[],
) => {
  if (!model) {
    return
  }
  if (typeof decorationsIdentifier === 'string') {
    model.deltaDecorations([decorationsIdentifier], [])
  } else {
    model.deltaDecorations(decorationsIdentifier, [])
  }
}

/**
 * 处理调试信息，更新编辑器的状态。
 * @param debugInfo 包含有关Monaco编辑器实例和模型的调试信息
 */
export const debug = (debugInfo: DebugInfo) => {
  const { instance, model } = debugInfo
  // 如果model和oldModel都是null，则认为不执行任何操作
  if (!model) {
    return
  }
  const currentLine = Number(debugInfo.currentline ?? 0)
  const previousLine = Number(debugInfo.preline ?? 0)
  // 如果之前有行被高亮，现在取消该行的高亮
  if (model && !model.isDisposed() && previousLine && previousLine > 0) {
    const decorations = model?.getLineDecorations(previousLine)
    const decorationsIds = decorations
      .filter((item) => item.options.className?.includes('debug-line'))
      .map((item) => item.id)
    cancelHighlightLine(model, decorationsIds)
  }
  // 高亮当前行，并且让它在编辑器视图中居中
  if (model && currentLine && currentLine > 0) {
    highlightLine(model, currentLine)
    scrollIntoView(instance, currentLine)
  }
}

// 滚动到编辑器中的指定行，如果该行不在视图中则居中展示
export const scrollIntoView = (instance: EditorInstance, lineNum: number) => {
  instance.revealLineInCenterIfOutsideViewport(lineNum, editor.ScrollType.Smooth)
}
