import { SV } from 'structv2'

SV.registerLayout('arrayBar', {
  defineOptions() {
    return {
      node: {
        default: {
          type: 'array-node',
          label: '[data]',
          size: [40, 80], // 柱状图默认尺寸，将在layout中动态调整
          labelOptions: {
            style: {
              fontSize: 12,
              fontWeight: 400,
              fill: '#2c3e50',
            },
          },
          style: {
            stroke: '#10b981',
            strokeWidth: 1,
            fill: '#d1fae5',
            cursor: 'pointer',
            radius: 4,
            // 悬停效果
            hover: {
              stroke: '#059669',
              strokeWidth: 2,
              fill: '#a7f3d0',
            },
          },
        },
      },
      marker: {
        headExternal: {
          type: 'pointer',
          anchor: 7,
          rotation: -Math.PI / 2, // 水平指向（0弧度 = 0度）
          style: {
            fill: '#e74c3c',
            stroke: '#c0392b',
            strokeWidth: 2,
          },
          labelOptions: {
            style: {
              fontSize: 12,
              fontWeight: 500,
              fill: '#7f8c8d',
            },
          },
        },
        external: {
          type: 'pointer',
          anchor: 2,
          offset: 30,
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
          style: { fontSize: 10, fill: '#666' },
        },
        indexRight: { position: 'right' },
      },
      behavior: {
        dragNode: false,
      },
    }
  },

  layout(elements) {
    const arr = elements

    // 固定柱状图宽度和间距
    const barWidth = 40 // 固定柱状图宽度
    const spacing = 10 // 固定间距

    // 智能高度计算逻辑 - 固定算法，不依赖最大值
    function calculateIntelligentHeight(value, minHeight, maxHeight) {
      // 使用对数函数实现前面变化快，后面变化慢的效果
      // 使用 log(1 + value) 来确保从0开始，且增长速度逐渐变慢
      const logValue = Math.log(1 + value)

      // 设置一个合理的缩放因子，让高度变化在合理范围内
      const scaleFactor = 30 // 控制整体缩放

      // 计算高度：使用对数函数，前面增长快，后面增长慢
      let actualHeight = logValue * scaleFactor

      // 确保最小和最大高度
      actualHeight = Math.max(actualHeight, minHeight)
      actualHeight = Math.min(actualHeight, maxHeight)

      return actualHeight
    }

    // 更新每个元素的尺寸和位置
    for (let i = 0; i < arr.length; i++) {
      const value = Number(arr[i].get('data')) || 0
      const minHeight = 10 // 基础高度
      const maxHeight = 300

      // 计算智能高度
      const intelligentHeight = calculateIntelligentHeight(value, minHeight, maxHeight)

      // 设置动态高度
      arr[i].set('size', [barWidth, intelligentHeight])

      // 计算x坐标
      if (i > 0) {
        arr[i].set('x', arr[i - 1].get('x') + arr[i - 1].get('size')[0] + spacing)
      } else {
        arr[i].set('x', 0)
      }

      arr[i].set('y', -intelligentHeight / 2)
    }
  },
})
