import { reqStructVisual } from '@/api/visual'
import { StructVisualData, LinkListDescription } from '@/api/visual/type'
import { LinkListData, LinkListNode } from '@/components/code-visual/visual/type/link-list'

const reqLinkListVisualData = async (
  debugID: string,
  description: LinkListDescription,
): Promise<LinkListData> => {
  const points = [description.next]
  if (description.prev) {
    points.push(description.prev)
  }
  // 请求可视化数据
  const req = {
    debugID: debugID,
    query: {
      struct: description.linkNode,
      values: [description.data],
      points: points,
    },
  }
  const result = await reqStructVisual(req)
  if (result.code != 200) {
    console.log(result.message)
    return {
      data: [],
      layouter: 'linkList',
    }
  }
  // 将后端返回的可视化数据转成av需要的可视化数据
  return convertLinkListVisualData(description, result.data)
}

const convertLinkListVisualData = (
  description: LinkListDescription,
  data: StructVisualData,
): LinkListData => {
  if (data.nodes.length == 0) {
    return {
      data: [],
      layouter: 'linkList',
    }
  }
  const nodes: LinkListNode[] = []

  const nextSet = new Set()

  // 转换节点数据
  for (let i = 0; i < data.nodes.length; i++) {
    const node = data.nodes[i]
    const lnode: LinkListNode = { id: '0', next: '' }
    lnode.id = node.id
    if (node.values && node.values.length != 0) {
      lnode.data = node.values[0].value
    }
    for (let j = 0; j < node.points.length; j++) {
      for (const p of node.points) {
        // 设置next指针
        if (p.name == description.next) {
          lnode.next = p.value
          nextSet.add(lnode.next)
        }
        // 设置pre指针
        if (p.name == description.prev) {
          lnode.prev = p.value
        }
      }
    }
    nodes.push(lnode)
  }

  let hasRoot = false
  // 寻找根节点，无任何next指向的节点为根节点
  for (let i = 0; i < nodes.length; i++) {
    if (!nextSet.has(nodes[i].id)) {
      nodes[i].root = true
      hasRoot = true
    }
  }

  // 如果没有root，则第一个节点为root
  if (!hasRoot) {
    nodes[0].root = true
  }

  // 设置变量数据
  for (let i = 0; i < data.points.length; i++) {
    const pointer = data.points[i]
    for (let j = 0; j < nodes.length; j++) {
      const node = nodes[j]
      // 当前节点的地址等于变量指向的地址，设置变量指针指向该节点
      if (node.id == pointer.value) {
        if (node.external == undefined) {
          node.external = []
        }
        nodes[j].external.push(pointer.name)
      }
    }
  }

  // 对节点按ID进行排序
  nodes.sort((a, b) => {
    // 将ID转换为数字进行比较
    const idA = parseInt(a.id)
    const idB = parseInt(b.id)
    return idA - idB
  })

  return {
    data: nodes,
    layouter: 'linkList',
  }
}

// 请求数组可视化数据
export default reqLinkListVisualData
