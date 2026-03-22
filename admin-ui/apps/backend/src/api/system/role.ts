import { requestClient } from '#/api/request';

export function getRoleList(params?: any) {
  return requestClient.get('/system/role/list', { params });
}
