import { requestClient } from '#/api/request';

export function getPostList(params?: any) {
  return requestClient.get('/system/post/list', { params });
}
