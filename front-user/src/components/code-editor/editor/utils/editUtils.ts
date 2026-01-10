/* eslint-disable no-bitwise */
import { KeyCode, KeyMod } from 'monaco-editor'

import { EditorInstance } from '../types'

/**
 *
 * @param editorInstance 编辑器实例
 * 旧编辑器核心是重写插入、删除、移动文本方法，
 * 此编辑器不自带某一行只读的功能（见 https://github.com/microsoft/monaco-editor/issues/953 ），原理是利用编辑器自身的undo，redo栈
 * 好像可以利用decorations实现，但是我们应该没有这种需求，首先就需要拿到需要只读的行，目前是没有的
 */
export const firstLineReadOnly = (editorInstance: EditorInstance) => {
  const model = editorInstance.getModel()
  model?.onDidChangeContent((e) => {
    // 不处理好像也行
    if (e.isUndoing) {
      return
    }

    /**
     * 某次修改中，可能涉及到多个地方的修改，比如拖动就涉及到原位置的修改和目标位置的修改，所以是数组。
     * 只要涉及到第一行的修改，就取消本次修改
     */
    const changeFirstLine = e.changes.some((change) => change.range.startLineNumber === 1)
    if (changeFirstLine) {
      // 使用微任务，无法恢复/撤销，使用宏任务的话输入内容会闪烁，直接调用的情况下ctrl + y一样能恢复
      queueMicrotask(() => {
        // 返回上一次编辑内容，不知道第一个参数啥意思，第二个参数是核心，不能改
        editorInstance.trigger('firstLineReadOnly', 'undo', null)
      })
    }
  })
}

/**
 * tips: 实现起来和上面差不多，只不过以decorations为中心（断点就是通过decorations实现的）
 * @param editorInstance 编辑器实例
 * @param readOnlyLine 只读行数组
 */
export const lineReadOnly = (editorInstance: EditorInstance, readOnlyLine: number[]) => {
  // TODO: 留给有缘人写
  editorInstance
  readOnlyLine
}

/**
 * 添加一个对ctrl+s的监听
 * @param editorInstance 通过monaco创建的编辑器实例
 * @param callback 执行回调函数
 */
export const addSaveCodeListener = (
  editorInstance: EditorInstance,
  callback = (value: string) => {
    console.log(`使用了ctrlS${value}`)
  },
) => {
  editorInstance.addCommand(KeyMod.CtrlCmd | KeyCode.KeyS, () => {
    const value = editorInstance.getValue()
    callback(value)
  })
}

/**
 * 设置编辑器当前model的内容（不可频繁调用，设置的时候编辑器会闪烁一下）
 * @param editorInstance 通过monaco创建的编辑器实例
 * @param value 要设置的内容
 */
export const setContent = (editorInstance: EditorInstance, value: string) => {
  editorInstance.setValue(value)
}

/**
 * 在失焦时保存，未投入使用
 * @param editorInstance 编辑器实例
 */
export const saveOnBlur = (editorInstance: EditorInstance) => {
  editorInstance.onDidBlurEditorWidget(() => {
    const value = editorInstance.getValue()
    setContent(editorInstance, value)
  })
}

export default {
  saveOnBlur,
  addSaveCodeListener,
  firstLineReadOnly,
  lineReadOnly,
  setContent,
}
