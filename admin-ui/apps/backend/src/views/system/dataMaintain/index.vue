<script lang="ts" setup>
import type { DataMaintainApi } from '#/api/system/data-maintain';

import { computed, onMounted, ref } from 'vue';

import { useAccess } from '@vben/access';
import { Page } from '@vben/common-ui';

import { message } from '#/adapter/tdesign';
import {
  fragmentDataMaintainTable,
  getDataMaintainDetailed,
  optimizeDataMaintainTable,
} from '#/api/system/data-maintain';
import CrudToolbar from '#/components/crud/crud-toolbar.vue';

import {
  InfoCircleFilledIcon,
  SearchIcon,
  ViewIcon,
} from 'tdesign-icons-vue-next';
import {
  Button,
  Form,
  FormItem,
  Input,
  Popup,
  Space,
  Table,
  Tag,
} from 'tdesign-vue-next';

import type { DataMaintainListItem, DataMaintainTableColumn } from './model';
import {
  createDataMaintainColumnOptions,
  createDataMaintainTableColumns,
} from './schemas';
import { useDataMaintainCrud } from './use-data-maintain-crud';

defineOptions({ name: 'SystemDataMaintain' });

const { hasAccessByCodes } = useAccess();
const canView = computed(() =>
  hasAccessByCodes(['system:dataMaintain:index', 'system:dataMaintain']),
);
const canDetailed = computed(() =>
  hasAccessByCodes(['system:dataMaintain:detailed', 'system:dataMaintain:index']),
);
const canOptimize = computed(() =>
  hasAccessByCodes(['system:dataMaintain:optimize', 'system:dataMaintain:index']),
);
const canFragment = computed(() =>
  hasAccessByCodes(['system:dataMaintain:fragment', 'system:dataMaintain:index']),
);

// Current backend only exposes index; detailed/optimize/fragment are prepared for future rollout.
const hasDetailedApi = false;
const hasOptimizeApi = false;
const hasFragmentApi = false;

const detailVisible = ref(false);
const detailLoading = ref(false);
const currentTable = ref<DataMaintainListItem>();
const detailColumns = ref<Array<{ field: string; type?: string; comment?: string }>>([]);

const columns: DataMaintainTableColumn[] = createDataMaintainTableColumns();
const columnOptions = createDataMaintainColumnOptions(columns);
const allColumnKeys = columnOptions.map((item) => item.value);
const visibleColumns = ref<string[]>([...allColumnKeys]);

const displayColumns = computed({
  get: () => [...visibleColumns.value],
  set: (value: string[]) => {
    visibleColumns.value = value;
  },
});

const {
  fetchTableData,
  handlePageChange,
  handleReset,
  handleSearch,
  loading,
  pagination,
  searchForm,
  tableData,
} = useDataMaintainCrud();

function handleUnimplementedAction(actionName: string) {
  message.info(`${actionName}接口暂未在当前后端开放`);
}

async function handleViewDetail(row: DataMaintainListItem) {
  currentTable.value = row;
  detailVisible.value = true;
  detailColumns.value = [];

  if (!hasDetailedApi) {
    return;
  }

  detailLoading.value = true;
  try {
    const response = await getDataMaintainDetailed({
      group_name: searchForm.group_name,
      table_name: row.name,
    });
    detailColumns.value = Object.values(response || {}).map((item) => ({
      comment: item.comment,
      field: item.field,
      type: item.type,
    }));
  } catch (error) {
    console.error(error);
    message.error('获取字段详情失败，请稍后重试');
  } finally {
    detailLoading.value = false;
  }
}

async function handleOptimize(row: DataMaintainListItem) {
  if (!hasOptimizeApi) {
    handleUnimplementedAction('优化表');
    return;
  }

  try {
    await optimizeDataMaintainTable({
      group_name: searchForm.group_name,
      table_name: row.name,
    });
    message.success('优化成功');
  } catch (error) {
    console.error(error);
    message.error('优化失败，请稍后重试');
  }
}

async function handleFragment(row: DataMaintainListItem) {
  if (!hasFragmentApi) {
    handleUnimplementedAction('清理碎片');
    return;
  }

  try {
    await fragmentDataMaintainTable({
      group_name: searchForm.group_name,
      table_name: row.name,
    });
    message.success('清理成功');
  } catch (error) {
    console.error(error);
    message.error('清理失败，请稍后重试');
  }
}

onMounted(() => {
  if (!canView.value) {
    return;
  }
  void fetchTableData();
});
</script>

