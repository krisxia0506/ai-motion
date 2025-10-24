import { useEffect, useState } from 'react';
import { novelApi } from '../services/novelApi';
import type { Novel } from '../types';

export const useNovel = (id: string) => {
  const [novel, setNovel] = useState<Novel | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchNovel = async () => {
      if (!id) {
        setLoading(false);
        return;
      }

      setLoading(true);
      setError(null);
      
      try {
        const fetchedNovel = await novelApi.getNovel(id);
        setNovel(fetchedNovel);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Failed to fetch novel');
      } finally {
        setLoading(false);
      }
    };

    fetchNovel();
  }, [id]);

  return { novel, loading, error };
};
