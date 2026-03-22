import { requestClient } from '#/api/request';

export function getUserList(params: any) {
  return requestClient.get('/system/user/index', { params });
}

export function getRecycleUserList(params: any) {
  return requestClient.get('/system/user/recycle', { params });
}

export function saveUser(data: any) {
  return requestClient.post('/system/user/save', data);
}

export function updateUser(id: number, data: any) {
  return requestClient.put(`/system/user/update/${id}`, data);
}

export function deleteUser(ids: number[]) {
  return requestClient.delete('/system/user/delete', { data: { ids } });
}

export function realDeleteUser(ids: number[]) {
  return requestClient.delete('/system/user/realDelete', { data: { ids } });
}

export function recoveryUser(ids: number[]) {
  return requestClient.put('/system/user/recovery', { ids });
}

export function changeUserStatus(data: { id: number; status: number }) {
  return requestClient.put('/system/user/changeStatus', data);
}

export function resetPassword(data: { id: number }) {
  return requestClient.put('/system/user/initUserPassword', data);
}

export function clearUserCache(data: { id: number }) {
  return requestClient.post('/system/user/clearCache', data);
}

// Set home page
export function setHomePage(data: { dashboard: string; id: number }) {
  return requestClient.post('/system/user/setHomePage', data);
}
