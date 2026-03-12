import { create } from 'zustand';
import { persist } from 'zustand/middleware';

interface AuthState {
  token: string | null;
  userInfo: any | null; // 可以根据 user.go 定义具体类型
  setToken: (token: string) => void;
  setUserInfo: (info: any) => void;
  logout: () => void;
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      token: null,
      userInfo: null,
      setToken: (token: string) => set({ token }),
      setUserInfo: (info: any) => set({ userInfo: info }),
      logout: () => {
        set({ token: null, userInfo: null });
        // 可选：触发清理后跳转路由
      },
    }),
    {
      name: 'auth-storage', // localStorage key
    }
  )
);
