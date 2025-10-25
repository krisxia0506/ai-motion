import { create } from 'zustand';
import type { Novel } from '../types';

interface NovelState {
  novels: Novel[];
  selectedNovel: Novel | null;
  loading: boolean;
  error: string | null;
  
  setNovels: (novels: Novel[]) => void;
  addNovel: (novel: Novel) => void;
  updateNovel: (id: string, updates: Partial<Novel>) => void;
  removeNovel: (id: string) => void;
  setSelectedNovel: (novel: Novel | null) => void;
  setLoading: (loading: boolean) => void;
  setError: (error: string | null) => void;
  clearError: () => void;
}

export const useNovelStore = create<NovelState>((set) => ({
  novels: [],
  selectedNovel: null,
  loading: false,
  error: null,

  setNovels: (novels) => set({ novels }),
  
  addNovel: (novel) => 
    set((state) => ({ novels: [...state.novels, novel] })),
  
  updateNovel: (id, updates) =>
    set((state) => ({
      novels: state.novels.map((novel) =>
        novel.id === id ? { ...novel, ...updates } : novel
      ),
      selectedNovel:
        state.selectedNovel?.id === id
          ? { ...state.selectedNovel, ...updates }
          : state.selectedNovel,
    })),
  
  removeNovel: (id) =>
    set((state) => ({
      novels: state.novels.filter((novel) => novel.id !== id),
      selectedNovel: state.selectedNovel?.id === id ? null : state.selectedNovel,
    })),
  
  setSelectedNovel: (novel) => set({ selectedNovel: novel }),
  setLoading: (loading) => set({ loading }),
  setError: (error) => set({ error }),
  clearError: () => set({ error: null }),
}));
