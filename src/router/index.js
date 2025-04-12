import { createRouter, createWebHistory } from 'vue-router';
import TableComponent from '../components/TableComponent.vue';
import Creator from '@/components/operation/Creator.vue';

const routes = [
  {
    path: '/',
    redirect: '/setup'
  },
  {
    path: '/setup',
    name: 'Setup',
    component: () => import('../components/SetupPage.vue'),
    meta: { requiresAuth: false } // SetupPage ����Ҫ��֤
  },
  {
    path: '/main',
    name: 'MainPage',
    component: () => import('../components/MainPage.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('../components/LoginPage.vue')
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('../components/Register.vue')
  },
  {
    path: '/table',
    name: 'Table',
    component: () => import('../components/TableComponent.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/',
    name: 'Table',
    component: TableComponent
  },
  {
    path: '/create',
    name: 'Create',
    component: Creator
  }
];

const router = createRouter({
  history: createWebHistory(),
  routes
});

// ·������
router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token');

  if (to.meta.requiresAuth && !token) {
    // �����Ҫ��֤��û��token���ض���SetupPage
    next('/setup');
  } else if ((to.path === '/login' || to.path === '/register') && token) {
    // ����ѵ�¼�ҷ��ʵ�¼��ע��ҳ���ض�����ҳ
    next('/setup');
  } else {
    next();
  }
});

export default router;