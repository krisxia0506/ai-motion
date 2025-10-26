import { create } from 'zustand';
import type { User } from '../types/auth';
import { authApi } from '../services/authApi';

interface AuthState {
  user: User | null;
  token: string | null;
  isAuthenticated: boolean;
  loading: boolean;
  error: string | null;
  
  setUser: (user: User | null) => void;
  setToken: (token: string | null) => void;
  setLoading: (loading: boolean) => void;
  setError: (error: string | null) => void;
  logout: () => void;
  initialize: () => void;
}

export const useAuthStore = create<AuthState>((set) => ({
  user: null,
  token: null,
  isAuthenticated: false,
  loading: false,
  error: null,

  setUser: (user) => set({ user, isAuthenticated: !!user }),
  setToken: (token) => set({ token }),
  setLoading: (loading) => set({ loading }),
  setError: (error) => set({ error }),
  
  logout: () => {
    authApi.logout();
    set({ user: null, token: null, isAuthenticated: false });
  },

  initialize: () => {
    const token = authApi.getStoredToken();
    const user = authApi.getStoredUser();
    if (token && user) {
      set({ token, user, isAuthenticated: true });
    }
  },
}));
