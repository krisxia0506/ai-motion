import { create } from 'zustand';
import { persist } from 'zustand/middleware';

type Theme = 'light' | 'dark';

interface UIState {
  theme: Theme;
  sidebarOpen: boolean;
  notificationCount: number;
  
  setTheme: (theme: Theme) => void;
  toggleTheme: () => void;
  setSidebarOpen: (open: boolean) => void;
  toggleSidebar: () => void;
  setNotificationCount: (count: number) => void;
  incrementNotificationCount: () => void;
  clearNotifications: () => void;
}

export const useUIStore = create<UIState>()(
  persist(
    (set) => ({
      theme: 'light',
      sidebarOpen: true,
      notificationCount: 0,

      setTheme: (theme) => set({ theme }),
      
      toggleTheme: () =>
        set((state) => ({ theme: state.theme === 'light' ? 'dark' : 'light' })),
      
      setSidebarOpen: (open) => set({ sidebarOpen: open }),
      
      toggleSidebar: () =>
        set((state) => ({ sidebarOpen: !state.sidebarOpen })),
      
      setNotificationCount: (count) => set({ notificationCount: count }),
      
      incrementNotificationCount: () =>
        set((state) => ({ notificationCount: state.notificationCount + 1 })),
      
      clearNotifications: () => set({ notificationCount: 0 }),
    }),
    {
      name: 'ai-motion-ui-storage',
      partialize: (state) => ({ theme: state.theme, sidebarOpen: state.sidebarOpen }),
    }
  )
);
