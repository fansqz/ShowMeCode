import { RouteRecordRaw } from 'vue-router';

// 常量路由
export const constantRoute: Readonly<RouteRecordRaw[]> = [
  {
    path: '/login',
    component: () => import('@/views/login/index.vue'),
    name: 'login',
    meta: {
      hidden: false,
    },
  },
  {
    path: '/404',
    component: () => import('@/views/404/index.vue'),
    name: '404',
    meta: {
      hidden: false,
    },
  },
  {
    path: '/',
    component: () => import('@/layout/index.vue'),
    name: 'layout',
    redirect: '/home',
    meta: {
      hidden: false,
    },
    children: [
      // 首页
      {
        path: '/home',
        name: 'home',
        meta: {
          title: '首页',
          icon: 'Monitor',
          hidden: false,
        },
        component: () => import('@/views/home/index.vue'),
      },
    ],
  },
];

// 异步路由
export const asyncRoute: Readonly<RouteRecordRaw[]> = [
  {
    path: '/manage',
    component: () => import('@/layout/index.vue'),
    name: 'manage',
    meta: {
      hidden: false,
    },
    children: [
      {
        path: '/manage/visual-document/bank',
        name: 'visual-document-bank',
        meta: {
          title: '可视化教程',
          icon: 'Reading',
          hidden: false,
        },
        component: () => import('@/views/visual-document/bank/index.vue'),
      },
      {
        path: '/manage/visual-document/:bankID',
        name: 'visual-document',
        meta: {
          title: '可视化教程',
          hidden: true,
        },
        component: () => import('@/views/visual-document/document/index.vue'),
      },
      // 权限管理
      {
        path: '/manage/permissions',
        name: 'permissions',
        meta: {
          title: '权限管理',
          icon: 'Lock',
          hidden: false,
        },
        children: [
          {
            path: '/manage/permissions/api',
            name: 'api',
            meta: {
              title: '接口管理',
              icon: 'Connection',
              hidden: false,
            },
            component: () => import('@/views/permissions/api/index.vue'),
          },
          {
            path: '/manage/permissions/menu',
            name: 'menu',
            meta: {
              title: '菜单管理',
              icon: 'PriceTag',
              hidden: false,
            },
            component: () => import('@/views/permissions/menu/index.vue'),
          },
          {
            path: '/manage/permissions/role',
            name: 'role',
            meta: {
              title: '角色管理',
              icon: 'User',
              hidden: false,
            },
            component: () => import('@/views/permissions/role/index.vue'),
          },
          {
            path: '/manage/permissions/user',
            name: 'user',
            meta: {
              title: '用户管理',
              icon: 'Avatar',
              hidden: false,
            },
            component: () => import('@/views/permissions/user/index.vue'),
          },
        ],
      },
    ],
  },
];

// 任意路由
export const anyRoute = [
  {
    path: '/:pathMath(.*)*',
    redirect: '/404',
    name: 'Any',
    meta: {
      hidden: false,
    },
  },
];
