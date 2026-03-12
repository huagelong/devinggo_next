import React, { useRef } from 'react';
import { PageContainer, ProTable, ActionType, ProColumns } from '@ant-design/pro-components';
import { Button, Popconfirm, message } from 'antd';
import { PlusOutlined, DeleteOutlined } from '@ant-design/icons';
import request from '../../../utils/request';

// API 交互模拟封装 (待真实对接)
export const getMenuList = (params: any) => request.get('/menu/index', { params });
export const deleteMenu = (id: number) => request.delete(`/menu/delete`, { data: { id } });

const MenuManagement: React.FC = () => {
  const actionRef = useRef<ActionType>();

  const columns: ProColumns<any>[] = [
    { title: '菜单标题', dataIndex: 'title' },
    { title: '图标', dataIndex: 'icon', search: false },
    { title: '路由路径', dataIndex: 'path' },
    { title: '组件', dataIndex: 'component' },
    { title: '权限标识', dataIndex: 'permission' },
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
      render: (text, record, _, action) => [
        <a key="edit" onClick={() => message.info('打开编辑抽屉')}>编辑</a>,
        <a key="addChild" onClick={() => message.info('添加子菜单')}>新增</a>,
        <Popconfirm
          key="delete"
          title="确认删除该菜单及其子菜单吗？"
          onConfirm={async () => {
             // await deleteMenu(record.id);
             message.success('模拟删除成功');
             action?.reload();
          }}
        >
          <a style={{ color: 'red' }}>删除</a>
        </Popconfirm>,
      ],
    },
  ];

  return (
    <PageContainer>
      <ProTable
        headerTitle="菜单管理"
        actionRef={actionRef}
        rowKey="id"
        // request={async (params) => {
        //   const res = await getMenuList(params);
        //   return { data: res.items, success: true, total: res.total };
        // }}
        dataSource={[
           // 模拟数据结构，后端数据需返回 children 字段由 antd 自动渲染树表
           { id: 1, title: '系统管理', path: '/system', status: 1, children: [
               { id: 2, title: '用户管理', path: '/system/user', status: 1 },
               { id: 3, title: '角色管理', path: '/system/role', status: 1 },
               { id: 4, title: '菜单管理', path: '/system/menu', status: 1 },
           ] }
        ]}
        pagination={false}
        columns={columns}
        search={{ labelWidth: 'auto' }}
        toolBarRender={() => [
          <Button key="button" icon={<PlusOutlined />} type="primary" onClick={() => message.info('打开新增抽屉')}>
            新建菜单
          </Button>,
        ]}
      />
    </PageContainer>
  );
};

export default MenuManagement;
