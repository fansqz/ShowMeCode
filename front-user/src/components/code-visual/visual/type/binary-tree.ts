export type BinaryTreeData = {
  data: BinaryTreeNode[]
  layouter: 'binaryTree'
}

export type BinaryTreeNode = {
  id: string
  root?: boolean
  external?: string[]
  // 数据
  data?: string
  child?: string[]
}
