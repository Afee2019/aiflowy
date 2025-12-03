import type { RouteRecordRaw } from 'vue-router';

import { $t } from '#/locales';

const routes: RouteRecordRaw[] = [
  {
    meta: {
      icon: 'ant-design:apartment-outlined',
      title: $t('datacenterTable.title'),
      hideInMenu: true,
      hideInTab: true,
      hideInBreadcrumb: true,
    },
    name: 'WorkflowDesign',
    path: '/workflow/design',
    component: () => import('#/views/ai/workflow/WorkflowDesign.vue'),
  },
  {
    meta: {
      title: '运行',
      openInNewWindow: true,
      noBasicLayout: true,
      hideInMenu: true,
    },
    name: 'RunPage',
    path: '/workflow/run',
    component: () => import('#/views/ai/workflow/RunPage.vue'),
  },
];

export default routes;
