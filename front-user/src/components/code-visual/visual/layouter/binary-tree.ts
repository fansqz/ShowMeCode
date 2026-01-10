import G6 from '@antv/g6'
import { SV, SourceNode, LayoutGroupOptions, Group, SVNode, LayoutOptions } from 'structv2'

SV.registerLayout('binaryTree', {
  defineOptions(_sourceData: SourceNode[]): LayoutGroupOptions {
    return {
      node: {
        default: {
          type: 'binary-tree-node',
          size: [60, 30],
          label: '[data]',
          labelOptions: {
            style: {
              fontSize: 14,
              fontWeight: 400,
              fill: '#2c3e50',
            },
          },
          style: {
            stroke: '#06b6d4',
            strokeWidth: 1,
            fill: '#cffafe',
            cursor: 'pointer',
            radius: 4,
            // 悬停效果
            hover: {
              stroke: '#0891b2',
              strokeWidth: 2,
              fill: '#a5f3fc',
            },
          },
        },
      },
      link: {
        child: {
          type: 'line',
          sourceAnchor: (index) => (index === 0 ? 3 : 1),
          targetAnchor: 0,
          style: {
            stroke: '#7f8c8d',
            strokeWidth: 1,
            lineAppendWidth: 6,
            cursor: 'pointer',
            endArrow: 'default',
            preventOverlap: true, //防重叠
            startArrow: {
              path: G6.Arrow.circle(2, -1),
              fill: '#7f8c8d',
            },
          },
        },
      },
      marker: {
        external: {
          type: 'pointer',
          anchor: 0,
          style: {
            fill: '#f39c12',
            stroke: '#d68910',
            strokeWidth: 1,
          },
        },
      },
      layout: {
        xInterval: 50,
        yInterval: 50,
      },
    }
  },

  /**
   * 对子树进行递归布局
   */
  layoutItem(node: SVNode, layoutOptions: LayoutOptions) {
    // 次双亲不进行布局
    if (!node) {
      return null
    }

    const bound = node.getBound()
    const width = bound.width
    const height = bound.height
    const group = new Group(node)
    let leftGroup = null
    let rightGroup = null
    let leftBound = null
    let rightBound = null

    if (node.visited) {
      return null
    }

    node.visited = true

    if (node.child && node.child[0]) {
      leftGroup = this.layoutItem(node.child[0], layoutOptions)
    }

    if (node.child && node.child[1]) {
      rightGroup = this.layoutItem(node.child[1], layoutOptions)
    }

    if (leftGroup) {
      leftBound = leftGroup.getBound()
      node.set('y', leftBound.y - layoutOptions.yInterval - height)
    }

    if (rightGroup) {
      rightBound = rightGroup.getBound()

      if (leftGroup) {
        rightGroup.translate(0, leftBound.y - rightBound.y)
      }

      rightBound = rightGroup.getBound()
      node.set('y', rightBound.y - layoutOptions.yInterval - height)
    }

    // 处理左右子树相交问题
    if (leftGroup && rightGroup) {
      const move = Math.abs(rightBound.x - layoutOptions.xInterval - leftBound.x - leftBound.width)
      if (move > 0) {
        leftGroup.translate(-move / 2, 0)
        rightGroup.translate(move / 2, 0)
      }
    }

    if (leftGroup) {
      leftBound = leftGroup.getBound()
      node.set('x', leftBound.x + leftBound.width + layoutOptions.xInterval / 2 - width)
    }

    if (rightGroup) {
      rightBound = rightGroup.getBound()
      node.set('x', rightBound.x - layoutOptions.xInterval / 2 - width)
    }

    if (leftGroup) {
      group.add(leftGroup)
    }

    if (rightGroup) {
      group.add(rightGroup)
    }

    return group
  },

  layout(elements: SVNode[], layoutOptions: LayoutOptions) {
    const groups = []
    // 解决多棵树相交问题
    for (const element of elements) {
      if (element.root) {
        const group = this.layoutItem(element, layoutOptions)
        groups.push(group)
      }
    }
    // 避免多棵相交
    for (let i = 0; i < groups.length - 1; i++) {
      const bound1 = groups[i].getBound()
      const bound2 = groups[i + 1].getBound()
      const move = Math.abs(bound2.x - layoutOptions.xInterval - bound1.x - bound1.width)
      groups[i + 1].translate(move, 0)
    }
  },
})
