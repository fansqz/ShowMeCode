// 进行axios二次封装
import axios from 'axios'
import { ElMessage } from 'element-plus'

import router from '@/router' // 请根据实际路径调整
import useUserStore from '@/store/modules/user'

const request = axios.create({
  baseURL: import.meta.env.VITE_SERVE + import.meta.env.VITE_APP_BASE_API,
  // 请求超时时间
  timeout: 50000,
})

// 添加请求和响应拦截器
request.interceptors.request.use((config) => {
  // 设置token
  const userStore = useUserStore()
  config.headers.token = userStore.token
  config.headers.visitorUid = generateVisitorId()
  // 返回配置对象
  return config
})

// 获取访客ID
function generateVisitorId(): string {
  let visitorId = localStorage.getItem('visitor-id')
  if (visitorId) {
    return visitorId
  }
  visitorId =
    Math.random().toString(36).substring(2, 15) + Math.random().toString(36).substring(2, 15)
  localStorage.setItem('visitor-id', visitorId)
  return visitorId
}

// 配置响应拦截器
request.interceptors.response.use(
  (response) => {
    const userStore = useUserStore()
    // 如果response中携带token，则更新token
    const token = response.headers['token']
    if (token) {
      userStore.token = token
    }
    // 简化数据
    return response.data
  },
  (error) => {
    const userStore = useUserStore()
    // 处理网络错误
    let message = ''
    const status = error.response.status
    switch (status) {
      case 401:
        userStore.userLogout()
        router?.push('/login')
        return Promise.reject(error)
      case 403:
        message = '无权访问'
        break
      case 404:
        message = '请求地址错误'
        break
      case 500:
        message = '服务器出现问题'
        break
      default:
        message = '网络出现问题'
        break
    }
    ElMessage({
      type: 'error',
      message,
    })
    return Promise.reject(error)
  },
)

export default request
