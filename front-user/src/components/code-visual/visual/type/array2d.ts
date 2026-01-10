// Array 传递的数据类型
export type Array2DData = {
  data: Array2DNode[]
  layouter: 'array2d'
}
export type Array2DNode = {
  id: string
  col: number
  row: number
  headExternal?: string
  // 外部指针
  colExternal?: string[]
  rowExternal?: string[]
  // 标题
  indexBottom?: string
  indexRight?: string
  // 数据
  data?: string
}
