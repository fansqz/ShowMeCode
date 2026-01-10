// Array 传递的数据类型
export type ArrayData = {
  data: ArrayNode[]
  layouter: 'array' | 'arrayBar'
}

export type ArrayNode = {
  type?: 'default' | 'virtual'
  id: string
  headExternal?: string
  // 外部指针
  external?: string[]
  // 标题，可以理解为数组的索引
  index?: string
  // 数据
  data?: string
}
