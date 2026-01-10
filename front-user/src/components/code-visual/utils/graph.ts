import { reqStructVisual } from '@/api/visual'
import { StructVisualData, GraphDescription } from '@/api/visual/type'
import { GraphData, GraphNode } from '@/components/code-visual/visual/type/graph'

const reqGraphVisualData = async (
  debugID: string,
  description: GraphDescription,
): Promise<GraphData> => {
  // 请求可视化数据
  const req = {
    debugID: debugID,
    query: {
      struct: description.graphNode,
      values: [description.data],
      points: description.nexts,
    },
  }
  const result = await reqStructVisual(req)
  if (result.code != 200) {
    return {
      data: [],
      layouter: 'graph',
    }
  }
  // 将后端返回的可视化数据转成av需要的可视化数据
  return convertGraphVisualData(description, result.data)
}

const convertGraphVisualData = (
  _description: GraphDescription,
  data: StructVisualData,
): GraphData => {
  const nodes: GraphNode[] = []

  for (let i = 0; i < data.nodes.length; i++) {
    const node = data.nodes[i]
    const gnode: GraphNode = {
      id: '0',
    }
    gnode.id = node.id
    gnode.data = node.values[0].value
    for (let j = 0; j < node.points.length; j++) {
      for (const n of node.points) {
        if (!gnode.next) {
          gnode.next = []
        }
        gnode.next.push(n.value)
      }
    }
    nodes.push(gnode)
  }

  // 设置变量数据
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
    layouter: 'graph',
  }
}

// 请求数组可视化数据
export default reqGraphVisualData
