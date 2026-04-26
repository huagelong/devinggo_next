<script lang="ts" setup>
import type { ConfigApi } from '#/api/system/config';
import type { ConfigFormModel, ConfigGroup } from './model';

import { computed, onMounted, reactive, ref } from 'vue';

import { Page } from '@vben/common-ui';
import { $t } from '@vben/locales';

import { message } from '#/adapter/tdesign';
import { logger } from '#/utils/logger';

import { AddIcon, InfoCircleIcon, SettingIcon } from 'tdesign-icons-vue-next';
import {
  Button,
  Form,
  FormItem,
  MessagePlugin,
  Space,
  Tabs,
  TabPanel,
} from 'tdesign-vue-next';

import {
  deleteConfigGroup,
  getConfigGroupList,
  getConfigList,
  updateConfigByKeys,
} from '#/api/system/config';

import ConfigFieldRenderer from './components/config-field-renderer.vue';
import ConfigFormModal from './components/config-form-modal.vue';
import ConfigGroupModal from './components/config-group-modal.vue';
import ConfigManageModal from './components/config-manage-modal.vue';
import type { ConfigFieldMeta } from './model';

const configGroupModalRef = ref<{ open: () => void }>();
const configFormModalRef = ref<{ open: () => void }>();
const configManageModalRef = ref<{ open: (groupId: number) => void }>();

const groupList = ref<ConfigGroup[]>([]);
const activeGroupKey = ref<number>();
const groupLoading = ref(false);
const configFieldsMap = reactive<Record<number, ConfigFieldMeta[]>>({});
const configFormMap = reactive<Record<number, ConfigFormModel>>({});

const hasGroups = computed(() => groupList.value.length > 0);

function normalizeOptions(data: unknown): ConfigFieldMeta['config_select_data'] {
  if (Array.isArray(data)) {
    return data as ConfigFieldMeta['config_select_data'];
  }
  if (typeof data === 'string' && data.trim()) {
    try {
      const parsed = JSON.parse(data);
      if (Array.isArray(parsed)) {
        return parsed as ConfigFieldMeta['config_select_data'];
      }
    } catch {
      return undefined;
    }
  }
  return undefined;
}

function normalizeSwitchValues(value: unknown) {
  if (typeof value === 'boolean') {
    return { checked: true, unchecked: false };
  }
  if (typeof value === 'number') {
    return { checked: 1, unchecked: 0 };
  }
  if (typeof value === 'string') {
    const normalized = value.trim().toLowerCase();
    if (normalized === '1' || normalized === '0') {
      return { checked: '1', unchecked: '0' };
    }
    if (normalized === 'true' || normalized === 'false') {
      return { checked: 'true', unchecked: 'false' };
    }
  }
  return { checked: true, unchecked: false };
}

function parseCheckboxValue(rawValue: unknown): string[] {
  if (Array.isArray(rawValue)) {
    return rawValue.map((item) => String(item));
  }
  if (typeof rawValue === 'string') {
    try {
      const parsed = JSON.parse(rawValue);
      if (Array.isArray(parsed)) {
        return parsed.map((item) => String(item));
      }
    } catch {
      if (rawValue.includes(',')) {
        return rawValue.split(',').map((value) => value.trim());
      }
    }
  }
  return [];
}

function parseKeyValueInput(rawValue: unknown) {
  if (Array.isArray(rawValue)) {
    return rawValue;
  }
  if (typeof rawValue === 'string' && rawValue.trim()) {
    try {
      const parsed = JSON.parse(rawValue);
      if (Array.isArray(parsed)) {
        return parsed;
      }
    } catch {
      return [{ key: '', value: rawValue }];
    }
  }
  return [{ key: '', value: '' }];
}

async function fetchGroups() {
  groupLoading.value = true;
  try {
    const list = await getConfigGroupList();
    groupList.value = list;
    if (list.length > 0) {
      const firstGroup = list[0]!;
      activeGroupKey.value = firstGroup.id;
      await fetchGroupConfigs(firstGroup.id);
    } else {
      activeGroupKey.value = undefined;
    }
  } catch (error) {
    logger.error(error);
    message.error($t('common.configGroupLoadFailed'));
  } finally {
    groupLoading.value = false;
  }
}

