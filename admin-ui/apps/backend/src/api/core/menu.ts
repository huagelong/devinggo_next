import type { RouteRecordStringComponent } from '@vben/types';

import { getSystemInfoApi } from './user';

/**
 * 获取用户所有菜单路由（从 /system/getInfo 的 routers 字段提取）
 */
export async function getAllMenusApi() {
  const info = await getSystemInfoApi();
  return info.routers as RouteRecordStringComponent[];
}
