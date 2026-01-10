export type GraphData = {
  data: GraphNode[]
  layouter: 'graph'
}

export type GraphNode = {
  id: string
  external?: string[]
  // 数据
  data?: string
  next?: string[]
}
