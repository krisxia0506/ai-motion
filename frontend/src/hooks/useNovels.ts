import { useState, useEffect, useCallback } from 'react';
import type { Novel, ApiError, ListQueryParams, Pagination } from '../types';
import { novelApi } from '../services';
import { useNovelStore } from '../store';

interface UseNovelsResult {
  novels: Novel[];
  loading: boolean;
  error: ApiError | null;
  pagination: Pagination | null;
  refetch: () => Promise<void>;
}

export const useNovels = (params?: ListQueryParams): UseNovelsResult => {
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<ApiError | null>(null);
  const [pagination, setPagination] = useState<Pagination | null>(null);
  const { novels, setNovels, setError: setStoreError } = useNovelStore();

  const fetchNovels = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);
      const response = await novelApi.listNovels(params);
      setNovels(response.data.data);
      setPagination(response.data.pagination);
    } catch (err) {
      const apiError = err as ApiError;
      setError(apiError);
      setStoreError(apiError.message);
    } finally {
      setLoading(false);
    }
  }, [params, setNovels, setStoreError]);

  useEffect(() => {
    fetchNovels();
  }, [fetchNovels]);

  return {
    novels,
    loading,
    error,
    pagination,
    refetch: fetchNovels,
  };
};
