/* eslint-disable @typescript-eslint/ban-ts-comment */
import { defineStore } from 'pinia'

import { reqAccountInfo } from '@/api/account'
import { reqLogin } from '@/api/auth'
import { reqGetURL } from '@/api/common'

const useUserStore = defineStore('User', {
  state: (): any => {
    return {
      token: localStorage.getItem('TOKEN'),
      avatar: '',
      email: '',
      phone: '',
      username: '',
      loginName: '',
    }
  },
  actions: {
    async userLogin(data: any) {
      const result = await reqLogin(data)
      if (result.code == 200) {
        this.token = result.data
        return 'ok'
      } else {
        return Promise.reject(new Error(result.message))
      }
    },
    async userInfo() {
      const result = await reqAccountInfo()
      if (result.code == 200) {
        const avatarPath = result.data.avatar
        // 读取头像
        if (avatarPath != '') {
          const result2 = await reqGetURL(avatarPath)
          if (result2.code == 200) {
            this.avatar = result2.data
          }
        }
        this.username = result.data.username
        this.email = result.data.email
        this.phone = result.data.phone
        this.loginName = result.data.loginName
      } else {
        return Promise.reject(new Error(result.message))
      }
    },
    userLogout() {
      localStorage.removeItem('TOKEN')
      this.token = ''
      this.username = ''
      this.email = ''
      this.phone = ''
      this.loginName = ''
      this.avatar = ''
    },
    // 判断用户是否登录
    isLoggedIn() {
      return !!this.token
    },
  },
  persist: true,
})

export default useUserStore
