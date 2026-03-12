import React, { useRef } from 'react';
import { PageContainer, ProTable, ActionType, ProColumns } from '@ant-design/pro-components';
import { Button, Popconfirm, message } from 'antd';
import { PlusOutlined } from '@ant-design/icons';
// import request from '../../../utils/request';

const DeptManagement: React.FC = () => {
  const actionRef = useRef<ActionType>();

  const columns: ProColumns<any>[] = [
    { title: '部门名称', dataIndex: 'name' },
    { title: '部门编码', dataIndex: 'code' },
    { title: '负责人', dataIndex: 'leader', search: false },
    { title: '电话', dataIndex: 'phone', search: false },
    { title: '邮箱', dataIndex: 'email', search: false },
    { title: '排序', dataIndex: 'sort', search: false },
    {
      title: '状态',
      dataIndex: 'status',
      valueType: 'select',
      valueEnum: {
        1: { text: '正常', status: 'Success' },
        2: { text: '停用', status: 'Error' },
      },
    },
    {
      title: '操作',
      valueType: 'option',
      width: 250,
      render: (text, record, _, action) => [
        <a key="edit" onClick={() => message.info('打开编辑')}>编辑</a>,
        <a key="addChild" onClick={() => message.info('添加下级')}>新增下级</a>,
        <Popconfirm
          key="delete"
          title="确认删除该部门吗？"
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
        headerTitle="组织架构管理"
        actionRef={actionRef}
        rowKey="id"
        cardBordered
        // 模拟树状数据
        dataSource={[
           { id: 1, name: 'DevingGo总部', code: 'HQ', leader: 'Admin', status: 1, children: [
               { id: 2, name: '研发部', code: 'RD', leader: '张三', status: 1 },
               { id: 3, name: '产品部', code: 'PM', leader: '李四', status: 1 },
           ] }
        ]}
        pagination={false}
        columns={columns}
        search={{
          labelWidth: 'auto',
          className: 'shadow-sm rounded-lg mb-4',
        }}
        options={{
          setting: { listsHeight: 400 },
        }}
        toolBarRender={() => [
          <Button key="button" icon={<PlusOutlined />} type="primary" onClick={() => message.info('新建根部门')}>
            新建部门节点
          </Button>,
        ]}
      />
    </PageContainer>
  );
};

export default DeptManagement;