async function fetchGroupConfigs(groupId: number) {
  try {
    const items = await getConfigList({
      group_id: groupId,
      orderBy: 'sort',
      orderType: 'asc',
    });
    const fields: ConfigFieldMeta[] = [];
    const form: ConfigFormModel = {};
    items.forEach((item) => {
      const field: ConfigFieldMeta = {
        id: item.id,
        input_type: item.input_type,
        key: item.key,
        label: item.name,
        remark: item.remark,
        sort: item.sort,
        config_select_data: normalizeOptions(item.config_select_data),
      };
      let value: ConfigFormModel[string];
      if (item.input_type === 'switch') {
        field.switchValues = normalizeSwitchValues(item.value);
      }
      if (item.input_type === 'checkbox') {
        value = parseCheckboxValue(item.value);
      } else if (item.input_type === 'key-value') {
        value = parseKeyValueInput(item.value);
      } else {
        value = item.value as ConfigFormModel[string];
      }
      form[item.key] = value;
      fields.push(field);
    });
    configFieldsMap[groupId] = fields;
    configFormMap[groupId] = form;
  } catch (error) {
    logger.error(error);
    message.error($t('common.configDataLoadFailed'));
  }
}

function handleTabChange(value: string | number) {
  const id = Number(value);
  if (Number.isNaN(id)) return;
  activeGroupKey.value = id;
  if (!configFieldsMap[id]) {
    void fetchGroupConfigs(id);
  }
}

function formatSubmitValue(
  field: ConfigFieldMeta,
  value: unknown,
): ConfigApi.UpdateByKeysPayload[string] {
  if (
    field.input_type === 'key-value' ||
    Array.isArray(value) ||
    (value !== null && typeof value === 'object')
  ) {
    try {
      return JSON.stringify(value ?? []);
    } catch {
      return JSON.stringify([]);
    }
  }
  return value as ConfigApi.UpdateByKeysPayload[string];
}

async function handleSubmit(groupId: number) {
  const currentForm = configFormMap[groupId];
  if (!currentForm) return;
  const fields = configFieldsMap[groupId] ?? [];
  const payload: ConfigApi.UpdateByKeysPayload = {};
  fields.forEach((field) => {
    payload[field.key] = formatSubmitValue(field, currentForm[field.key]);
  });

  try {
    await updateConfigByKeys(payload);
    MessagePlugin.success($t('common.configUpdateSuccess'));
    await fetchGroupConfigs(groupId);
  } catch (error) {
    logger.error(error);
    MessagePlugin.error($t('common.configUpdateFailed'));
  }
}

async function handleDeleteGroup(groupId: number) {
  const group = groupList.value.find((item) => item.id === groupId);
  if (!group) return;
  if (groupId === 1 || groupId === 2) {
    MessagePlugin.info($t('common.defaultGroupCannotDelete'));
    return;
  }
  try {
    await deleteConfigGroup({ id: groupId });
    MessagePlugin.success($t('common.groupDeleteSuccess'));
    await fetchGroups();
  } catch (error) {
    logger.error(error);
    MessagePlugin.error($t('common.groupDeleteFailed'));
  }
}

onMounted(() => {
  void fetchGroups();
});
</script>

