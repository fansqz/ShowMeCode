import { descriptions } from '@/constants/description.ts'

import { BaseResponse } from '../type'

// 结构体导向可视化请求
export type StructVisualRequest = {
  debugID: string
  query: StructVisualQuery
}

export type StructVisualQuery = {
  struct: string
  values: string[]
  points: string[]
}

export type StructVisualResponse = BaseResponse & {
  data: StructVisualData
}

export type StructVisualData = {
  // 可视化节点
  nodes: VisualNode[]
  // 指针列表
  points: VisualVariable[]
}

// 数组可视化请求
export type ArrayVisualRequest = {
  debugID: string
  query: ArrayVisualQuery
}

export type ArrayVisualQuery = {
  arrayName: string
  pointNames: string[]
}

export type ArrayVisualResponse = BaseResponse & {
  data: ArrayVisualData
}

// 数组指针的可视化数据，比如数组
export type ArrayVisualData = {
  // 结构体列表
  array: VisualVariable[]
  // 指针列表
  points: VisualVariable[]
}

export type Array2DVisualRequest = {
  debugID: string
  query: Array2DVisualQuery
}

export type Array2DVisualQuery = {
  arrayName: string
  rowPointNames: string[]
  colPointNames: string[]
}

export type Array2DVisualResponse = BaseResponse & {
  data: Array2DVisualData
}

// 数组指针的可视化数据，比如数组
export type Array2DVisualData = {
  // 结构体列表
  array: VisualVariable[][]
  // 指针列表
  rowPoints: VisualVariable[]
  colPoints: VisualVariable[]
}

// 可视化节点
export type VisualNode = {
  name: string
  id: string
  type: string
  // 值域
  values: VisualVariable[]
  // 指针域
  points: VisualVariable[]
}

// 可视化变量
export type VisualVariable = {
  name: string
  type: string
  value: string
}

export type VisaulDocumentResponse = BaseResponse & {
  data: string
}

export type VisualDescriptionResponse = BaseResponse & {
  data: {
    visualType: descriptions
    description: VisualDescription
  }
}

// 可视化描述类型
export type VisualDescription =
  | ArrayDescription
  | Array2DDescription
  | BinaryTreeDescription
  | LinkListDescription
  | GraphDescription
  | any

// 二维数组可视化描述
export type Array2DDescription = {
  arrayName: string
  rowPointNames: string[]
  colPointNames: string[]
}

// 数组可视化描述
export type ArrayDescription = {
  arrayName: string
  pointNames: string[]
  displayType?: 'array' | 'array-bar'
}

// 二叉树可视化描述
export type BinaryTreeDescription = {
  // 二叉树节点结构体名称
  treeNode: string
  // 数据域
  data: string
  // 左子树和右边子树属性名称
  left: string
  right: string
}

// 图的可视化描述
export type GraphDescription = {
  // 二叉树节点结构体名称
  graphNode: string
  // 数据域
  data: string
  nexts: string[]
}

// 链表的可视化描述
export type LinkListDescription = {
  // 链表节点
  linkNode: string
  // 数据域
  data: string
  next: string
  prev?: string
}
