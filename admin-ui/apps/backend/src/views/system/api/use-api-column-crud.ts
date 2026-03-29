import { ref } from 'vue';

import type { ApiColumnListItem, ApiColumnType } from './model';
import type { ApiColumnApi } from '#/api/system/api-column';

import {
  getApiColumnPageList,
  getRecycleApiColumnList,
} from '#/api/system/api-column';
import { useCrudPage } from '#/composables/crud/use-crud-page';

import { createApiColumnSearchForm } from './schemas';

const emptyResponse: ApiColumnApi.ListResponse = {
  items: [],
  pageInfo: { total: 0 },
};

export function useApiColumnCrud() {
  const currentApiId = ref<number | null>(null);
  const currentType = ref<ApiColumnType>(1);

  const crud = useCrudPage<
    ApiColumnListItem,
    ReturnType<typeof createApiColumnSearchForm>,
    ApiColumnApi.ListResponse
  >({
    defaultSearchForm: createApiColumnSearchForm,
    fetchList: (params, context) => {
      if (!currentApiId.value) {
        return Promise.resolve(emptyResponse);
      }
      return context.isRecycleBin
        ? getRecycleApiColumnList(params)
        : getApiColumnPageList(params);
    },
    buildParams: (form) => {
      const params: ApiColumnApi.ListQuery = {
        api_id: currentApiId.value ?? undefined,
        type: currentType.value,
      };
      if (form.name) params.name = form.name;
      if (form.data_type) params.data_type = form.data_type;
      if (typeof form.status === 'number') params.status = form.status;
      if (typeof form.is_required === 'number') params.is_required = form.is_required;
      if (form.created_at?.length === 2 && form.created_at[0]) {
        params.created_at = form.created_at;
      }
      return params as Record<string, unknown>;
    },
  });

  function setContext(apiId: number, type: ApiColumnType) {
    currentApiId.value = apiId;
    currentType.value = type;
  }

  return {
    ...crud,
    currentApiId,
    currentType,
    setContext,
  };
}
