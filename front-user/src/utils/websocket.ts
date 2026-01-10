import useUserStore from '@/store/modules/user'

import { EventDispatcher, EventListenerCallback } from './event-dispatcher'

const dispatcher = new EventDispatcher<string>()
let ws: WebSocket | null = null
let heartBeat: any
const user = useUserStore()

// 主动断开websocket
function closeWebSocket() {
  if (ws) {
    clearInterval(heartBeat)
    ws.close()
    ws = null
    heartBeat = undefined
  }
}

/**
 * 初始化websocket
 */
function initWebSocket(fun: () => any) {
  if (ws) {
    return
  } // 连接存在，退出初始化操作
  ws = new WebSocket(`${import.meta.env.VITE_WS_SERVE}/${user.token}`)

  // 消息响应
  ws.onmessage = (e) => {
    dispatcher.dispatch('onmessage', e.data)
  }

  // 连接关闭
  ws.onclose = (e) => {
    console.log(`WS关闭 ${new Date()}`, e)
    ws = null
    dispatcher.dispatch('onclose', 'websocket close')
  }

  // 连接打开
  ws.onopen = (e) => {
    fun()
    console.log(`WS开启 ${new Date()}`, e)
    dispatcher.dispatch('onopent', 'websocket open')
    // ws成功连接时，将当前的连接模式设置为ws
    heartBeat = window.setInterval(() => {
      // 心跳包做心跳检测
      sendWebSocket({ type: 'ping', content: null })
      console.log('xdb-heartbeat', new Date())
    }, 40000)
  }

  // 连接异常
  ws.onerror = (e) => {
    console.log(`WS发生错误 ${new Date()}`, e)
    dispatcher.dispatch('onerror', 'websocket error')
    clearInterval(heartBeat)
    heartBeat = undefined
    ws = null
  }
}

/**
 * 发送消息到websocket
 */
function sendWebSocket(data: any) {
  if (ws?.readyState === WebSocket.OPEN) {
    ws.send(JSON.stringify(data))
  } else {
    initWebSocket(() => {
      ws?.send(JSON.stringify(data))
    })
  }
}

// 注册websocket事件的监听
function onWebSocketEvent(event: string, callback: EventListenerCallback) {
  dispatcher.on(event, callback)
}

// 将方法暴露出去
export { initWebSocket, sendWebSocket, closeWebSocket, onWebSocketEvent }
