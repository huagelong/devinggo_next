import { createRouter, createRoute, createRootRoute } from '@tanstack/react-router';
import BasicLayout from './components/layout/BasicLayout';
import LoginPage from './pages/login';

// Root Route
const rootRoute = createRootRoute();

// --- 授权相关 ---
const loginRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/login',
  component: LoginPage,
});

// --- 被 BasicLayout 包裹的受保护路由 ---
const layoutRoute = createRoute({
  getParentRoute: () => rootRoute,
  id: 'layout',
  component: BasicLayout,
});

import DashboardPage from './pages/dashboard';

const indexRoute = createRoute({
  getParentRoute: () => layoutRoute,
  path: '/',
  component: () => <div className="text-xl">Welcome to Admin UI</div>,
});

const dashboardRoute = createRoute({
  getParentRoute: () => layoutRoute,
  path: '/dashboard',
  component: DashboardPage,
});

import UserManagementPage from './pages/system/user';

const systemUserRoute = createRoute({
  getParentRoute: () => layoutRoute,
  path: '/system/user',
  component: UserManagementPage,
});

import MenuManagementPage from './pages/system/menu';
import DeptManagementPage from './pages/system/dept';
import RoleManagementPage from './pages/system/role';

const systemRoleRoute = createRoute({
  getParentRoute: () => layoutRoute,
  path: '/system/role',
  component: RoleManagementPage,
});

const systemDeptRoute = createRoute({
  getParentRoute: () => layoutRoute,
  path: '/system/dept',
  component: DeptManagementPage,
});

const systemMenuRoute = createRoute({
  getParentRoute: () => layoutRoute,
  path: '/system/menu',
  component: MenuManagementPage,
});

// Assemble the route tree
const routeTree = rootRoute.addChildren([
  loginRoute,
  layoutRoute.addChildren([
    indexRoute,
    dashboardRoute,
    systemUserRoute,
    systemRoleRoute,
    systemDeptRoute,
    systemMenuRoute,
  ]),
]);

// Create the router
export const router = createRouter({ routeTree });

// Register your router for maximum type safety
declare module '@tanstack/react-router' {
  interface Register {
    router: typeof router;
  }
}
