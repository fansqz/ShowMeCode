import { reqArrayVisual } from '@/api/visual'
import {
  ArrayVisualData,
  VisualVariable,
  ArrayDescription,
  ArrayVisualRequest,
} from '@/api/visual/type'
import { ArrayData, ArrayNode } from '@/components/code-visual/visual/type/array'

const reqArrayVisualData = async (
  debugID: string,
  description: ArrayDescription,
): Promise<ArrayData> => {
  const req: ArrayVisualRequest = {
    debugID: debugID,
    query: {
      arrayName: description.arrayName,
      pointNames: description.pointNames,
    },
  }
  const result = await reqArrayVisual(req)
  if (result.code != 200) {
    return {
      data: [],
      layouter: 'array',
    }
  }
  // 将后端返回的可视化数据转成av需要的可视化数据
  return convertArrayVisualData(description, result.data)
}

const convertArrayVisualData = (
  description: ArrayDescription,
  data: ArrayVisualData,
): ArrayData => {
  // 将第一个struct[0]作为数组的数据
  const struct = data.array
  const nodes: ArrayNode[] = []
  for (let i = 0; i < struct.length; i++) {
    const node: ArrayNode = {
      id: struct[i].name,
      index: struct[i].name,
      data: struct[i].value,
    }
    if (i === 0) {
      // 头指针
      node.headExternal = description.arrayName
    }
    nodes.push(node)
  }

  const virtualNodePoints: VisualVariable[] = []

  // 设置指针
  for (let i = 0; i < data.points.length; i++) {
    const index = data.points[i].value
    if (nodes.length <= Number(index) || Number(index) < 0) {
      virtualNodePoints.push(data.points[i])
      continue
    }
    if (nodes[index].external != undefined) {
      nodes[index].external.push(data.points[i].name)
    } else {
      nodes[index].external = [data.points[i].name]
    }
  }

  // 设置虚拟节点
  const maxPoint = virtualNodePoints.reduce((max, point) => {
    if (Number(point.value) > max) {
      return Number(point.value)
    }
    return max
  }, 0)
  let count = 0
  for (let i = nodes.length; i <= maxPoint; i++) {
    nodes.push({
      id: `[${String(i + 1)}]`,
      index: `[${String(i + 1)}]`,
      data: '',
      type: 'virtual',
    })
    // 只支持5个虚拟节点
    count++
    if (count > 5) {
      break
    }
  }

  // 设置虚拟节点的外部指针
  for (let i = 0; i < virtualNodePoints.length; i++) {
    const index = virtualNodePoints[i].value
    if (nodes.length <= Number(index) || Number(index) < 0) {
      continue
    }
    if (nodes[index].external != undefined) {
      nodes[index].external.push(virtualNodePoints[i].name)
    } else {
      nodes[index].external = [virtualNodePoints[i].name]
    }
  }

  // 根据展示方式选择布局器
  const layouter = description.displayType === 'array-bar' ? 'arrayBar' : 'array'

  return {
    data: nodes,
    layouter: layouter,
  }
}

// 请求数组可视化数据
export default reqArrayVisualData
