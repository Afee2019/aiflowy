import type { RouteRecordRaw } from 'vue-router';

const routes: RouteRecordRaw[] = [
  {
    name: 'PersonalCenter',
    path: '/personalCenter',
    component: () => import('#/views/personalCenter/index.vue'),
    meta: {
      icon: 'mdi:account-outline',
      order: 5,
      title: '个人中心',
      hideInMenu: true,
    },
  },
];

export default routes;
