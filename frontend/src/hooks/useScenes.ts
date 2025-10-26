import { useState, useEffect, useCallback } from 'react';
import type { Scene, ApiError } from '../types';
import { sceneApi } from '../services';
import { useSceneStore } from '../store';

interface UseScenesResult {
  scenes: Scene[];
  loading: boolean;
  error: ApiError | null;
  refetch: () => Promise<void>;
}

export const useScenes = (novelId: string, chapterId?: string): UseScenesResult => {
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<ApiError | null>(null);
  const { scenes, setScenes, setError: setStoreError } = useSceneStore();

  const fetchScenes = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);
      const response = chapterId
        ? await sceneApi.getScenesByChapter(novelId, chapterId)
        : await sceneApi.getScenesByNovel(novelId);
      setScenes(response.data);
    } catch (err) {
      const apiError = err as ApiError;
      setError(apiError);
      setStoreError(apiError.message);
    } finally {
      setLoading(false);
    }
  }, [novelId, chapterId, setScenes, setStoreError]);

  useEffect(() => {
    if (novelId) {
      fetchScenes();
    }
  }, [novelId, fetchScenes]);

  return {
    scenes: chapterId
      ? scenes.filter((s: Scene) => s.chapterId === chapterId)
      : scenes.filter((s: Scene) => s.novelId === novelId),
    loading,
    error,
    refetch: fetchScenes,
  };
};
