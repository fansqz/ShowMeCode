<template>
  <!-- 手动触发指引的按钮 -->
  <div
    class="debug-guide-trigger"
    :style="{
      top: buttonPosition.top + 'px',
      right: buttonPosition.right + 'px',
    }"
    @mousedown="startDrag"
  >
    <el-button
      type="primary"
      :icon="QuestionFilled"
      circle
      size="small"
      @click="handleButtonClick"
      class="guide-button"
    />
  </div>

  <!-- Element UI 用户引导组件 -->
  <el-tour v-model="tourVisible" @close="handleTourClose" @finish="handleTourFinish">
    <el-tour-step :target="getBreakpointArea" title="步骤1: 添加断点" placement="right">
      <div class="tour-content">
        <p>在代码编辑器的行号左侧点击，可以添加断点。断点用于在特定位置暂停程序执行，方便调试。</p>
        <ul>
          <li>
            <span class="debug-icon">🔴</span>
            <strong>断点</strong>：在行号左侧点击设置断点
          </li>
        </ul>
      </div>
    </el-tour-step>

    <el-tour-step target=".debug-button" title="步骤2: 启动调试" placement="top">
      <div class="tour-content">
        <ul>
          <li>
            <span class="debug-icon iconfont icon-caozuo-tiaoshi"></span>
            <strong>启动调试</strong>：点击开始可视化调试
          </li>
        </ul>
      </div>
    </el-tour-step>

    <el-tour-step target=".debug-buttons" title="步骤3: 调试控制" placement="top">
      <div class="tour-content">
        <p>使用这些按钮来控制调试过程：</p>
        <ul>
          <li>
            <span class="debug-icon iconfont icon-continue1"></span>
            <strong>继续执行</strong>：继续运行到下一个断点
          </li>
          <li>
            <span class="debug-icon iconfont icon-step"></span>
            <strong>单步跳过</strong>：执行当前行，不进入函数内部
          </li>
          <li>
            <span class="debug-icon iconfont icon-step-in"></span>
            <strong>单步进入</strong>：进入函数内部执行
          </li>
          <li>
            <span class="debug-icon iconfont icon-step-out"></span>
            <strong>单步退出</strong>：跳出当前函数
          </li>
        </ul>
      </div>
    </el-tour-step>

    <el-tour-step target=".console-terminal" title="步骤4: 查看调试信息" placement="top">
      <div class="tour-content">
        <p>在调试终端中，你可以查看：</p>
        <ul>
          <li>• 程序输出和错误信息</li>
          <li>• 当前变量的值</li>
          <li>• 调用栈信息</li>
        </ul>
      </div>
    </el-tour-step>

    <el-tour-step target=".terminal-button" title="步骤5: 结束调试" placement="top">
      <div class="tour-content">
        <ul>
          <li>
            <span class="debug-icon iconfont icon-indeterminate"></span>
            <strong>终止调试</strong>：点击终止按钮来结束调试
          </li>
        </ul>
      </div>
    </el-tour-step>

    <el-tour-step target=".debug-guide-trigger" title="步骤6: 重新查看指引" placement="left">
      <div class="tour-content">
        <p>调试指引已完成！如果您需要重新查看指引，可以点击右侧的这个按钮。</p>
      </div>
    </el-tour-step>
  </el-tour>
</template>

