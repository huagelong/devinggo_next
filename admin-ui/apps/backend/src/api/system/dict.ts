import { requestClient } from '#/api/request';

export function getDictList(code: string) {
  return requestClient.get('/system/dataDict/list', { params: { code } });
}
