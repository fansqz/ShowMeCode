// 路由鉴权
import router from '@/router';
import setting from './setting';
// 进度条
import nprogress from 'nprogress';
import 'nprogress/nprogress.css';
nprogress.configure({ showSpinner: false });
import useUserStore from './store/modules/user';
import pinia from './store';
const userStore = useUserStore(pinia);

// 全局前置守卫
router.beforeEach(async (to: any, _from: any, next: any) => {
  // 顶部标题
  document.title = setting.title + '-' + to.meta.title;
  // 进度条
  nprogress.start();
  // 判断用户是否登录
  const token = userStore.token;

  // 如果是login
  if (to.path == '/login') {
    if (token) {
      next({ path: '/home' });
    } else {
      next();
    }
    return;
  }

  // 没有token则路由到login
  if (!token) {
    next({ path: '/login', query: { redirect: to.path } });
    return;
  }

  // 获取用户名称，如果获取不到用户名称则进行获取
  try {
    await userStore.userInfo();
    // 等到用户信息加载完毕，异步路由注册以后再放行
    next();
    return;
  } catch (error) {
    // 意味着token过期
    userStore.userLogout();
    next({ path: '/login', query: { redirect: to.path } });
  }
});

// 全局后置守卫
router.afterEach(() => {
  nprogress.done();
});
