import { defineStore } from 'pinia'

import {
  Array2DDescription,
  ArrayDescription,
  BinaryTreeDescription,
  GraphDescription,
  LinkListDescription,
  VisualDescription,
} from '@/api/visual/type'
import { descriptions } from '@/constants/description.ts'

type VisualState = {
  // 可视化是否开启
  action: boolean
  // 是否ai自动识别可视化
  isAIEnabled: boolean
  // 可视化描述类型
  descriptionType: descriptions
  // 可视化描述映射
  descriptions: {
    [descriptions.Array]: ArrayDescription
    [descriptions.Array2D]: Array2DDescription
    [descriptions.BinaryTree]: BinaryTreeDescription
    [descriptions.LinkList]: LinkListDescription
    [descriptions.Graph]: GraphDescription
  }
}

const useVisualStore = defineStore('visual', {
  state: (): VisualState => ({
    // 可视化是否开启
    action: false,
    // 是否ai自动识别可视化
    isAIEnabled: true,
    descriptionType: descriptions.Array,
    descriptions: {
      [descriptions.Array]: {
        arrayName: 'Arr',
        pointNames: ['Point1', 'Point2'],
        displayType: 'array-bar',
      },
      [descriptions.Array2D]: {
        arrayName: 'Arr',
        rowPointNames: ['i'],
        colPointNames: ['j'],
      },

      [descriptions.BinaryTree]: {
        treeNode: 'TreeNode',
        data: 'Val',
        left: 'Left',
        right: 'Right',
      },
      [descriptions.LinkList]: {
        // 链表节点
        linkNode: 'LinkNode',
        // 数据域
        data: 'Val',
        next: 'Next',
        prev: 'Prev',
      },
      [descriptions.Graph]: {
        graphNode: 'GraphNode',
        // 数据域
        data: 'Val',
        nexts: ['Point1', 'Point2'],
      },
    },
  }),
  actions: {
    setDescription(descriptionType: descriptions, description: VisualDescription) {
      this.descriptions[descriptionType] = description as any
    },
    getDescription(descriptionType: descriptions) {
      return this.descriptions[descriptionType]
    },
  },
})

export default useVisualStore
