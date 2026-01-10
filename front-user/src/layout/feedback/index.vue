<template>
  <div
    :style="{ right: `${right}px`, bottom: `${bottom}px` }"
    :class="{ feedback: true, feedback_move: moved }"
    @click="handleClick"
    @mousedown="startDrag"
    @touchstart="startTouchDrag"
  >
    <p class="iconfont icon-feedback feedback-icon"></p>
    <el-text class="feedback-text"> 反馈 </el-text>
  </div>
</template>

<script setup>
  import { ref, defineProps, onMounted, onBeforeUnmount } from 'vue'

  // 图标的高度和宽度
  const width = 16
  const height = 40

  // 屏幕高度和宽度
  let windowWidth = window.innerWidth
  let windowHeight = window.innerHeight

  const right = ref(windowWidth - width - 10) // 初始距右边距离
  const bottom = ref(300) // 初始距底部距离

  const isDragging = ref(false)
  let offsetX, offsetY, moved

  const startDrag = (event) => {
    disableTextSelection()
    isDragging.value = true
    let startX = event.clientX
    let startY = event.clientY
    // 根据悬浮球的当前位置重新计算 offsetX 和 offsetY
    offsetX = window.innerWidth - right.value - 16 - startX
    offsetY = window.innerHeight - bottom.value - 16 - startY
    moved = false

    window.addEventListener('mousemove', onMouseMove)
    window.addEventListener('mouseup', stopDrag)
  }

  const startTouchDrag = (event) => {
    disableTextSelection()
    event.preventDefault()
    const touch = event.touches[0]
    isDragging.value = true
    let startX = touch.clientX
    let startY = touch.clientY
    // 同样重新计算 offsetX 和 offsetY
    offsetX = window.innerWidth - right.value - 16 - startX
    offsetY = window.innerHeight - bottom.value - 16 - startY
    moved = false

    window.addEventListener('touchmove', onTouchMove, { passive: false })
    window.addEventListener('touchend', stopTouchDrag)
  }

  const onMouseMove = (event) => {
    if (!isDragging.value) return
    updatePosition(event.clientX, event.clientY)
  }

  const onTouchMove = (event) => {
    event.preventDefault()
    if (!isDragging.value) return
    const touch = event.touches[0]
    updatePosition(touch.clientX, touch.clientY)
  }

  const updatePosition = (x, y) => {
    moved = true
    let newRight = Math.max(
      0,
      Math.min(window.innerWidth - width, window.innerWidth - x - offsetX - width / 2),
    )
    newRight = Math.min(newRight, window.innerWidth - width - 10)
    const newBottom = Math.max(
      0,
      Math.min(window.innerHeight - height, window.innerHeight - y - offsetY - height / 2),
    )

    // 直接更新值，避免重复计算
    right.value = newRight
    bottom.value = newBottom
  }

  const stopDrag = () => {
    enableTextSelection()
    isDragging.value = false
    window.removeEventListener('mousemove', onMouseMove)
    window.removeEventListener('mouseup', stopDrag)
  }

  const stopTouchDrag = () => {
    enableTextSelection()
    isDragging.value = false
    window.removeEventListener('touchmove', onTouchMove)
    window.removeEventListener('touchend', stopTouchDrag)
  }

  const handleResize = () => {
    // 当屏幕大小改变时，更新偏移量
    if (right.value == windowWidth - width) {
      right.value = window.innerWidth - width
    } else {
      right.value = right.value * ((window.innerWidth - width) / (windowWidth - width))
    }
    if (bottom.value == windowHeight - height) {
      bottom.value = window.innerHeight - height
    } else {
      bottom.value = bottom.value * ((window.innerHeight - height) / (windowHeight - height))
    }
    // 按照比例扩大或缩小
    windowWidth = window.innerWidth
    windowHeight = window.innerHeight
  }

  function disableTextSelection() {
    document.body.style.userSelect = 'none'
    document.body.style.webkitUserSelect = 'none' // 针对 Safari 浏览器
    document.body.style.msUserSelect = 'none' // 针对 IE 浏览器
    document.body.style.mozUserSelect = 'none' // 针对 Firefox 浏览器
  }

  function enableTextSelection() {
    document.body.style.userSelect = ''
    document.body.style.webkitUserSelect = ''
    document.body.style.msUserSelect = ''
    document.body.style.mozUserSelect = ''
  }

  const handleClick = () => {
    // 处理点击逻辑
    if (!moved) {
      window.open('https://wj.qq.com/s2/17342588/c8c0/')
    }
  }

  onMounted(() => {
    window.addEventListener('resize', handleResize)
  })

  onBeforeUnmount(() => {
    window.removeEventListener('resize', handleResize)
  })
</script>

<style scoped lang="scss">
  .feedback {
    position: absolute;
    z-index: 99999;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    width: 24px;
    height: 40px;
    cursor: pointer;
    background: rgba($base-background-color, 0.9);
    border-radius: 10px;
    box-shadow: $box-shadow-light;

    &:hover {
      background: $base-background-color;
      box-shadow: $box-shadow;
      transform: scale(1.05);
    }

    .feedback-icon {
      width: 20px;
      height: 20px;
      margin-bottom: 2px;
      font-size: 20px;
      color: $primary-color;
    }

    .feedback-text {
      margin: 0;
      font-size: 10px;
      line-height: 1;
      color: $primary-color;
    }
  }
</style>
