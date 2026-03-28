import type { DeptApi } from '#/api/system/dept';

import { reactive, ref } from 'vue';

import { message } from '#/adapter/tdesign';
import { getDeptPageList, getRecycleDeptList } from '#/api/system/dept';

import { createDeptSearchForm } from './schemas';

export function useDeptPage() {
  const searchForm = reactive(createDeptSearchForm());
  const tableData = ref<DeptApi.ListTreeItem[]>([]);
  const loading = ref(false);
  const selectedRowKeys = ref<Array<number | string>>([]);
  const isRecycleBin = ref(false);

  function buildParams() {
    const params: Partial<DeptApi.ListQuery> = {};
    if (searchForm.name) params.name = searchForm.name;
    if (searchForm.leader) params.leader = searchForm.leader;
    if (searchForm.phone) params.phone = searchForm.phone;
    if (searchForm.status !== undefined) params.status = searchForm.status;
    if (searchForm.created_at?.length === 2 && searchForm.created_at[0]) {
      params.created_at = searchForm.created_at;
    }
    return params;
  }

  function clearSelectedRowKeys() {
    selectedRowKeys.value = [];
  }

  function handleSelectChange(keys: Array<number | string>) {
    selectedRowKeys.value = keys;
  }

  function resetSearchForm() {
    Object.assign(searchForm, createDeptSearchForm());
  }

  async function fetchTableData() {
    loading.value = true;
    try {
      const params = buildParams();
      tableData.value = isRecycleBin.value
        ? await getRecycleDeptList(params)
        : await getDeptPageList(params);
    } catch (error) {
      console.error(error);
      message.error('部门列表加载失败，请稍后重试');
    } finally {
      loading.value = false;
    }
  }

  function handleSearch() {
    void fetchTableData();
  }

  function handleReset() {
    resetSearchForm();
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
