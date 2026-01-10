// 路由鉴权
import nprogress from 'nprogress'

import router from '@/router'

// 进度条
import 'nprogress/nprogress.css'
nprogress.configure({ showSpinner: false })
import pinia from './store'
import useUserStore from './store/modules/user'
const userStore = useUserStore(pinia)

// 允许不使用登陆态访问的路径
const allowPaths = ['/register', '/login', '/coding', '/home', '/learn']

// 全局前置守卫
router.beforeEach(async (to: any, _from: any, next: any) => {
  // 进度条
  nprogress.start()

  // 校验token是否正确，如果正确则继续访问
  const token = userStore.token
  if (token) {
    try {
      await userStore.userInfo()
      next()
      return
    } catch (error) {
      userStore.userLogout()
    }
  }

  // 校验是否在允许路径
  if (checkIsAllow(to.path)) {
    next()
    return
  }

  // 其余路径没有登录态，路由到login页
  if (!token) {
    next({ path: '/login', query: { redirect: to.path } })
    return
  }
})

// 全局后置守卫
router.afterEach(() => {
  nprogress.done()
})

const checkIsAllow = (path: string): boolean => {
  for (const item of allowPaths) {
    if (path.startsWith(item)) {
      return true
    }
  }
  return false
}
