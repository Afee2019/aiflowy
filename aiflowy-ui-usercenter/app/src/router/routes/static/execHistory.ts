import type { RouteRecordRaw } from 'vue-router';

const routes: RouteRecordRaw[] = [
  {
    name: 'ExecHistory',
    path: '/execHistory',
    component: () => import('#/views/execHistory/index.vue'),
    meta: {
      icon: 'mdi:history',
      order: 2,
      title: '执行记录',
    },
  },
];

export default routes;
