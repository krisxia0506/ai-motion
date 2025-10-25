import { create } from 'zustand';
import type { GenerationTask } from '../types';

interface GenerationState {
  tasks: GenerationTask[];
  activeTasks: GenerationTask[];
  completedTasks: GenerationTask[];
  
  addTask: (task: GenerationTask) => void;
  updateTask: (id: string, updates: Partial<GenerationTask>) => void;
  removeTask: (id: string) => void;
  getTaskById: (id: string) => GenerationTask | undefined;
  getTasksByScene: (sceneId: string) => GenerationTask[];
  clearCompletedTasks: () => void;
}

export const useGenerationStore = create<GenerationState>((set, get) => ({
  tasks: [],
  activeTasks: [],
  completedTasks: [],

  addTask: (task) =>
    set((state) => {
      const tasks = [...state.tasks, task];
      return {
        tasks,
        activeTasks: tasks.filter((t) => t.status === 'processing'),
        completedTasks: tasks.filter((t) => t.status === 'completed'),
      };
    }),

  updateTask: (id, updates) =>
    set((state) => {
      const tasks = state.tasks.map((task) =>
        task.id === id ? { ...task, ...updates } : task
      );
      return {
        tasks,
        activeTasks: tasks.filter((t) => t.status === 'processing'),
        completedTasks: tasks.filter((t) => t.status === 'completed'),
      };
    }),

  removeTask: (id) =>
    set((state) => {
      const tasks = state.tasks.filter((task) => task.id !== id);
      return {
        tasks,
        activeTasks: tasks.filter((t) => t.status === 'processing'),
        completedTasks: tasks.filter((t) => t.status === 'completed'),
      };
    }),

  getTaskById: (id) => {
    const { tasks } = get();
    return tasks.find((task) => task.id === id);
  },

  getTasksByScene: (sceneId) => {
    const { tasks } = get();
    return tasks.filter((task) => task.sceneId === sceneId);
  },

  clearCompletedTasks: () =>
    set((state) => ({
      tasks: state.tasks.filter((t) => t.status !== 'completed'),
      completedTasks: [],
    })),
}));
