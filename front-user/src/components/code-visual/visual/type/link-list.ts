export type LinkListData = {
  data: LinkListNode[]
  layouter: 'linkList'
}

export type LinkListNode = {
  id: string
  root?: boolean
  external?: string[]
  // 数据
  data?: string
  next: string
  prev?: string
}
