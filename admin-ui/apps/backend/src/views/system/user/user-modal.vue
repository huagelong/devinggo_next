<script lang="ts" setup>
import { useVbenModal } from '@vben/common-ui';

import { MessagePlugin } from 'tdesign-vue-next';

import { useVbenForm } from '#/adapter/form';
import { getDeptTree } from '#/api/system/dept';
import { getPostList } from '#/api/system/post';
import { getRoleList } from '#/api/system/role';
import { saveUser, updateUser } from '#/api/system/user';

const emit = defineEmits(['success']);

const [Form, formApi] = useVbenForm({
  showDefaultActions: false,
  commonConfig: {
    labelWidth: 80,
  },
  wrapperClass: 'grid-cols-2',
  schema: [
    {
      fieldName: 'id',
      label: 'ID',
      component: 'Input',
      dependencies: {
        show: false,
        triggerFields: [''],
      },
    },
    {
      fieldName: 'avatar',
      label: '头像',
      component: 'Input',
      formItemClass: 'col-span-2',
    },
    {
      fieldName: 'username',
      label: '账户',
      component: 'Input',
      rules: 'required',
      componentProps: { placeholder: '请输入账户' },
    },
    {
      fieldName: 'dept_ids',
      label: '所属部门',
      component: 'ApiTreeSelect',
      componentProps: {
        api: getDeptTree,
        labelField: 'name',
        valueField: 'id',
        childrenField: 'children',
        multiple: true,
        placeholder: '请选择所属部门',
      },
      rules: 'required',
    },
    {
      fieldName: 'password',
      label: '密码',
      component: 'InputPassword',
      dependencies: {
        show: (values) => !values?.id,
        triggerFields: ['id'],
      },
      rules: 'required',
      componentProps: { placeholder: '请输入密码' },
    },
    {
      fieldName: 'nickname',
      label: '昵称',
      component: 'Input',
      componentProps: { placeholder: '请输入昵称' },
    },
    {
      fieldName: 'role_ids',
      label: '角色',
      component: 'ApiSelect',
      componentProps: {
        api: async () => {
          const res = await getRoleList();
          return res?.items || res || [];
        },
        labelField: 'name',
        valueField: 'id',
        multiple: true,
        placeholder: '请选择角色',
      },
      rules: 'required',
    },
    {
      fieldName: 'phone',
      label: '手机',
      component: 'Input',
      componentProps: { placeholder: '请输入手机' },
    },
    {
      fieldName: 'post_ids',
      label: '岗位',
      component: 'ApiSelect',
      componentProps: {
        api: async () => {
          const res = await getPostList();
          return res?.items || res || [];
        },
        labelField: 'name',
        valueField: 'id',
        multiple: true,
        placeholder: '请选择岗位',
      },
      rules: 'required',
    },
    {
      fieldName: 'email',
      label: '邮箱',
      component: 'Input',
      componentProps: { placeholder: '请输入邮箱' },
    },
    {
      fieldName: 'status',
      label: '状态',
      component: 'RadioGroup',
      defaultValue: 1,
      componentProps: {
        options: [
          { label: '正常', value: 1 },
          { label: '停用', value: 2 },
        ],
      },
      formItemClass: 'col-span-2',
    },
    {
      fieldName: 'user_type',
      label: '用户类型',
      component: 'Select',
      defaultValue: '100',
      componentProps: {
        options: [{ label: '系统用户', value: '100' }],
        placeholder: '请选择用户类型',
      },
      rules: 'required',
      formItemClass: 'col-span-2',
    },
    {
      fieldName: 'remark',
      label: '备注',
      component: 'Textarea',
      formItemClass: 'col-span-2',
      componentProps: { placeholder: '请输入备注' },
    },
  ],
});

const [Modal, modalApi] = useVbenModal({
  onConfirm: async () => {
    try {
      const { valid } = await formApi.validate();
      if (!valid) return;
      const values = await formApi.getValues();
      modalApi.setState({ confirmLoading: true });
      values.id ? await updateUser(values.id, values) : await saveUser(values);
      MessagePlugin.success(values.id ? '更新成功' : '新增成功');
      emit('success');
      modalApi.close();
    } catch (error) {
      console.error(error);
    } finally {
      modalApi.setState({ confirmLoading: false });
    }
  },
  class: 'w-[800px]',
});

function open(data?: any) {
  modalApi.setState({ title: data?.id ? '编辑管理员' : '新增管理员' });
  modalApi.open();
  formApi.resetForm().then(() => {
    if (data) {
      formApi.setValues(data);
    }
  });
}

defineExpose({
  open,
});
</script>

<template>
  <Modal>
    <Form />
  </Modal>
</template>
