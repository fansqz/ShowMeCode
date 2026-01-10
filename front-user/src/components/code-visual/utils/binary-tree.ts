import { reqStructVisual } from '@/api/visual'
import { StructVisualData, BinaryTreeDescription } from '@/api/visual/type'
import { BinaryTreeData, BinaryTreeNode } from '@/components/code-visual/visual/type/binary-tree'

const reqBinaryTreeVisualData = async (
  debugID: string,
  description: BinaryTreeDescription,
): Promise<BinaryTreeData> => {
  // 请求可视化数据
  const req = {
    debugID: debugID,
    query: {
      struct: description.treeNode,
      values: [description.data],
      points: [description.left, description.right],
    },
  }
  const result = await reqStructVisual(req)
  if (result.code != 200) {
    return {
      data: [],
      layouter: 'binaryTree',
    }
  }
  // 将后端返回的可视化数据转成av需要的可视化数据
  return convertqBinaryTreeVisualData(description, result.data)
}

const convertqBinaryTreeVisualData = (
  description: BinaryTreeDescription,
  data: StructVisualData,
): BinaryTreeData => {
  const nodes: BinaryTreeNode[] = []

  // 遍历所有二叉树节点，因为二叉树的渲染要从根节点开始一直到子节点，所以必须根节点在上面
  const nextIDSet = new Set()
  for (let i = 0; i < data.nodes.length; i++) {
    const node = data.nodes[i]
    const btNode: BinaryTreeNode = {
      id: '0',
    }
    btNode.id = node.id
    btNode.data = node.values[0].value
    for (let j = 0; j < node.points.length; j++) {
      if (btNode.child == null) {
        btNode.child = []
      }
      if (node.points[j].name == description.left) {
        btNode.child[0] = node.points[j].value
      } else {
        btNode.child[1] = node.points[j].value
      }
      nextIDSet.add(node.points[j].value)
    }
    nodes.push(btNode)
  }

  // 找出根节点，如果没有谁指向该节点，这个节点就是根节点。
  let existRoot = false
  for (const node of nodes) {
    if (!nextIDSet.has(node.id)) {
      node.root = true
      existRoot = true
    }
  }
  // 不存在根节点，则第一个节点设置为根节点
  if (!existRoot) {
    nodes[0].root = true
  }

  // 设置指针
  for (let i = 0; i < data.points.length; i++) {
    const pointer = data.points[i]
    for (let j = 0; j < nodes.length; j++) {
      const node = nodes[j]
      if (node.id == pointer.value) {
        if (node.external == undefined) {
          node.external = []
        }
        nodes[j].external.push(pointer.name)
      }
    }
  }
  return {
    data: nodes,
    layouter: 'binaryTree',
  }
}

// 请求数组可视化数据
export default reqBinaryTreeVisualData
