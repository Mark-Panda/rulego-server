import { createRouter, createWebHashHistory } from 'vue-router';
import { SESSIONSTORAGE_KEYS, getSession } from '@src/utils/sessionstorage';

const routes = [
  {
    path: '/',
    redirect: '/workflow-list',
    component: () => import('@src/layout/default-layout.vue'),
    children: [
      {
        path: '/workflow-list',
        name: 'workflow-list',
        component: () => import('@src/pages/workflow-list/workflow-list.vue'),
      },
      {
        path: '/runtime-logs',
        name: 'runtime-logs',
        component: () => import('@src/pages/runtime-logs/runtime-logs.vue'),
      },
      {
        path: '/share-node-list',
        name: 'share-node-list',
        component: () =>
          import('@src/pages/share-node-list/share-node-list.vue'),
      },
      {
        path: '/components-installed',
        name: 'components-installed',
        component: () =>
          import('@src/pages/components-management/components-installed.vue'),
      },
      {
        path: '/components-market',
        name: 'components-market',
        component: () =>
          import('@src/pages/components-management/components-market.vue'),
      },
      {
        path: '/components-rules',
        name: 'components-rules',
        component: () =>
          import('@src/pages/components-management/components-rules.vue'),
      },
    ],
  },
  {
    path: '/workflow',
    name: 'workflow',
    component: () => import('@src/pages/workflow/workflow.vue'),
  },
  {
    path: '/login',
    name: 'login',
    component: () => import('@src/pages/login/login.vue'),
  },
];

const router = createRouter({
  history: createWebHashHistory(),
  routes,
});

router.beforeEach((to, from, next) => {
  const token = getSession(SESSIONSTORAGE_KEYS.TOKEN);
  
  // 如果用户已登录
  if (token) {
    if (to.path === '/login') {
      // 已登录用户访问登录页，重定向到首页
      next('/');
    } else {
      next();
    }
  } else {
    // 用户未登录
    if (to.path !== '/login') {
      // 未登录用户访问非登录页，重定向到登录页
      next('/login');
    } else {
      next();
    }
  }
});

export default router;
