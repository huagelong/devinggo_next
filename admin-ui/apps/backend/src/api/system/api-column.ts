import { requestClient } from '#/api/request';
import type { BatchIdsPayload, IdType } from '#/types/common';
import type { PageQuery, PageResponse } from '#/types/paging';

export namespace ApiColumnApi {
  export interface ListItem {
    api_id: number;
    created_at?: string;
    data_type: string | number;
    default_value?: string;
    description?: string;
    id: number;
    is_required: number;
    name: string;
    remark?: string;
    status: number;
    type: number;
  }

  export interface ListQuery extends Partial<PageQuery> {
    api_id?: IdType;
    created_at?: string[];
    data_type?: string | number;
    is_required?: number;
    name?: string;
    status?: number;
    type?: number;
  }

  export interface SubmitPayload {
    api_id: number;
    data_type: string | number;
    default_value?: string;
    description?: string;
    id?: number;
    is_required: number;
    name: string;
    remark?: string;
    status: number;
    type: number;
  }

  export interface ChangeStatusPayload {
    id: number;
    status: number;
  }

  export interface DownloadResponse {
    data: Blob;
    headers?:
      | Record<string, string | undefined>
      | {
          get?: (name: string) => null | string;
          [key: string]: unknown;
        };
  }

  export type ListResponse = PageResponse<ListItem>;
  export type BatchPayload = BatchIdsPayload<number>;
}

export function getApiColumnPageList(params: ApiColumnApi.ListQuery) {
  return requestClient.get<ApiColumnApi.ListResponse>('/system/apiColumn/index', {
    params,
  });
}

export function getRecycleApiColumnList(params: ApiColumnApi.ListQuery) {
  return requestClient.get<ApiColumnApi.ListResponse>('/system/apiColumn/recycle', {
    params,
  });
}

export function getApiColumnDetail(id: number) {
  return requestClient.post<ApiColumnApi.ListItem>(`/system/apiColumn/read/${id}`);
}

export function saveApiColumn(data: ApiColumnApi.SubmitPayload) {
  return requestClient.post<void>('/system/apiColumn/save', data);
}

export function updateApiColumn(id: number, data: ApiColumnApi.SubmitPayload) {
  return requestClient.put<void>(`/system/apiColumn/update/${id}`, data);
}

export function deleteApiColumn(ids: number[]) {
  return requestClient.delete<void>('/system/apiColumn/delete', { data: { ids } });
}

export function realDeleteApiColumn(ids: number[]) {
  return requestClient.delete<void>('/system/apiColumn/realDelete', {
    data: { ids },
  });
}

export function recoveryApiColumn(ids: number[]) {
  return requestClient.put<void>('/system/apiColumn/recovery', { ids });
}

export function changeApiColumnStatus(data: ApiColumnApi.ChangeStatusPayload) {
  return requestClient.put<void>('/system/apiColumn/changeStatus', data);
}

export function importApiColumnFile(
  file: File,
  payload: { api_id: number; type: number },
) {
  const formData = new FormData();
  formData.append('file', file);
  formData.append('api_id', String(payload.api_id));
  formData.append('type', String(payload.type));
  return requestClient.post<void>('/system/apiColumn/import', formData);
}

export function exportApiColumnList(data: Omit<ApiColumnApi.ListQuery, keyof PageQuery>) {
  return requestClient.download<ApiColumnApi.DownloadResponse>('/system/apiColumn/export', {
    data,
    method: 'POST',
    responseReturn: 'raw',
  });
}

export function downloadApiColumnTemplate() {
  return requestClient.download<ApiColumnApi.DownloadResponse>(
    '/system/apiColumn/downloadTemplate',
    {
      method: 'GET',
      responseReturn: 'raw',
    },
  );
}
