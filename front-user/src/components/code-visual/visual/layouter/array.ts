import { SV, SourceNode, LayoutGroupOptions, SVNode, LayoutOptions } from 'structv2'

// 数组
SV.registerLayout('array', {
  defineLeakRule(_models: SVNode[]): SVNode[] {
    return []
  },

  defineOptions(_sourceData: SourceNode[]): LayoutGroupOptions {
    return {
      node: {
        default: {
          type: 'array-node',
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
            stroke: '#5dade2',
            strokeWidth: 1,
            fill: '#ebf3fd',
            cursor: 'pointer',
            radius: 4,
            // 悬停效果
            hover: {
              stroke: '#3498db',
              strokeWidth: 2,
              fill: '#d6eaf8',
            },
          },
        },
        virtual: {
          type: 'array-node',
          label: '[data]',
          size: [60, 30],
          labelOptions: {
            style: {
              fontSize: 14,
              fontWeight: 400,
              fill: '#95a5a6',
            },
          },
          style: {
            stroke: '#bdc3c7',
            strokeWidth: 1,
            lineDash: [3, 4],
            fill: 'transparent',
            radius: 4,
          },
        },
      },
      marker: {
        headExternal: {
          type: 'pointer',
          anchor: 3,
          style: {
            fill: '#e74c3c',
            stroke: '#c0392b',
            strokeWidth: 1,
          },
        },
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
      indexLabel: {
        index: {
          position: 'bottom',
          style: {
            fontSize: 12,
            fill: '#7f8c8d',
            fontWeight: 400,
          },
        },
        indexRight: {
          position: 'right',
          style: {
            fontSize: 12,
            fill: '#7f8c8d',
            fontWeight: 400,
          },
        },
      },
      behavior: {
        dragNode: false,
      },
    }
  },

  layout(nodes: SVNode[], _layoutOptions: LayoutOptions) {
    const arr = nodes

    for (let i = 0; i < arr.length; i++) {
      const width = arr[i].get('size')[0]
      if (i > 0) {
        arr[i].set('x', arr[i - 1].get('x') + width)
      }
    }
  },
})
