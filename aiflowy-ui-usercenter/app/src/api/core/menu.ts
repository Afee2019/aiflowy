import type { RouteRecordStringComponent } from '@aiflowy/types';

import { requestClient } from '#/api/request';

/**
 * 获取用户所有菜单
 */
export async function getAllMenusApi() {
  return requestClient.get<RouteRecordStringComponent[]>(
    '/api/v1/sysMenu/treeV2',
  );
}
