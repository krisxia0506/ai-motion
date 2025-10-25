import { useState, useEffect } from 'react';
import { Novel, ApiError } from '../types';
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

  const fetchNovel = async () => {
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
  };

  useEffect(() => {
    if (id) {
      fetchNovel();
    }
  }, [id]);

  return {
    novel: selectedNovel,
    loading,
    error,
    refetch: fetchNovel,
  };
};
