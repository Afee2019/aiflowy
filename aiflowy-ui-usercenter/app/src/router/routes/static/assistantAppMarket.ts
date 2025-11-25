import type { RouteRecordRaw } from 'vue-router';

const routes: RouteRecordRaw[] = [
  {
    name: 'AssistantAppMarket',
    path: '/assistantAppMarket',
    component: () => import('#/views/assistantAppMarket/index.vue'),
    meta: {
      icon: 'mdi:storefront-outline',
      order: 4,
      title: '助理应用市场',
    },
  },
];

export default routes;