<template>
  <Page auto-content-height>
    <div class="flex h-full flex-col gap-3">
      <div class="rounded-md bg-white p-4">
        <Form :data="searchForm" label-width="90px" colon>
          <div class="grid grid-cols-4 gap-x-4">
            <FormItem label="数据库组" name="group_name">
              <Input
                v-model="searchForm.group_name"
                placeholder="默认 default"
                clearable
              />
            </FormItem>
            <FormItem label="表名" name="name">
              <Input v-model="searchForm.name" placeholder="请输入表名" clearable />
            </FormItem>
          </div>
          <div class="flex justify-end gap-2 pt-2">
            <Button theme="default" @click="handleReset">重置</Button>
            <Button theme="primary" @click="handleSearch">
              <template #icon><SearchIcon /></template>
              查询
            </Button>
          </div>
        </Form>
      </div>

      <div class="flex min-h-0 flex-1 flex-col rounded-md bg-white p-4">
        <div class="mb-3 flex items-center justify-between">
          <Space>
            <Button variant="outline" @click="handleUnimplementedAction('查看字段')">
              查看字段
            </Button>
            <Button variant="outline" @click="handleUnimplementedAction('优化表')">
              优化表
            </Button>
            <Button variant="outline" @click="handleUnimplementedAction('清理碎片')">
              清理碎片
            </Button>
            <Popup placement="bottom" trigger="hover" content="首版仅接入列表能力，扩展动作待后端接口开放后补齐。">
              <InfoCircleFilledIcon class="cursor-help text-gray-400" />
            </Popup>
          </Space>

          <CrudToolbar
            v-model="displayColumns"
            :column-options="columnOptions"
            :is-recycle-bin="false"
            @refresh="fetchTableData"
          />
        </div>

        <div v-if="!canView" class="rounded-md border border-dashed border-gray-300 p-6 text-center text-gray-500">
          无权限访问数据维护列表（需要 `system:dataMaintain:index`）。
        </div>

        <div v-else class="min-h-0 flex-1">
          <Table
            row-key="name"
            hover
            stripe
            :columns="columns"
            :column-controller-visible="false"
            :display-columns="displayColumns"
            :data="tableData"
            :loading="loading"
            :pagination="pagination"
            @page-change="handlePageChange"
          >
            <template #rows="{ row }">
              {{ row.rows ?? '-' }}
            </template>
            <template #comment="{ row }">
              <span :title="row.comment || '-'">{{ row.comment || '-' }}</span>
            </template>
            <template #create_time="{ row }">
              {{ row.create_time || '-' }}
            </template>
            <template #action="{ row }">
              <Space>
                <Button
                  v-if="canDetailed"
                  size="small"
                  variant="text"
                  @click="handleViewDetail(row)"
                >
                  详情
                </Button>
                <Button
                  v-if="canOptimize"
                  size="small"
                  theme="warning"
                  variant="text"
                  @click="handleOptimize(row)"
                >
                  优化
                </Button>
                <Button
                  v-if="canFragment"
                  size="small"
                  theme="danger"
                  variant="text"
                  @click="handleFragment(row)"
                >
                  碎片整理
                </Button>
              </Space>
            </template>
          </Table>

          <div
            v-if="detailVisible"
            class="mt-3 rounded-md border border-gray-100 bg-gray-50 p-4"
          >
            <div class="mb-2 flex items-center justify-between">
              <div class="text-sm font-medium text-gray-700">
                表详情：{{ currentTable?.name || '-' }}
              </div>
              <Button size="small" variant="text" @click="detailVisible = false">
                收起
              </Button>
            </div>

            <div class="mb-3 grid grid-cols-3 gap-3 text-sm text-gray-600">
              <div>引擎：{{ currentTable?.engine || '-' }}</div>
              <div>字符集：{{ currentTable?.collation || '-' }}</div>
              <div>行数：{{ currentTable?.rows ?? '-' }}</div>
            </div>

            <div v-if="!hasDetailedApi" class="text-sm text-gray-500">
              当前后端未开放字段详情接口，已预留展示区域。
            </div>

            <Table
              v-else
              row-key="field"
              size="small"
              :loading="detailLoading"
              :data="detailColumns"
              :columns="[
                { colKey: 'field', title: '字段名', width: 220 },
                { colKey: 'type', title: '类型', width: 180 },
                { colKey: 'comment', title: '注释', minWidth: 240 },
              ]"
            />

            <div class="mt-3 flex items-center gap-2 text-xs text-gray-500">
              <Tag theme="warning" variant="light">能力预留</Tag>
              <span>详细字段、优化、碎片整理待后端接口开放后无缝启用。</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </Page>
</template>
