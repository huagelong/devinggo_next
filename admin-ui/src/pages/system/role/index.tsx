import React, { useRef } from 'react';
import { PageContainer, ProTable, ActionType, ProColumns } from '@ant-design/pro-components';
import { Button, Popconfirm, message } from 'antd';
import { PlusOutlined, SettingOutlined } from '@ant-design/icons';
// import request from '../../../utils/request';

const RoleManagement: React.FC = () => {
  const actionRef = useRef<ActionType>();

  const columns: ProColumns<any>[] = [
    { title: '角色ID', dataIndex: 'id', search: false },
    { title: '角色名称', dataIndex: 'name' },
    { title: '角色代码', dataIndex: 'code' },
    { title: '备注', dataIndex: 'remark', search: false },
    {
      title: '状态',
      dataIndex: 'status',
      valueType: 'select',
      valueEnum: {
        1: { text: '正常', status: 'Success' },
        2: { text: '停用', status: 'Error' },
      },
    },
    { title: '创建时间', dataIndex: 'createdAt', valueType: 'dateTime', search: false },
    {
      title: '操作',
      valueType: 'option',
      width: 300,
      render: (text, record, _, action) => [
        <a key="edit" onClick={() => message.info('打开编辑')}>编辑</a>,
        <a key="menuAuth" onClick={() => message.info('分配菜单权限')}><SettingOutlined /> 菜单权限</a>,
        <a key="dataAuth" onClick={() => message.info('分配数据权限')}><SettingOutlined /> 数据权限</a>,
        <Popconfirm
          key="delete"
          title="确认删除该角色吗？"
          onConfirm={() => message.success('模拟删除成功')}
        >
          <a style={{ color: 'red' }}>删除</a>
        </Popconfirm>,
      ],
    },
  ];

  return (
    <PageContainer ghost header={{ title: '' }}>
      <ProTable
        headerTitle="访问角色管理"
        actionRef={actionRef}
        rowKey="id"
        cardBordered
        // request={...}
        dataSource={[
           { id: 1, name: '超级管理员', code: 'super_admin', remark: '系统超级管理员，不可删除', status: 1, createdAt: Date.now() },
           { id: 2, name: '普通用户', code: 'common_user', remark: '普通的系统用户', status: 1, createdAt: Date.now() },
        ]}
        pagination={{ pageSize: 10, showSizeChanger: true }}
        columns={columns}
        search={{
          labelWidth: 'auto',
          className: 'shadow-sm rounded-lg mb-4',
        }}
        options={{
          setting: { listsHeight: 400 },
        }}
        toolBarRender={() => [
          <Button key="button" icon={<PlusOutlined />} type="primary" onClick={() => message.info('新建角色')}>
            新增角色绑定
          </Button>,
        ]}
      />
    </PageContainer>
  );
};

export default RoleManagement;
