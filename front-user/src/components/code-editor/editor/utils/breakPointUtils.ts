import { editor, Range } from 'monaco-editor'

import { Model } from '../types'

/**
 * 判断是否存在断点
 * @param decorations 通过model.getLineDecorations获取的结果
 * @param transparent 是否验证判断存在透明度断点是拥有透明度的断点
 * @returns 存在或不存在
 */
export const existBreakPoint = (decorations: editor.IModelDecoration[], transparent = false) => {
  const str = transparent ? 'hide' : 'active'
  // 判断这行有没有断点
  return decorations.some((item) => item.options.glyphMarginClassName?.split(' ').includes(str))
}

// TODO: 重构下，封装抽取这几个字符串，不好维护
type getBreakPointsIdsOptions = 'all' | 'BP' | 'TBP'
export const getBreakPointsIds = (
  decorations: editor.IModelDecoration[],
  options: getBreakPointsIdsOptions = 'all',
) => {
  const map = {
    BP: 'active',
    TBP: 'hide',
    all: 'breakpoint',
  }
  return decorations
    .filter((item) => item.options.glyphMarginClassName?.split(' ').includes(map[options]))
    .map((item) => item.id)
}

/**
 *
 * @param model 编辑器的model
 * @description 删除所有存在透明度的断点
 */
export const removeTDecorations = (model?: Model | null) => {
  if (!model) {
    return
  }
  const allDecorations = model.getAllDecorations()
  const TBPIds = getBreakPointsIds(allDecorations, 'TBP')
  model.deltaDecorations(TBPIds, [])
}

/**
 *
 * @param lineStr 行字符串
 * @return 检查是否通过
 * @description 传入编辑器某一行的内容，返回检查通过与否
 */
export const checkBreakPoints = (lineStr: string) => {
  const patterns = [
    /^[\s]*$/, // 0.全为空白符
    /^[\s]*[(){}]+[\s]*$/, // 1.只有括号的行
    /^[\s]*#[a-zA-Z]+.*$/, // 2.形如"#include"、"#define"开头的代码行
    /.*\/\*.*\*\//, // 3."/* */"在同一行
    /^[\s(){}]*(\/\/|\/\*)/, // 4."//"或者"/*"之前全为空白符或者括号
    /\*\/[\s(){}]*$/, // 5."*/"之后全为空白符或者括号
  ]
  // every是一个false终止，some是一个true终止
  return patterns.every((value) => !value.test(lineStr))
}

/**
 * 获取断点的选项（添加断点需要的数组）
 * @param lineNumber 行号
 * @param mode 断点模式
 * @returns 断点选项
 */
export const getBreakPointOption = (lineNumber: number, mode: 'hide' | 'active' = 'active') => {
  return <editor.IModelDeltaDecoration>{
    range: new Range(lineNumber, 1, lineNumber, 1),
    options: {
      isWholeLine: true,
      glyphMarginClassName: `breakpoint ${mode}`,
      stickiness: editor.TrackedRangeStickiness.NeverGrowsWhenTypingAtEdges,
      glyphMarginHoverMessage: { value: '断点' },
    },
  }
}

/**
 *
 * @param decorations 通过model.getAllDecorations获取
 * @returns lineNumber数组
 */
export const getBreakPointLineNumber = (decorations: editor.IModelDeltaDecoration[]) =>
  decorations
    .filter((decoration) => decoration.options.glyphMarginClassName?.split(' ').includes('active'))
    // 这里startLineNumber和endLineNumber都是一样的，因为new Range(lineNumber, 1, lineNumber, 1)
    .map((decoration) => decoration.range.startLineNumber)
