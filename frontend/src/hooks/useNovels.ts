import { useEffect } from 'react';
import { useNovelStore } from '../store/novelStore';
import { novelApi } from '../services/novelApi';

export const useNovels = () => {
  const { novels, loading, error, setNovels, setLoading, setError } = useNovelStore();

  useEffect(() => {
    const fetchNovels = async () => {
      setLoading(true);
      setError(null);
      
      try {
        const fetchedNovels = await novelApi.listNovels();
        setNovels(fetchedNovels);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Failed to fetch novels');
      } finally {
        setLoading(false);
      }
    };

    if (novels.length === 0) {
      fetchNovels();
    }
  }, []);

  return { novels, loading, error };
};
