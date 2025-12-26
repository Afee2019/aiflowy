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
    path: '/ai/workflow/design',
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
    path: '/ai/workflow/run',
    component: () => import('#/views/ai/workflow/RunPage.vue'),
  },
  {
    meta: {
      title: $t('aiWorkflowExecRecord.moduleName'),
      hideInMenu: true,
      hideInTab: true,
      hideInBreadcrumb: true,
    },
    name: 'ExecRecord',
    path: '/ai/workflow/executeRecords',
    component: () =>
      import('#/views/ai/workflow/execute/WorkflowExecResultList.vue'),
  },
  {
    meta: {
      title: $t('aiWorkflowRecordStep.moduleName'),
      hideInMenu: true,
      hideInTab: true,
      hideInBreadcrumb: true,
    },
    name: 'RecordStep',
    path: '/ai/workflow/executeSteps',
    component: () =>
      import('#/views/ai/workflow/execute/WorkflowExecStepList.vue'),
  },
];

export default routes;
