import { create } from 'zustand';
import type { Novel } from '../types';

interface NovelState {
  novels: Novel[];
  selectedNovel: Novel | null;
  loading: boolean;
  error: string | null;
  setNovels: (novels: Novel[]) => void;
  addNovel: (novel: Novel) => void;
  removeNovel: (id: string) => void;
  setSelectedNovel: (novel: Novel | null) => void;
  setLoading: (loading: boolean) => void;
  setError: (error: string | null) => void;
}

export const useNovelStore = create<NovelState>((set) => ({
  novels: [],
  selectedNovel: null,
  loading: false,
  error: null,
  setNovels: (novels) => set({ novels }),
  addNovel: (novel) => set((state) => ({ novels: [...state.novels, novel] })),
  removeNovel: (id) => set((state) => ({ novels: state.novels.filter((n) => n.id !== id) })),
  setSelectedNovel: (novel) => set({ selectedNovel: novel }),
  setLoading: (loading) => set({ loading }),
  setError: (error) => set({ error }),
}));