<template>
  <Page auto-content-height>
    <div class="flex h-full flex-col gap-5">
      <div class="flex flex-col gap-5 lg:flex-row">
        <div class="config-card lg:w-1/2">
          <div class="config-toolbar flex items-center justify-between">
            <Space>
              <Button
                theme="primary"
                variant="outline"
                @click="configGroupModalRef?.open()"
              >
                <template #icon><AddIcon /></template>
                {{ $t('system.config.addGroup') }}
              </Button>
              <Button
                theme="default"
                variant="outline"
                :disabled="!activeGroupKey"
                @click="
                  activeGroupKey &&
                  configManageModalRef?.open(Number(activeGroupKey))
                "
              >
                <template #icon><SettingIcon /></template>
                {{ $t('system.config.manageTitle') }}
              </Button>
              <Button
                v-if="activeGroupKey && activeGroupKey > 2"
                theme="danger"
                variant="outline"
                @click="handleDeleteGroup(activeGroupKey)"
              >
                {{ $t('system.config.deleteGroup') }}
              </Button>
            </Space>
          </div>
          <div v-if="hasGroups" class="config-content">
            <Tabs
              v-model="activeGroupKey"
              :loading="groupLoading"
              class="config-tabs"
              @change="handleTabChange"
            >
              <TabPanel
                v-for="group in groupList"
                :key="group.id"
                :value="group.id"
                :label="group.name"
              >
                <Form
                  v-if="configFormMap[group.id]"
                  :data="configFormMap[group.id]!"
                  class="config-edit-form"
                  label-width="200px" layout="inline" colon
                >
                  <div class="config-form-list">
                    <FormItem
                      v-for="field in configFieldsMap[group.id] ?? []"
                      :key="field.id"
                      :label="field.label"
                      :name="field.key"
                      class="config-field-item"
                    >
                      <ConfigFieldRenderer
                        v-model="configFormMap[group.id]![field.key]"
                        :field="field"
                      />
                      <div
                        v-if="field.remark"
                        class="config-field-remark"
                      >
                        {{ field.remark }}
                      </div>
                    </FormItem>
                  </div>
                  <div class="config-save-bar">
                    <Button theme="primary" @click="handleSubmit(group.id)">
                      {{ $t('system.config.saveConfig') }}
                    </Button>
                  </div>
                </Form>
                <div
                  v-else
                  class="flex min-h-[200px] items-center justify-center text-gray-500"
                >
                  <InfoCircleIcon class="mr-2" /> {{ $t('system.config.configLoading') }}
                </div>
              </TabPanel>
            </Tabs>
          </div>
          <div
            v-else
            class="flex min-h-[200px] items-center justify-center text-gray-500"
          >
            {{ $t('system.config.noConfigGroup') }}
          </div>
        </div>

        <div class="config-card lg:w-1/2">
          <div class="config-right-head flex items-center justify-between">
            <h3 class="text-lg font-semibold">{{ $t('system.config.addConfigTitle') }}</h3>
            <Button theme="primary" @click="configFormModalRef?.open()">
              <template #icon><AddIcon /></template>
              {{ $t('common.create') }}
            </Button>
          </div>
          <p class="config-add-tip text-sm">
            {{ $t('system.config.addConfigTip') }}
          </p>
          <div class="config-empty-state">
            <InfoCircleIcon class="config-empty-icon" />
            <span>{{ $t('system.config.addConfigTitle') }}</span>
          </div>
        </div>
      </div>
    </div>

    <ConfigGroupModal ref="configGroupModalRef" @success="fetchGroups" />
    <ConfigFormModal ref="configFormModalRef" @success="fetchGroups" />
    <ConfigManageModal ref="configManageModalRef" />
  </Page>
</template>

<style scoped>
.config-card {
  border: 1px solid var(--td-component-border, #e7e7e7);
  border-radius: 10px;
  background: #fff;
  padding: 16px 18px;
  box-shadow: 0 6px 18px rgb(15 23 42 / 4%);
}

.config-toolbar {
  margin-bottom: 14px;
  padding-bottom: 12px;
  border-bottom: 1px solid var(--td-component-stroke, #f0f1f2);
}

.config-toolbar :deep(.t-space) {
  flex-wrap: wrap;
  row-gap: 8px;
}

.config-content {
  min-height: 420px;
}

.config-form-list {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.config-edit-form {
  display: block;
}

.config-field-item {
  margin-bottom: 0 !important;
  padding: 12px 8px;
  border-bottom: 1px dashed #f0f2f5;
}

.config-field-remark {
  margin-left: 10px;
  font-size: 13px;
  color: #8a94a6;
  line-height: 1.5;
}

.config-field-item:last-child {
  border-bottom: none;
}

.config-edit-form :deep(.t-form__item) {
  width: 100%;
  margin-right: 0;
}

.config-edit-form :deep(.t-form__label) {
  padding-right: 12px;
  color: var(--td-text-color-secondary, #6b7280);
  font-weight: 500;
  line-height: 40px;
}

.config-edit-form :deep(.t-form__controls) {
  flex: 1;
  min-width: 0;
}

.config-edit-form :deep(.t-input),
.config-edit-form :deep(.t-textarea__inner),
.config-edit-form :deep(.t-select),
.config-edit-form :deep(.t-select__wrap) {
  border-radius: 8px;
}

.config-save-bar {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
}

.config-right-head {
  margin-bottom: 16px;
}

.config-add-tip {
  margin: 0;
  color: var(--td-text-color-secondary, #6b7280);
  line-height: 1.7;
}

.config-empty-state {
  margin-top: 20px;
  min-height: 320px;
  border: 1px dashed #d9e3f3;
  border-radius: 12px;
  background: linear-gradient(180deg, #f8fbff 0%, #ffffff 100%);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 10px;
  color: #5f6f85;
}

.config-empty-icon {
  font-size: 24px;
  color: #3b82f6;
}

@media (max-width: 1023px) {
  .config-card {
    padding: 14px;
  }

  .config-content {
    min-height: auto;
  }

  .config-empty-state {
    min-height: 200px;
  }
}
</style>
