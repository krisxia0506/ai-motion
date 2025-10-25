import { useState, useCallback } from 'react';
import type { GenerationTask, GenerateSceneRequest, ApiError } from '../types';
import { sceneApi, generationApi } from '../services';
import { useGenerationStore } from '../store';

interface UseGenerationResult {
  generating: boolean;
  error: ApiError | null;
  generateScene: (request: GenerateSceneRequest) => Promise<void>;
  cancelGeneration: (taskId: string) => Promise<void>;
  activeTasks: GenerationTask[];
}

export const useGeneration = (): UseGenerationResult => {
  const [generating, setGenerating] = useState(false);
  const [error, setError] = useState<ApiError | null>(null);
  const { addTask, removeTask, activeTasks } = useGenerationStore();

  const generateScene = useCallback(async (request: GenerateSceneRequest) => {
    try {
      setGenerating(true);
      setError(null);
      const response = await sceneApi.generateScene(request);
      
      const task: GenerationTask = {
        id: response.data.id,
        type: request.type,
        sceneId: request.sceneId,
        status: 'processing',
        progress: 0,
        startedAt: new Date().toISOString(),
        resultUrl: response.data.url,
      };
      
      addTask(task);
    } catch (err) {
      const apiError = err as ApiError;
      setError(apiError);
    } finally {
      setGenerating(false);
    }
  }, [addTask]);

  const cancelGeneration = useCallback(async (taskId: string) => {
    try {
      await generationApi.cancelTask(taskId);
      removeTask(taskId);
    } catch (err) {
      const apiError = err as ApiError;
      setError(apiError);
    }
  }, [removeTask]);

  return {
    generating,
    error,
    generateScene,
    cancelGeneration,
    activeTasks,
  };
};