<script lang="ts" setup>
  import { ref, onMounted } from 'vue'
  import { QuestionFilled } from '@element-plus/icons-vue'

  // 引导状态
  const tourVisible = ref(false)

  // 按钮位置状态
  const buttonPosition = ref({
    top: 0, // 初始位置，会在onMounted中设置为垂直居中
    right: 20, // 默认距离右边20px
  })

  // 拖拽状态
  const isDragging = ref(false)
  const dragStart = ref({ x: 0, y: 0 })
  const startPosition = ref({ top: 0, right: 0 })
  const hasMoved = ref(false) // 记录是否真的移动了

  // 获取断点区域元素的函数 - 计算左边断点区域位置
  const getBreakpointArea = () => {
    const editorElement = document.querySelector(
      '.coding-workspace__editor .monaco-editor',
    ) as HTMLElement
    if (!editorElement) {
      return document.querySelector('.coding-workspace__editor') || document.body
    }

    // 获取编辑器的位置和尺寸
    const rect = editorElement.getBoundingClientRect()

    // 计算断点区域的位置（左边行号区域）
    // 通常Monaco编辑器的行号区域宽度约为50-60px
    const lineNumberWidth = 60 // 行号区域宽度
    const breakpointAreaWidth = lineNumberWidth // 断点区域宽度

    // 计算断点区域的位置
    const breakpointRect = {
      top: rect.top,
      left: rect.left,
      right: rect.left + breakpointAreaWidth,
      bottom: rect.bottom,
      width: breakpointAreaWidth,
      height: rect.height,
      x: rect.left,
      y: rect.top,
    }

    // 创建一个虚拟的断点区域元素
    const virtualElement = {
      getBoundingClientRect: () => breakpointRect,
      // 添加其他必要的属性
      offsetParent: editorElement.offsetParent,
      offsetTop: editorElement.offsetTop,
      offsetLeft: editorElement.offsetLeft,
      offsetWidth: breakpointAreaWidth,
      offsetHeight: editorElement.offsetHeight,
      clientTop: editorElement.clientTop,
      clientLeft: editorElement.clientLeft,
      clientWidth: breakpointAreaWidth,
      clientHeight: editorElement.clientHeight,
      scrollTop: editorElement.scrollTop,
      scrollLeft: editorElement.scrollLeft,
      scrollWidth: breakpointAreaWidth,
      scrollHeight: editorElement.scrollHeight,
    }

    return virtualElement
  }

  // 处理按钮点击
  const handleButtonClick = (e: MouseEvent) => {
    // 如果正在拖拽或刚刚移动过，不触发指导
    if (isDragging.value || hasMoved.value) {
      e.preventDefault()
      e.stopPropagation()
      return
    }
    showGuide()
  }

  // 显示引导
  const showGuide = () => {
    tourVisible.value = true
  }

  // 开始拖拽
  const startDrag = (e: MouseEvent) => {
    e.preventDefault()
    e.stopPropagation()
    isDragging.value = true
    hasMoved.value = false
    dragStart.value = { x: e.clientX, y: e.clientY }
    startPosition.value = { ...buttonPosition.value }

    document.addEventListener('mousemove', handleDrag)
    document.addEventListener('mouseup', stopDrag)
  }

  // 处理拖拽
  const handleDrag = (e: MouseEvent) => {
    if (!isDragging.value) return

    e.preventDefault()
    e.stopPropagation()

    const deltaX = e.clientX - dragStart.value.x
    const deltaY = e.clientY - dragStart.value.y

    // 检查是否真的移动了（超过5像素才算拖拽）
    const distance = Math.sqrt(deltaX * deltaX + deltaY * deltaY)
    if (distance > 5) {
      hasMoved.value = true
    }

    // 计算新位置
    let newTop = startPosition.value.top + deltaY
    let newRight = startPosition.value.right - deltaX

    // 限制在视窗范围内
    const maxTop = window.innerHeight - 40 // 按钮高度约40px
    const maxRight = window.innerWidth - 40 // 按钮宽度约40px

    newTop = Math.max(0, Math.min(newTop, maxTop))
    newRight = Math.max(0, Math.min(newRight, maxRight))

    buttonPosition.value = {
      top: newTop,
      right: newRight,
    }
  }

  // 停止拖拽
  const stopDrag = () => {
    isDragging.value = false

    // 只有真正移动了才保存位置
    if (hasMoved.value) {
      localStorage.setItem('debug-guide-button-position', JSON.stringify(buttonPosition.value))
    }

    // 延迟重置移动状态，避免点击事件被误判
    setTimeout(() => {
      hasMoved.value = false
    }, 50)

    document.removeEventListener('mousemove', handleDrag)
    document.removeEventListener('mouseup', stopDrag)
  }

  // 处理引导关闭
  const handleTourClose = () => {
    tourVisible.value = false
    localStorage.setItem('debug-guide', JSON.stringify({ firstVisit: false }))
  }

  // 处理引导完成
  const handleTourFinish = () => {
    tourVisible.value = false
    localStorage.setItem('debug-guide', JSON.stringify({ firstVisit: false }))
  }

  // 检查是否需要自动显示引导
  onMounted(() => {
    // 加载保存的按钮位置
    const savedPosition = localStorage.getItem('debug-guide-button-position')
    if (savedPosition) {
      try {
        buttonPosition.value = JSON.parse(savedPosition)
      } catch (error) {
        console.warn('Failed to parse saved button position:', error)
        // 如果加载失败，设置为垂直居中
        buttonPosition.value.top = window.innerHeight * 0.5 - 20
      }
    } else {
      // 如果没有保存的位置，设置为垂直居中
      buttonPosition.value.top = window.innerHeight * 0.5 - 20
    }

    const flag = localStorage.getItem('debug-guide')
    if (!flag) {
      // 延迟显示，确保页面元素已渲染
      setTimeout(() => {
        showGuide()
      }, 1000)
    }
  })
