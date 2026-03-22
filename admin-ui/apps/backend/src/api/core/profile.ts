import { requestClient } from '#/api/request';

/**
 * 个人资料接口
 */

/**
 * 更新个人信息
 */
export async function updateUserInfoApi(data: any) {
  return requestClient.post('/system/user/updateInfo', data);
}

/**
 * 修改密码
 */
export async function modifyPasswordApi(data: any) {
  return requestClient.post('/system/user/modifyPassword', data);
}

/**
 * 上传头像 (图片)
 */
export async function uploadImageApi(data: FormData) {
  return requestClient.post('/system/uploadImage', data, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  });
}

/**
 * 获取登录日志
 */
export async function getLoginLogListApi(params: any) {
  return requestClient.get('/system/common/getLoginLogList', { params });
}

/**
 * 获取操作日志
 */
export async function getOperationLogListApi(params: any) {
  return requestClient.get('/system/common/getOperationLogList', { params });
}
