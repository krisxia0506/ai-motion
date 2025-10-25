import { create } from 'zustand';
import type { Scene } from '../types';

interface SceneState {
  scenes: Scene[];
  selectedScene: Scene | null;
  loading: boolean;
  error: string | null;
  
  setScenes: (scenes: Scene[]) => void;
  addScene: (scene: Scene) => void;
  updateScene: (id: string, updates: Partial<Scene>) => void;
  removeScene: (id: string) => void;
  setSelectedScene: (scene: Scene | null) => void;
  getScenesByNovel: (novelId: string) => Scene[];
  getScenesByChapter: (chapterId: string) => Scene[];
  setLoading: (loading: boolean) => void;
  setError: (error: string | null) => void;
  clearError: () => void;
}

export const useSceneStore = create<SceneState>((set, get) => ({
  scenes: [],
  selectedScene: null,
  loading: false,
  error: null,

  setScenes: (scenes) => set({ scenes }),
  
  addScene: (scene) =>
    set((state) => ({ scenes: [...state.scenes, scene] })),
  
  updateScene: (id, updates) =>
    set((state) => ({
      scenes: state.scenes.map((scene) =>
        scene.id === id ? { ...scene, ...updates } : scene
      ),
      selectedScene:
        state.selectedScene?.id === id
          ? { ...state.selectedScene, ...updates }
          : state.selectedScene,
    })),
  
  removeScene: (id) =>
    set((state) => ({
      scenes: state.scenes.filter((scene) => scene.id !== id),
      selectedScene:
        state.selectedScene?.id === id ? null : state.selectedScene,
    })),
  
  setSelectedScene: (scene) => set({ selectedScene: scene }),
  
  getScenesByNovel: (novelId) => {
    const { scenes } = get();
    return scenes.filter((scene) => scene.novelId === novelId);
  },
  
  getScenesByChapter: (chapterId) => {
    const { scenes } = get();
    return scenes.filter((scene) => scene.chapterId === chapterId);
  },
  
  setLoading: (loading) => set({ loading }),
  setError: (error) => set({ error }),
  clearError: () => set({ error: null }),
}));
