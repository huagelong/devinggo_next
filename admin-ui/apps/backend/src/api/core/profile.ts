import file2md5 from 'file2md5';

import { requestClient } from '#/api/request';

type UploadRequestData = Record<string, boolean | number | string>;

async function fileToMd5(file: File) {
  return file2md5(file);
}

export async function buildImageUploadFormData(
  file: File,
  requestData: UploadRequestData = {},
) {
  const hash = await fileToMd5(file);
  const formData = new FormData();

  formData.append('image', file);
  formData.append('isChunk', 'false');
  formData.append('hash', hash);

  Object.entries(requestData).forEach(([name, value]) => {
    formData.append(name, String(value));
  });

  return formData;
}

export async function updateUserInfoApi(data: any) {
  return requestClient.post('/system/user/updateInfo', data);
}

export async function modifyPasswordApi(data: any) {
  return requestClient.post('/system/user/modifyPassword', data);
}

export async function uploadImageApi(data: FormData) {
  return requestClient.post('/system/uploadImage', data);
}

export async function uploadImageFileApi(
  file: File,
  requestData: UploadRequestData = {},
) {
  const formData = await buildImageUploadFormData(file, requestData);
  return uploadImageApi(formData);
}

export async function getLoginLogListApi(params: any) {
  return requestClient.get('/system/common/getLoginLogList', { params });
}

export async function getOperationLogListApi(params: any) {
  return requestClient.get('/system/common/getOperationLogList', { params });
}
