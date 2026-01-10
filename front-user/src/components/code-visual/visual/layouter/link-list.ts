import G6 from '@antv/g6'
import { SV, Group, SVNode, LayoutOptions } from 'structv2'

// 链表
SV.registerLayout('linkList', {
  sourcesPreprocess(sources) {
    if (sources.length == 0) {
      return sources
    }
    const root = sources[0]

    if (root.external) {
      root.rootExternal = root.external
      delete root.external
    }

    return sources
  },

  defineOptions() {
    return {
      node: {
        default: {
          type: 'link-list-node',
          label: '[data]',
          size: [60, 30],
          labelOptions: {
            style: {
              fontSize: 14,
              fontWeight: 400,
              fill: '#2c3e50',
            },
          },
          style: {
            stroke: '#8b5cf6',
            strokeWidth: 1,
            fill: '#f3e8ff',
            cursor: 'pointer',
            radius: 4,
            // 悬停效果
            hover: {
              stroke: '#7c3aed',
              strokeWidth: 2,
              fill: '#ede9fe',
            },
          },
        },
      },
      link: {
        next: {
          type: 'line',
          strokeWidth: 1,
          sourceAnchor: 6,
          targetAnchor: 5,
          style: {
            stroke: '#7f8c8d',
            endArrow: 'default',
            startArrow: {
              path: G6.Arrow.circle(2, -1),
              fill: '#7f8c8d',
            },
          },
        },
        prev: {
          type: 'line',
          strokeWidth: 1,
          sourceAnchor: 7,
          targetAnchor: 2,
          style: {
            stroke: '#7f8c8d',
            endArrow: 'default',
            startArrow: {
              path: G6.Arrow.circle(2, -1),
              fill: '#7f8c8d',
            },
          },
        },
      },
      marker: {
        rootExternal: {
          type: 'pointer',
          anchor: 0,
          offset: 8,
          style: {
            fill: '#f08a5d',
          },
        },
        external: {
          type: 'pointer',
          anchor: 0,
          offset: 8,
          style: {
            fill: '#f08a5d',
          },
        },
      },
      layout: {
        xInterval: 80,
        yInterval: 80,
      },
    }
  },

  /**
   * 对子树进行递归布局
   * @param node
   * @param parent
   */
  layoutItem(level, node, layoutOptions): [Group, number] {
    if (!node) {
      return null
    }

    const bound = node.getBound()
    const width = bound.width
    const group = new Group(node)
    let nextBound = null

    if (node.visited) {
      return [group, node.level]
    }

    node.visited = true
    node.level = level
    if (node.next) {
      const nextNodeVisited = node.next.visited
      const [nextGroup, nextNodeLevel] = this.layoutItem(level, node.next, layoutOptions)
      nextBound = nextGroup.getBound()
      node.set('x', nextBound.x - layoutOptions.xInterval - width)

      // 如果next本身已经渲染过，就不加入当前group中
      if (!nextNodeVisited) {
        group.add(nextGroup)
      }

      // 如果当前节点有写一个节点，当前节点的level等于下一个节点的level
      node.level = nextNodeLevel
    }

    return [group, node.level]
  },

  layout(elements: SVNode[], layoutOptions: LayoutOptions) {
    if (elements.length == 0) {
      return
    }
    const results = []
    // 收集result
    for (let i = 0; i < elements.length; i++) {
      const element = elements[i]
      if (element.root) {
        const result = this.layoutItem(i, element, layoutOptions)
        results.push(result)
      }
    }

    // 按照level排序
    results.sort((a, b) => a[1] - b[1])

    // 获取groups
    const groups = results.map((result) => result[0])

    // 避免多个链表相交
    for (let i = groups.length - 1; i > 0; i--) {
      const bound1 = groups[i].getBound()
      const bound2 = groups[i - 1].getBound()
      const move = Math.abs(bound1.y - layoutOptions.yInterval - bound2.y - bound2.height)
      groups[i - 1].translate(0, -(move + 20))
    }
  },
})