</script>

<style lang="scss" scoped>
  // 自定义引导内容样式
  .tour-content {
    min-width: 280px;
    max-width: 350px;
    padding: 8px 0;

    p {
      margin: 0 0 8px;
      font-size: 14px;
      line-height: 1.6;
      color: $text-color-regular;
    }

    ul {
      padding-left: 0;
      margin: 8px 0 0;

      li {
        display: flex;
        align-items: center;
        margin: 6px 0;
        font-size: 14px;
        line-height: 1.5;
        color: $text-color-regular;
        list-style: none;

        .debug-icon {
          display: inline-block;
          flex-shrink: 0;
          width: 20px;
          height: 20px;
          margin-right: 8px;
          font-size: 14px;
          line-height: 20px;
          color: var(--el-color-primary);
          text-align: center;
          background-color: var(--el-color-primary-light-9);
          border-radius: 4px;

          &.icon-continue1 {
            color: $success-color; // 绿色，对应继续执行
          }

          &.icon-step,
          &.icon-step-in,
          &.icon-step-out {
            color: $info-color; // 蓝色，对应其他调试按钮
          }

          &.icon-indeterminate {
            color: $error-color; // 红色，对应终止调试按钮
          }

          &.icon-play {
            color: $success-color; // 绿色，对应启动调试按钮
          }

          &.icon-caozuo-tiaoshi {
            color: $error-color;
          }
        }
      }
    }

    ul li strong {
      margin-right: 4px;
      color: $text-color-primary;
    }

    // 断点演示样式
    .breakpoint-demo {
      padding: 12px;
      margin: 12px 0;
      background-color: var(--el-color-warning-light-9);
      border: 1px solid var(--el-color-warning-light-7);
      border-radius: 6px;

      .breakpoint-item {
        display: flex;
        align-items: center;
        margin: 6px 0;

        .breakpoint-icon {
          display: inline-block;
          width: 20px;
          height: 20px;
          margin-right: 8px;
          font-size: 16px;
          line-height: 20px;
          text-align: center;
          background-color: var(--el-color-white);
          border-radius: 50%;
          box-shadow: $box-shadow-light;
        }

        .breakpoint-text {
          font-size: 14px;
          font-weight: 500;
          color: $text-color-regular;
        }
      }
    }
  }

  // 调试指引触发按钮样式
  .debug-guide-trigger {
    position: fixed;
    z-index: 1000;
    cursor: move;
    user-select: none;
    transform: translateY(-50%);

    .guide-button {
      color: var(--el-color-white);
      // 使用Element UI的蓝色主题色
      background-color: var(--el-color-primary);
      border-color: var(--el-color-primary);
      box-shadow: 0 2px 8px rgba($primary-color, 0.3);
      transition: all 0.3s ease;

      &:hover {
        background-color: var(--el-color-primary-light-3);
        border-color: var(--el-color-primary-light-3);
        box-shadow: 0 4px 12px rgba($primary-color, 0.4);
        transform: scale(1.05);
      }

      &:active {
        background-color: var(--el-color-primary-dark-2);
        border-color: var(--el-color-primary-dark-2);
        transform: scale(0.98);
      }

      &:focus {
        color: var(--el-color-white);
        background-color: var(--el-color-primary);
        border-color: var(--el-color-primary);
      }
    }
  }
</style>
