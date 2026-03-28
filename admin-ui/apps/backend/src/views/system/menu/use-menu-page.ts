import type { MenuApi } from '#/api/system/menu';

import { reactive, ref } from 'vue';

import { message } from '#/adapter/tdesign';
import { getMenuTreeList, getRecycleMenuTreeList } from '#/api/system/menu';

import { createMenuSearchForm } from './schemas';

export function useMenuPage() {
  const searchForm = reactive(createMenuSearchForm());
  const tableData = ref<MenuApi.TreeItem[]>([]);
  const loading = ref(false);
  const isRecycleBin = ref(false);
  const selectedRowKeys = ref<Array<number | string>>([]);

  function buildParams() {
    const params: Partial<MenuApi.ListQuery> = {};
    if (searchForm.name) params.name = searchForm.name;
    if (searchForm.code) params.code = searchForm.code;
    if (searchForm.level) params.level = searchForm.level;
    if (searchForm.status !== undefined) params.status = searchForm.status;
    if (searchForm.created_at?.length === 2 && searchForm.created_at[0]) {
      params.created_at = searchForm.created_at;
    }
    return params;
  }

  async function fetchTableData() {
    loading.value = true;
    try {
      const params = buildParams();
      tableData.value = isRecycleBin.value
        ? await getRecycleMenuTreeList(params)
        : await getMenuTreeList(params);
    } catch (error) {
      console.error(error);
      message.error('菜单列表加载失败，请稍后重试');
    } finally {
      loading.value = false;
    }
  }

  function handleSelectChange(keys: Array<number | string>) {
    selectedRowKeys.value = keys;
  }

  function clearSelectedRowKeys() {
    selectedRowKeys.value = [];
  }

  function handleSearch() {
    void fetchTableData();
  }

  function handleReset() {
    Object.assign(searchForm, createMenuSearchForm());
    void fetchTableData();
  }

  function toggleRecycleBin(next?: boolean) {
    isRecycleBin.value = typeof next === 'boolean' ? next : !isRecycleBin.value;
    clearSelectedRowKeys();
    void fetchTableData();
  }

  return {
    clearSelectedRowKeys,
    fetchTableData,
    handleReset,
    handleSearch,
    handleSelectChange,
    isRecycleBin,
    loading,
    searchForm,
    selectedRowKeys,
    tableData,
    toggleRecycleBin,
  };
}
