import { Source } from 'structv2'

import { VisualDescription } from '@/api/visual/type'
import { descriptions } from '@/constants/description'

import reqArrayVisualData from './array'
import reqArray2DVisualData from './array2d'
import reqBinaryTreeVisualData from './binary-tree'
import reqGraphVisualData from './graph'
import reqLinkListVisualData from './link-list'

export const reqVisualData = async (
  debugID: string,
  type: descriptions,
  description: VisualDescription,
): Promise<Source> => {
  if (type == descriptions.Array) {
    return reqArrayVisualData(debugID, description)
  } else if (type == descriptions.Array2D) {
    return reqArray2DVisualData(debugID, description)
  } else if (type == descriptions.BinaryTree) {
    return reqBinaryTreeVisualData(debugID, description)
  } else if (type == descriptions.Graph) {
    return reqGraphVisualData(debugID, description)
  } else if (type == descriptions.LinkList) {
    return reqLinkListVisualData(debugID, description)
  }
}
