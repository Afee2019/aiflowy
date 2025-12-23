import type { RouteRecordRaw } from 'vue-router';

const routes: RouteRecordRaw[] = [
  {
    name: 'ExecHistory',
    path: '/execHistory',
    component: () => import('#/views/execHistory/index.vue'),
    meta: {
      icon: 'svg:exec-history',
      order: 2,
      title: '执行记录',
    },
  },
  {
    name: 'ExecHistoryDetails',
    path: '/execHistory/:id',
    component: () => import('#/views/execHistory/details/index.vue'),
    meta: {
      title: '执行记录',
      hideInMenu: true,
      hideInTab: true,
      hideInBreadcrumb: true,
      activePath: '/execHistory',
    },
  },
];

export default routes;
