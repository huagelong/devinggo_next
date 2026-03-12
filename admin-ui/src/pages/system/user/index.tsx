import React, { useRef } from 'react';
import { PageContainer, ProTable, ActionType, ProColumns } from '@ant-design/pro-components';
import { Button, Popconfirm, message, Space, Tag } from 'antd';
import { PlusOutlined, DeleteOutlined } from '@ant-design/icons';
// import request from '../../../utils/request';

const UserManagement: React.FC = () => {
  const actionRef = useRef<ActionType>();

  const columns: ProColumns<any>[] = [
    { title: '用户ID', dataIndex: 'id', search: false },
    { title: '用户名', dataIndex: 'username' },
    { title: '手机号', dataIndex: 'phone' },
    { title: '邮箱', dataIndex: 'email', search: false },
    {
      title: '所属部门',
      dataIndex: 'deptName',
      search: false,
      render: (_, record) => <Tag color="blue">{record.deptName}</Tag>,
    },
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
      width: 200,
      render: (text, record, _, action) => [
        <a key="edit" onClick={() => message.info('打开编辑')}>编辑</a>,
        <a key="resetPwd" onClick={() => message.info('重置密码')}>重置密码</a>,
        <Popconfirm
          key="delete"
          title="确认将该用户移入回收站吗？"
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
        headerTitle="用户信息"
        actionRef={actionRef}
        rowKey="id"
        cardBordered
        // request={...}
        dataSource={[
           { id: 1, username: 'admin', phone: '13888888888', email: 'admin@dev.go', deptName: 'DevingGo总部', status: 1, createdAt: Date.now() },
           { id: 2, username: 'testuser', phone: '13999999999', email: 'test@dev.go', deptName: '研发部', status: 1, createdAt: Date.now() },
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
          <Button key="button" icon={<PlusOutlined />} type="primary" onClick={() => message.info('新建用户')}>
            新增关联用户
          </Button>,
        ]}
      />
    </PageContainer>
  );
};

export default UserManagement;
