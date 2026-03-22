import { requestClient } from '#/api/request';

export function getDeptTree(params?: any) {
  return requestClient.get('/system/dept/tree', { params });
}
