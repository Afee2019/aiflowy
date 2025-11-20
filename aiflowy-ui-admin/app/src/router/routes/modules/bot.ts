import type { RouteRecordRaw } from 'vue-router';

const routes: RouteRecordRaw[] = [
  {
    name: 'BotRun',
    path: '/ai/bots/run/:id',
    component: () => import('#/views/ai/bots/run.vue'),
    meta: {
      title: 'bot',
      noBasicLayout: true,
      openInNewWindow: true,
    },
  },
];

export default routes;
