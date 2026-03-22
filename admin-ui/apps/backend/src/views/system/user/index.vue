<script lang="ts" setup>
import { defineOptions, ref } from 'vue';

import { Page } from '@vben/common-ui';

import { Button, MessagePlugin, Popconfirm, Switch } from 'tdesign-vue-next';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  changeUserStatus,
  deleteUser,
  getRecycleUserList,
  getUserList,
  realDeleteUser,
  recoveryUser,
  resetPassword,
} from '#/api/system/user';

import DeptTree from './dept-tree.vue';
import UserModal from './user-modal.vue';

defineOptions({ name: 'SystemUser' });

const currentDeptId = ref<number | string>('');
const isRecycleBin = ref(false);

const userModalRef = ref();

const [Grid, gridApi] = useVbenVxeGrid({
  formOptions: {
    handleReset: () => {
      currentDeptId.value = '';
      gridApi.reload();
    },
    schema: [
      {
        fieldName: 'username',
        label: '账户',
        component: 'Input',
        componentProps: { placeholder: '请输入账户' },
      },
      {
        fieldName: 'phone',
        label: '手机',
        component: 'Input',
        componentProps: { placeholder: '请输入手机' },
      },
      {
        fieldName: 'status',
        label: '状态',
        component: 'Select',
        componentProps: {
          placeholder: '请选择状态',
          options: [
            { label: '正常', value: 1 },
            { label: '停用', value: 2 },
          ],
        },
      },
    ],
  },
  gridOptions: {
    toolbarConfig: {
      custom: true,
      refresh: true,
      zoom: true,
    },
    proxyConfig: {
      ajax: {
        query: async ({ page }, formValues) => {
          const params = {
            page: page.currentPage,
            limit: page.pageSize,
            ...formValues,
          };
          if (currentDeptId.value) {
            params.dept_id = currentDeptId.value;
          }
          if (isRecycleBin.value) {
            return await getRecycleUserList(params);
          }
          return await getUserList(params);
        },
      },
    },
    columns: [
      { type: 'checkbox', width: 60, align: 'center' },
      {
        field: 'avatar',
        title: '头像',
        width: 80,
        cellRender: { name: 'CellImage' },
        align: 'center',
      },
      { fieldName: 'username', title: '账户', minWidth: 100, align: 'center' },
      { field: 'nickname', title: '昵称', minWidth: 100, align: 'center' },
      { fieldName: 'phone', title: '手机', minWidth: 120, align: 'center' },
      { field: 'email', title: '邮箱', minWidth: 150, align: 'center' },
      {
        field: 'status',
        title: '状态',
        width: 100,
        slots: { default: 'status' },
        align: 'center',
      },
      { field: 'user_type', title: '用户类型', width: 100, align: 'center' },
      { field: 'created_at', title: '注册时间', width: 160, align: 'center' },
      {
        field: 'action',
        title: '操作',
        width: 220,
        fixed: 'right',
        slots: { default: 'action' },
        align: 'center',
      },
    ],
  },
});

function handleDeptSelect(deptId: number | string) {
  currentDeptId.value = deptId;
  gridApi.reload();
}

function handleAdd() {
  userModalRef.value?.open();
}

function handleEdit(row: any) {
  userModalRef.value?.open(row);
}

async function handleDelete(row: any) {
  try {
    await (isRecycleBin.value
      ? realDeleteUser([row.id])
      : deleteUser([row.id]));
    MessagePlugin.success('删除成功');
    gridApi.reload();
  } catch (error) {
    console.error(error);
  }
}

async function handleBatchDelete() {
  const records = gridApi.grid.getCheckboxRecords();
  if (records.length === 0) {
    MessagePlugin.warning('请选择需要操作的数据');
    return;
  }
  const ids = records.map((item: any) => item.id);
  try {
    await (isRecycleBin.value ? realDeleteUser(ids) : deleteUser(ids));
    MessagePlugin.success('操作成功');
    gridApi.reload();
  } catch (error) {
    console.error(error);
  }
}

async function handleRecovery(row: any) {
  try {
    await recoveryUser([row.id]);
    MessagePlugin.success('恢复成功');
    gridApi.reload();
  } catch (error) {
    console.error(error);
  }
}

