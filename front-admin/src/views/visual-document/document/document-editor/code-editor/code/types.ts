import { editor } from 'monaco-editor';
import { Ref } from 'vue';

// editor实例
export type EditorInstance = editor.IStandaloneCodeEditor;

export type Model = editor.ITextModel;

export type ThemeData = editor.IStandaloneThemeData;

export type onUpdateBP = (breakpoints: number[], lineNum: number, mode: 'add' | 'del') => void;

// vscode创建的请求
export type VsCode = {
  // 目标元素
  target: Ref<HTMLElement | undefined>;
  // ===========以下是可选的===============
  // 内容修改时的回调
  onContentChanged?: (value: string) => void;
  // 失去焦点的回调
  onEditorBlur?: (value: string) => void;
  // ctrl + s的回调
  onCtrlS?: (value: string) => void;
  onSetBP?: (breakpoints?: number[]) => void;
  // 断点改变
  onUpdateBP?: onUpdateBP;
};

export type DebugInfo = {
  model: Model | null;
  oldModel?: Model | null;
  /** 当前高亮行 */
  currentline?: number;
  /** 上一次的高亮行 */
  preline?: number;
  instance: EditorInstance;
};
