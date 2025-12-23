import type { RouteRecordRaw } from 'vue-router';

const routes: RouteRecordRaw[] = [
  {
    name: 'ChatAssistant',
    path: '/chatAssistant',
    component: () => import('#/views/chatAssistant/index.vue'),
    meta: {
      icon: 'svg:chat-assistant',
      order: 0,
      title: '聊天助理',
    },
  },
];

export default routes;
