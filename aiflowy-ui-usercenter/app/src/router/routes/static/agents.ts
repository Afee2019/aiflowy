import type { RouteRecordRaw } from 'vue-router';

const routes: RouteRecordRaw[] = [
  {
    name: 'Agents',
    path: '/agents',
    component: () => import('#/views/agents/index.vue'),
    meta: {
      icon: 'mdi:robot-outline',
      order: 1,
      title: '智能体',
    },
  },
];

export default routes;
