import { createRouter, createRoute, createRootRoute } from '@tanstack/react-router';
import BasicLayout from './components/layout/BasicLayout';

// Root Route
const rootRoute = createRootRoute();

// --- 授权相关 ---
// 因为不强制后端关联，目前暂时手写临时页面占位
const loginRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/login',
  component: () => <div className="flex h-screen items-center justify-center">Login Page Placeholder</div>,
});

// --- 被 BasicLayout 包裹的受保护路由 ---
const layoutRoute = createRoute({
  getParentRoute: () => rootRoute,
  id: 'layout',
  component: BasicLayout,
});

const indexRoute = createRoute({
  getParentRoute: () => layoutRoute,
  path: '/',
  component: () => <div className="text-xl">Welcome to Admin UI</div>,
});

const dashboardRoute = createRoute({
  getParentRoute: () => layoutRoute,
  path: '/dashboard',
  component: () => <div>Dashboard Analysis</div>,
});

const systemUserRoute = createRoute({
  getParentRoute: () => layoutRoute,
  path: '/system/user',
  component: () => <div>System Users</div>,
});

const systemRoleRoute = createRoute({
  getParentRoute: () => layoutRoute,
  path: '/system/role',
  component: () => <div>System Roles</div>,
});

const systemMenuRoute = createRoute({
  getParentRoute: () => layoutRoute,
  path: '/system/menu',
  component: () => <div>System Menus</div>,
});

// Assemble the route tree
const routeTree = rootRoute.addChildren([
  loginRoute,
  layoutRoute.addChildren([
    indexRoute,
    dashboardRoute,
    systemUserRoute,
    systemRoleRoute,
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
