import React from 'react';
import { ProLayout } from '@ant-design/pro-components';
import { Outlet, useNavigate, useLocation } from '@tanstack/react-router';
import { Dropdown, MenuProps } from 'antd';
import { LogoutOutlined, UserOutlined } from '@ant-design/icons';
import { useAuthStore } from '../../stores/authStore';
import { useTranslation } from 'react-i18next';

export const BasicLayout: React.FC = () => {
  const navigate = useNavigate();
  const location = useLocation();
  const { userInfo, logout } = useAuthStore();
  const { t } = useTranslation();

  const handleLogout = () => {
    logout();
    navigate({ to: '/login' });
  };

  const userMenuItems: MenuProps['items'] = [
    {
      key: 'profile',
      icon: <UserOutlined />,
      label: '个人中心',
    },
    {
      key: 'logout',
      icon: <LogoutOutlined />,
      label: t('system.logout'),
      onClick: handleLogout,
    },
  ];

  return (
    <div style={{ height: '100vh', display: 'flex', flexDirection: 'column' }}>
      <ProLayout
        title="Admin UI"
        logo={<div className="h-8 w-8 bg-blue-500 rounded-full" />}
        layout="mix"
        splitMenus={false}
        fixSiderbar
        avatarProps={{
          src: userInfo?.avatar || 'https://gw.alipayobjects.com/zos/antfincdn/efFD%24IOql2/weixintupian_20170331104822.jpg',
          title: userInfo?.username || 'Admin',
          size: 'small',
          render: (_props, dom) => {
            return (
              <Dropdown menu={{ items: userMenuItems }} trigger={['click']}>
                <div style={{ cursor: 'pointer' }}>{dom}</div>
              </Dropdown>
            );
          },
        }}
        route={{
          path: '/',
          routes: [
            {
              path: '/dashboard',
              name: t('menu.dashboard'),
            },
            {
              path: '/system',
              name: t('menu.system'),
              routes: [
                { path: '/system/user', name: t('menu.user') },
                { path: '/system/role', name: t('menu.role') },
                { path: '/system/dept', name: '部门管理' },
                { path: '/system/menu', name: t('menu.menu') },
              ],
            },
          ],
        }}
        location={{
          pathname: location.pathname,
        }}
        menuItemRender={(item, dom) => (
          <a
            onClick={(e) => {
              e.preventDefault();
              if (item.path) {
                navigate({ to: item.path });
              }
            }}
          >
            {dom}
          </a>
        )}
      >
        <div style={{ minHeight: 'calc(100vh - 110px)', padding: '24px' }}>
          <Outlet />
        </div>
      </ProLayout>
    </div>
  );
};

export default BasicLayout;
