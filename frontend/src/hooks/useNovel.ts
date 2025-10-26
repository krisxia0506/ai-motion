import { useState, useEffect, useCallback } from 'react';
import type { Novel, ApiError } from '../types';
import { novelApi } from '../services';
import { useNovelStore } from '../store';

interface UseNovelResult {
  novel: Novel | null;
  loading: boolean;
  error: ApiError | null;
  refetch: () => Promise<void>;
}

export const useNovel = (id: string): UseNovelResult => {
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<ApiError | null>(null);
  const { selectedNovel, setSelectedNovel, setError: setStoreError } = useNovelStore();

  const fetchNovel = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);
      const response = await novelApi.getNovel(id);
      setSelectedNovel(response.data);
    } catch (err) {
      const apiError = err as ApiError;
      setError(apiError);
      setStoreError(apiError.message);
    } finally {
      setLoading(false);
    }
  }, [id, setSelectedNovel, setStoreError]);

  useEffect(() => {
    if (id) {
      fetchNovel();
    }
  }, [id, fetchNovel]);

  return {
    novel: selectedNovel,
    loading,
    error,
    refetch: fetchNovel,
  };
};
