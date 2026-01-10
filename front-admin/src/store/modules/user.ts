/* eslint-disable @typescript-eslint/ban-ts-comment */
import { defineStore } from 'pinia';
import { reqLogin } from '@/api/auth';
import { reqAccountInfo } from '@/api/account';
// 引入路由
import { constantRoute, asyncRoute, anyRoute } from '@/router/routers';
//@ts-ignore
import cloneDeep from 'lodash/cloneDeep';
import { reqGetURL } from '@/api/common';

const useUserStore = defineStore('User', {
  state: (): any => {
    return {
      token: localStorage.getItem('TOKEN'),
      menuRoutes: [],
      email: '',
      phone: '',
      sex: '',
      username: '',
      loginName: '',
    };
  },
  actions: {
    async userLogin(data: any) {
      const result = await reqLogin(data);
      if (result.code == 200) {
        this.token = result.data;
        // 本地存一份
        localStorage.setItem('TOKEN', result.data);
        return 'ok';
      } else {
        return Promise.reject(new Error(result.message));
      }
    },
    async userInfo() {
      const result = await reqAccountInfo();
      if (result.code == 200) {
        let avatarPath = result.data.avatar;
        // 读取头像
        if (avatarPath != '') {
          let result2 = await reqGetURL(avatarPath);
          if (result2.code == 200) {
            this.avatar = result2.data;
          }
        }
        this.username = result.data.username;
        this.email = result.data.email;
        this.sex = result.data.email;
        this.phone = result.data.phone;
        this.loginName = result.data.loginName;
        const userAsyncRoute = filterAsyncRoute(cloneDeep(asyncRoute), result.data.menus);
        this.menuRoutes = [...constantRoute, ...userAsyncRoute, ...anyRoute];
      } else {
        return Promise.reject(new Error(result.message));
      }
    },
    userLogout() {
      this.token = '';
      this.username = '';
      this.email = '';
      this.sex = '';
      this.phone = '';
      this.loginName = '';
      localStorage.removeItem('TOKEN');
    },
  },
  getters: {},
});

function filterAsyncRoute(asyncRoute: any[], routeList: string[]) {
  return asyncRoute.filter((item: any) => {
    if (routeList.includes(item.name)) {
      if (item.children && item.children.length > 0) {
        item.children = filterAsyncRoute(item.children, routeList);
      }
      return true;
    }
  });
}

export default useUserStore;
