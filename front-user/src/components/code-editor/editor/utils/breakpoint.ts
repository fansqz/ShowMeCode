import { ElMessageBox } from 'element-plus'
import throttle from 'lodash.throttle'
import { editor, Range } from 'monaco-editor'

import type { onUpdateBP, EditorInstance } from '../types'

import {
  checkBreakPoints,
  existBreakPoint,
  getBreakPointLineNumber,
  getBreakPointOption,
  getBreakPointsIds,
  removeTDecorations,
} from './breakPointUtils'

// 获取所有断点
export function getAllBreakpoint(editorInstance: EditorInstance): number[] {
  if (!editorInstance) {
    return []
  }
  const model = editorInstance.getModel()
  if (!model) {
    return
  }
  const allDecorations = model.getAllDecorations()
  return getBreakPointLineNumber(allDecorations).slice().sort()
}

export function setBreakPoint(editorInstance: EditorInstance, onUpdateBP?: onUpdateBP) {
  let model = editorInstance.getModel()
  removeTDecorations(model)

  editorInstance.onMouseUp((e) => {
    model = editorInstance.getModel()
    if (!model) {
      return
    }
    if (e.target.type === editor.MouseTargetType.GUTTER_GLYPH_MARGIN) {
      // 点击了断点，删除存在透明度的断点
      if (!e.target.position) {
        return
      }
      const { lineNumber } = e.target.position
      const decorations = model.getLineDecorations(lineNumber)
      const existBP = existBreakPoint(decorations)
      if (existBP) {
        // 存在断点，点击取消断点
        const BPIds = getBreakPointsIds(decorations, 'BP')
        model.deltaDecorations(BPIds, [])
        const allDecorations = model.getAllDecorations()
        onUpdateBP?.(getBreakPointLineNumber(allDecorations), lineNumber, 'del')
        return
      }
      const existTBP = existBreakPoint(decorations, true)
      if (existTBP) {
        // 存在拥有透明度的断点，点击设置断点
        const TBPIds = getBreakPointsIds(decorations, 'TBP') // 删除存在透明度的断点
        model.deltaDecorations(TBPIds, [
          {
            range: new Range(lineNumber, 1, lineNumber, 1),
            options: {
              isWholeLine: true,
              // className: "content",
              glyphMarginClassName: 'breakpoint active',
              stickiness: editor.TrackedRangeStickiness.NeverGrowsWhenTypingAtEdges,
              glyphMarginHoverMessage: { value: '断点' },
            },
          },
        ])
        const allDecorations = model.getAllDecorations()
        onUpdateBP?.(getBreakPointLineNumber(allDecorations), lineNumber, 'add')

        const checkBP = checkBreakPoints(model.getLinesContent()[lineNumber - 1]) // lineNumber是从1开始的
        if (!checkBP) {
          // 可能不能执行到的断点给予提示
          ElMessageBox({
            message: `line ${lineNumber}: 在此处设置的断点可能不会被执行到`,
            type: 'warning',
            showClose: true,
          })
        }
      }
    }
  })
  editorInstance.onMouseMove(
    throttle((e: editor.IEditorMouseEvent) => {
      // 离开上述所有区域，删除存在透明度的断点
      model = editorInstance.getModel()
      if (!model) {
        return
      }
      removeTDecorations(model)
      switch (e.target.type) {
        // 移动到行号区域
        case editor.MouseTargetType.GUTTER_LINE_NUMBERS: {
          if (!e.target.position) {
            return
          }
          const { lineNumber } = e.target.position
          const decorations = model.getLineDecorations(lineNumber)
          // {decorations}为什么是数组？因为一行可以有多个Decoration
          const existBP = existBreakPoint(decorations)
          if (!existBP) {
            // 不存在断点就可以显示存在透明度的断点
            model.deltaDecorations([], [getBreakPointOption(lineNumber, 'hide')])
          }
          // 存在断点，无事发生
          break
        }
        // 移动到断点区域
        case editor.MouseTargetType.GUTTER_GLYPH_MARGIN: {
          // 暂时先这样写吧
          if (!e.target.position) {
            return
          }
          const { lineNumber } = e.target.position
          const decorations = model.getLineDecorations(lineNumber)
          const existBP = existBreakPoint(decorations)
          if (!existBP) {
            // 不存在断点就可以显示存在透明度的断点
            model.deltaDecorations([], [getBreakPointOption(lineNumber, 'hide')])
          }
          // 存在断点，无事发生
          break
        }
        default: {
          removeTDecorations(model)
        }
      }
    }, 100),
  )
}

/**
 * 更新断点
 * @param editorInstance 编辑器实例
 * @returns breakpoints 断点数组
 */
export function updateBreakPoint(editorInstance: EditorInstance) {
  const model = editorInstance.getModel()
  if (!model) {
    return []
  }
  const decorations = model.getAllDecorations()
  const breakpoints = getBreakPointLineNumber(decorations)
  return breakpoints
}

/**
 * 初始化断点
 * @param model 要初始化的model
 * @param range 断点数组
 */
export const initBP = (editorInstance: EditorInstance, bp?: number[]) => {
  const model = editorInstance.getModel()
  if (!model) {
    return []
  }
  if (bp) {
    // 个人理解是第一个参数传之前的断点identifier，用来删除，第二个参数用来创建断点
    model.deltaDecorations(
      [],
      bp.map((lineNum) => {
        return {
          range: new Range(lineNum, 1, lineNum, 1),
          options: {
            isWholeLine: true,
            // className: "content", // 设置行的背景颜色（自己设置css，给类名）
            glyphMarginClassName: 'breakpoint active',
            stickiness: editor.TrackedRangeStickiness.NeverGrowsWhenTypingAtEdges,
            glyphMarginHoverMessage: { value: '断点' },
          },
        }
      }),
    )
  }
}
