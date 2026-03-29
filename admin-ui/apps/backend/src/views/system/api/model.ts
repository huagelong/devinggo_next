import type { ApiManageApi } from '#/api/system/api';
import type { ApiColumnApi } from '#/api/system/api-column';
import type { OptionItem, IdType } from '#/types/common';

export interface ApiSearchFormModel {
  group_id?: IdType;
  name: string;
  access_name: string;
  request_mode?: number | string;
  status?: number;
  created_at: string[];
}

export type ApiListItem = ApiManageApi.ListItem;

export interface ApiFormModel
  extends Omit<ApiManageApi.SubmitPayload, 'group_id' | 'request_mode'> {
  group_id?: IdType;
  request_mode?: number | string;
}

export type ApiColumnOptionItem = OptionItem<string>;

export interface ApiTableColumn {
  align?: 'left' | 'center' | 'right';
  colKey: string;
  fixed?: 'left' | 'right';
  minWidth?: number;
  title?: string;
  type?: 'multiple' | 'single';
  width?: number;
}

export type ApiColumnType = 1 | 2;

export interface ApiColumnSearchFormModel {
  name: string;
  data_type?: string | number;
  status?: number;
  is_required?: number;
  created_at: string[];
}

export type ApiColumnListItem = ApiColumnApi.ListItem;

export interface ApiColumnFormModel
  extends Omit<ApiColumnApi.SubmitPayload, 'data_type'> {
  id?: number;
  data_type?: string | number;
}

export interface ApiColumnTableColumn {
  align?: 'left' | 'center' | 'right';
  colKey: string;
  fixed?: 'left' | 'right';
  minWidth?: number;
  title?: string;
  type?: 'multiple' | 'single';
  width?: number;
}
