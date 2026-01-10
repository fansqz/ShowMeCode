import { reqArray2DVisual } from '@/api/visual'
import { Array2DVisualData, Array2DDescription, Array2DVisualRequest } from '@/api/visual/type'
import { Array2DData, Array2DNode } from '@/components/code-visual/visual/type/array2d'

const reqArray2DVisualData = async (
  debugID: string,
  description: Array2DDescription,
): Promise<Array2DData> => {
  // 请求可视化数据
  const req: Array2DVisualRequest = {
    debugID: debugID,
    query: {
      arrayName: description.arrayName,
      colPointNames: description.colPointNames,
      rowPointNames: description.rowPointNames,
    },
  }
  const result = await reqArray2DVisual(req)
  if (result.code != 200) {
    return {
      data: [],
      layouter: 'array2d',
    }
  }
  // 将后端返回的可视化数据转成av需要的可视化数据
  return convertArrayVisualData(description, result.data)
}

const convertArrayVisualData = (
  description: Array2DDescription,
  data: Array2DVisualData,
): Array2DData => {
  const struct = data.array
  const nodes: Array2DNode[] = []
  for (let i = 0; i < struct.length; i++) {
    const array = struct[i]
    for (let j = 0; j < array.length; j++) {
      const arrayNode = array[j]
      const node: Array2DNode = {
        row: i,
        col: j,
        id: `[${i}][${j}]`,
        data: arrayNode.value,
      }
      if (i === 0 && j === 0) {
        // 头指针
        node.headExternal = description.arrayName
      }
      // 设置index
      if (i == struct.length - 1) {
        node.indexBottom = `[${j}]`
      }
      if (j == struct[0].length - 1) {
        node.indexRight = `[${i}]`
      }
      nodes.push(node)
    }
  }

  // 设置指针
  for (let i = 0; i < data.rowPoints.length; i++) {
    const index = data.rowPoints[i].value
    const pointName = data.rowPoints[i].name
    for (let j = 0; j < nodes.length; j++) {
      if (nodes[j].row == Number(index) && nodes[j].col == 0) {
        if (nodes[j].rowExternal != undefined) {
          nodes[j].rowExternal.push(pointName)
        } else {
          nodes[j].rowExternal = [pointName]
        }
      }
    }
  }

  for (let i = 0; i < data.colPoints.length; i++) {
    const index = data.colPoints[i].value
    const pointName = data.colPoints[i].name
    for (let j = 0; j < nodes.length; j++) {
      if (nodes[j].col == Number(index) && nodes[j].row == 0) {
        if (nodes[j].colExternal != undefined) {
          nodes[j].colExternal.push(pointName)
        } else {
          nodes[j].colExternal = [pointName]
        }
      }
    }
  }
  return {
    data: nodes,
    layouter: 'array2d',
  }
}

// 请求数组可视化数据
export default reqArray2DVisualData
