import type { ConfigApi } from '#/api/system/config';
import type { OptionItem } from '#/types/common';

export interface ConfigGroup extends ConfigApi.ConfigGroupItem {}

export type ConfigListItem = ConfigApi.ConfigItem;

export interface ConfigFieldMeta {
  config_select_data?: ConfigApi.ConfigSelectOption[];
  id: number;
  input_type: ConfigApi.InputType;
  key: string;
  label: string;
  remark?: string;
  sort?: number;
  switchValues?: {
    checked: string | number | boolean;
    unchecked: string | number | boolean;
  };
}

export type ConfigKeyValueItem = { key: string; value: string };

export interface ConfigFormModel {
  [key: string]:
    | string
    | number
    | boolean
    | string[]
    | Record<string, unknown>
    | ConfigKeyValueItem[]
    | undefined;
}

export type ConfigGroupOption = OptionItem<number>;

export const inputComponentOptions: ConfigApi.ConfigSelectOption[] = [
  { label: '文本框', value: 'input' },
  { label: '文本域', value: 'textarea' },
  { label: '下拉选择', value: 'select' },
  { label: '单选', value: 'radio' },
  { label: '多选', value: 'checkbox' },
  { label: '开关', value: 'switch' },
  { label: '图片上传', value: 'upload' },
  { label: '键值对', value: 'key-value' },
  { label: '富文本', value: 'editor' },
];