async function handleBatchRecovery() {
  const records = gridApi.grid.getCheckboxRecords();
  if (records.length === 0) {
    MessagePlugin.warning('请选择需要操作的数据');
    return;
  }
  const ids = records.map((item: any) => item.id);
  try {
    await recoveryUser(ids);
    MessagePlugin.success('恢复成功');
    gridApi.reload();
  } catch (error) {
    console.error(error);
  }
}

async function handleStatusChange(row: any, val: boolean) {
  const status = val ? 1 : 2;
  try {
    await changeUserStatus({ id: row.id, status });
    MessagePlugin.success('更新状态成功');
    gridApi.reload();
  } catch (error) {
    console.error(error);
  }
}

async function handleResetPassword(row: any) {
  try {
    await resetPassword({ id: row.id });
    MessagePlugin.success('密码重置成功');
  } catch (error) {
    console.error(error);
  }
}

function handleSuccess() {
  gridApi.reload();
}

function toggleRecycleBin() {
  isRecycleBin.value = !isRecycleBin.value;
  gridApi.reload();
}
</script>

<template>
  <Page auto-content-height class="p-4 bg-gray-100">
    <div class="flex flex-row gap-4 h-full layout-container">
      <div class="h-full bg-white shadow-sm rounded flex-shrink-0">
        <DeptTree @select="handleDeptSelect" />
      </div>
      <div class="flex-1 min-w-0 h-full">
        <Grid>
          <template #toolbar-tools>
            <Button v-if="!isRecycleBin" theme="primary" @click="handleAdd">
              新增
            </Button>
            <Button
              v-if="!isRecycleBin"
              theme="danger"
              variant="outline"
              @click="handleBatchDelete"
            >
              删除
            </Button>
            <Button v-if="!isRecycleBin" variant="outline"> 导入 </Button>
            <Button v-if="!isRecycleBin" variant="outline"> 导出 </Button>

            <Button
              v-if="isRecycleBin"
              theme="success"
              @click="handleBatchRecovery"
            >
              恢复
            </Button>
            <Button
              v-if="isRecycleBin"
              theme="danger"
              variant="outline"
              @click="handleBatchDelete"
            >
              彻底删除
            </Button>

            <Button variant="outline" @click="toggleRecycleBin" class="ml-2">
              {{ isRecycleBin ? '返回列表' : '显示回收站' }}
            </Button>
          </template>

          <template #status="{ row }">
            <Switch
              :value="row.status === 1"
              @change="(val) => handleStatusChange(row, val as boolean)"
              :disabled="isRecycleBin"
            />
          </template>

          <template #action="{ row }">
            <div class="flex gap-2">
              <template v-if="!isRecycleBin">
                <Button
                  theme="primary"
                  variant="text"
                  size="small"
                  @click="handleEdit(row)"
                >
                  编辑
                </Button>
                <Popconfirm
                  content="确认删除该用户吗？"
                  @confirm="handleDelete(row)"
                >
                  <Button theme="danger" variant="text" size="small">
                    删除
                  </Button>
                </Popconfirm>
                <Popconfirm
                  content="确认重置该用户密码吗？"
                  @confirm="handleResetPassword(row)"
                >
                  <Button theme="warning" variant="text" size="small">
                    重置密码
                  </Button>
                </Popconfirm>
              </template>
              <template v-else>
                <Popconfirm
                  content="确认恢复该用户吗？"
                  @confirm="handleRecovery(row)"
                >
                  <Button theme="primary" variant="text" size="small">
                    恢复
                  </Button>
                </Popconfirm>
                <Popconfirm
                  content="确认彻底删除该用户吗？"
                  @confirm="handleDelete(row)"
                >
                  <Button theme="danger" variant="text" size="small">
                    彻底删除
                  </Button>
                </Popconfirm>
              </template>
            </div>
          </template>
        </Grid>
      </div>
    </div>

    <UserModal ref="userModalRef" @success="handleSuccess" />
  </Page>
</template>

<style scoped>
.layout-container {
  overflow: hidden;
}

:deep(.vben-vxe-grid) {
  height: 100%;
}
</style>
