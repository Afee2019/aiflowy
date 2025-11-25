import type { RouteRecordRaw } from 'vue-router';

const routes: RouteRecordRaw[] = [
  {
    name: 'AssetLibrary',
    path: '/assetLibrary',
    component: () => import('#/views/assetLibrary/index.vue'),
    meta: {
      icon: 'mdi:folder-open-outline',
      order: 3,
      title: '素材库',
    },
  },
];

export default routes;
