import type { RouteRecordRaw } from 'vue-router';

const routes: RouteRecordRaw[] = [
  {
    name: 'Workflow',
    path: '/workflow',
    component: () => import('#/views/bots/index.vue'),
    meta: {
      icon: 'svg:bot',
      order: 1,
      title: '智能体',
    },
  },
  {
    name: 'WorkflowDetail',
    path: '/workflow/:id',
    component: () => import('#/views/bots/bot/index.vue'),
    meta: {
      title: '智能体',
      hideInMenu: true,
      hideInTab: true,
      hideInBreadcrumb: true,
      activePath: '/workflow',
    },
  },
];

export default routes;
